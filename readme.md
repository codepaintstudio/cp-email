# CP-Email

## 项目简介

CP-Email 是一个基于 Node.js 和 Express 的邮件群发工具，旨在解决批量发送定制化邮件的需求。该工具支持多种邮箱服务（如 QQ 邮箱、163 邮箱等），并提供了友好的 API 接口，方便集成到其他系统中。

## 功能特性

- **批量邮件发送**：支持一次性向多个收件人发送邮件
- **个性化定制**：支持在邮件内容中使用变量，实现每封邮件的个性化定制
- **云存储支持**：集成阿里云 OSS，支持图片上传和管理

## 安装与部署

1. 克隆项目到本地：

   ```bash
   git clone git@github.com:codepaintstudio/cp-email.git
   cd cp-email/backfront
   ```

2. 安装依赖：

   ```bash
   npm install
   ```

3. 配置环境变量：

   ```bash
   # 复制 .env.example 文件并重命名为 .env
   # 在 .env 文件中填写阿里云相关配置
   ALIYUN_ACCESS_KEY_ID='你的阿里云AccessKeyID'
   ALIYUN_ACCESS_KEY_SECRET='你的阿里云AccessKeySecret'
   ALIYUN_BUCKET=imgurl-oss
   ALIYUN_REGION=oss-cn-chengdu
   ```

4. 启动服务：
   ```bash
   node app.js
   ```

服务默认运行在 `4443` 端口。

## API 接口说明

### 1. 登录验证

**POST** `/api/login/verify`

验证邮箱账号密码是否正确，并获取访问令牌。

**请求参数：**

```json
{
  "email": "your-email@example.com",
  "password": "your-password"
}
```

### 2. 发送邮件

**POST** `/api/email/sendemail`

发送批量邮件。

**请求参数：**

```json
{
  "email": "your-email@example.com",
  "password": "your-password",
  "subject": "邮件主题",
  "receiverItemsArray": [
    [
      "recipient1@example.com",
      "时间戳",
      { "变量名1": "变量值1", "变量名2": "变量值2" }
    ],
    [
      "recipient2@example.com",
      "时间戳",
      { "变量名1": "变量值1", "变量名2": "变量值2" }
    ]
  ],
  "content": "邮件内容，可使用 {变量名1} 的形式插入变量"
}
```

### 3. 模板下载

**GET** `/api/template/xlsx`

下载收件人信息模板文件。

### 4. 图片上传

**POST** `/api/uploads/aliyun`

上传图片到阿里云 OSS。

## 使用示例

### 前端调用示例

```javascript
// 发送邮件示例
axios.post("http://127.0.0.1:4443/api/email/sendemail", {
  email: "your-email@example.com",
  password: "your-password",
  subject: "欢迎加入我们",
  receiverItemsArray: [
    [
      "recipient1@example.com",
      "2023-01-01",
      { name: "张三", position: "开发工程师" },
    ],
    [
      "recipient2@example.com",
      "2023-01-01",
      { name: "李四", position: "产品经理" },
    ],
  ],
  content: "亲爱的{name}，欢迎应聘{position}职位！",
})
```

## 支持的邮箱服务商

目前支持以下邮箱服务商：

- QQ 邮箱 (smtp.qq.com)
- 163 邮箱 (smtp.163.com)

如需添加其他邮箱服务商，可在 `handle/nodemailerHandle.js` 文件中的 `selectSmtpConfig` 函数中添加相应配置。

## 注意事项

1. 邮箱密码为授权码，而非登录密码，请在邮箱设置中获取授权码
3. 使用阿里云 OSS 功能需要正确配置环境变量

## 许可证

本项目采用 MIT 许可证，详情请见 [LICENSE](LICENSE) 文件。
