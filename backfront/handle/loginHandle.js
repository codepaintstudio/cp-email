// 此文件实现登录功能
const jwt = require('jsonwebtoken')
const nodemailer = require('nodemailer')
const { selectSmtpConfig } = require('./nodemailerHandle')
const { sendRes, RES_CODE } = require('./errorHandle')
const config = require('../config')

const verifyLogin = async (req, res) => {
  const { email, password } = req.body

  if (!email || !password) {
    return sendRes(res, RES_CODE.PARAM_ERROR, '邮箱和密码不能为空！')
  }

  const smtpConfig = selectSmtpConfig(email)
  if (!smtpConfig) {
    return sendRes(res, RES_CODE.PARAM_ERROR, '暂不支持此类邮箱！')
  }

  const smtpServer = await nodemailer.createTransport({
    ...smtpConfig,
    auth: { user: email, pass: password },
  })

  smtpServer.verify(function (error, success) {
    if (error) {
      return sendRes(res, RES_CODE.PARAM_ERROR, '请检查账号密码是否正确！')
    }
    const token = jwt.sign({ email }, config.jwtSecretKey, {
      expiresIn: config.expiresIn,
    })
    return sendRes(res, RES_CODE.SUCCESS, {
      msg: '账号密码验证成功',
      token: 'Bearer ' + token,
    })
  })
}

module.exports = {
  verifyLogin,
}
