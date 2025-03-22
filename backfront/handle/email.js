//此文件实现发送邮箱业务
const { sendRes, RES_CODE } = require('./errorHandle');
const { createSmtpServer, sendMail } = require('./nodemailerHandle')
/**
 * SMTP服务请求队列： 先来先服务
*/
const smtpServerQueue = [];

//因为是请求队列，这个函数其实不需要，如果后续有限制发送间隔可以用此方案
const frozen = async (time) => {
  await setTimeout(() => {}, time)
}

/** 去除数组中所有假值 */
const compactArr = (arr) => {
  const newArr = [];
  arr.forEach((item) => {
    if(item) newArr.push(item);
  });
  return newArr;
}

/** 拼接两个数组后转为字符串 */
const joinArrays = (arr1, arr2) => {
  if (arr2.length) {
    const newArr = [];
    arr1.forEach((item, index) => {
      newArr.push(arr1[index], arr2[index]);
    });
    return newArr.join('');
  } else {
    return arr1.join('');
  }
}

let running = false;
/**
 * 执行队列
*/
const executeServerQueue = async (fun) => {
  // 如果传入参数 则把fun压入服务队列中
  if (fun) {
    smtpServerQueue.push(fun);
  }
  console.log(`调用执行队列，当前请求服务队列有${smtpServerQueue.length}个服务`, smtpServerQueue);
  // 检查服务队列中是否有值，有值且当前没有函数在执行则执行剩余函数
  if (smtpServerQueue.length && !running) {
    running = true;
    console.log('执行！');
    await smtpServerQueue[0]();
    // 这里可以限制一下发送间隔
    // await frozen(1000);
    running = false;
    console.log('执行结束！');
    //等待函数执行完进行出队
    smtpServerQueue.shift();
    executeServerQueue();
  }
}

/**
 * 预校验
 * TODO： 需要校验长度是否合法
*/
const preError = (req) => {
  const {
    email = null,
    password = null,
    subject = null,
    receiverItemsArray = null,
    content = null,
  } = req;
  if (!email) {
    console.log('邮箱不能为空！');
    return '邮箱不能为空！';
  }
  if (!password) {
    console.log('密码不能为空！');
    return '密码不能为空！';
  }
  if (!subject) {
    console.log('邮件主题不能为空！');
    return '邮件主题不能为空！';
  }
  if (!receiverItemsArray) {
    console.log('联系人列表不能为空！');
    return '联系人列表不能为空！';
  }
  if (!content) {
    console.log('正文内容不能为空！');
    return '正文内容不能为空';
  }
  return '';
}

/** 
 * 批量发送邮件
*/
const sendEmails = async (req, res) => {
  if (preError(req.body)) {
    sendRes(
      res,
      RES_CODE.FIELD_IS_EMPTY,
      {
        msg: preError(req.body),
      }
    )
    return;
  }
  const {
    email,
    password,
    subject,
    receiverItemsArray,
    content,
  } = req.body;
  // console.log('开始发送邮件');
  const smtpServerObj = await createSmtpServer(email, password);
  //没创建成功直接返回前端信息并退出
  if (!smtpServerObj.success) {
    sendRes(
      res,
      RES_CODE.ACCOUNT_OR_PASSWORD_ERROR,
      {
        msg: smtpServerObj.info,
      }
    )
    return;
  }
  const smtpServer = smtpServerObj.smtpServer;
  // 处理邮件发送
  const emailList = receiverItemsArray.map(items => {
    const [receiverEmail, timestamp, variables] = items;
    
    // 动态替换内容
    let replacedContent = content;
    for (const [varName, varValue] of Object.entries(variables)) {
      const regex = new RegExp(`{${varName}}`, 'g');
      replacedContent = replacedContent.replace(regex, varValue ?? '');
    }

    // 构造消息
    const message = {
      from: `<${email}>`,
      to: receiverEmail,
      subject: subject,
      text: replacedContent,
      html: replacedContent 
    };

    return sendMail(message, smtpServerObj.smtpServer);
  });
  console.log(emailList)
  const messages = compactArr(await Promise.all(emailList));
  sendRes(
    res,
    RES_CODE.SUCCESS,
    {
      successList: messages,
    }
  );
  //关闭smtp服务 TODO: check
  await new Promise(async (res) => {
    await smtpServer.close();
    res('');
  })
  console.log('服务关闭');
}

module.exports = {
  executeServerQueue,
  sendEmails,
}