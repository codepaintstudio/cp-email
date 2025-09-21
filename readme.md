# CP-Email

## 项目简介

CP-Email 是一个基于 Go 后端和 Vue.js 前端的邮件群发工具，旨在解决批量发送定制化邮件的需求。该工具支持多种邮箱服务（如 QQ 邮箱、163 邮箱等），并提供了友好的 Web 界面和 API 接口。

## 功能特性

- **批量邮件发送**：支持一次性向多个收件人发送邮件
- **个性化定制**：支持在邮件内容中使用变量，实现每封邮件的个性化定制
- **用户认证**：基于 JWT 的用户登录验证系统
- **模板管理**：支持 Excel 模板上传和下载
- **邮件统计**：提供邮件发送统计信息
- **文件上传**：支持文件上传功能

## 项目结构

```
cp-email/
├── client/          # Vue.js 前端项目
│   ├── src/
│   └── package.json
└── server/          # Go 后端项目
    ├── cmd/server/  # 主程序入口
    ├── internal/    # 内部包
    ├── configs/     # 配置文件
    └── go.mod
```

## 安装与部署

### 前端部署

1. 进入前端目录：

   ```bash
   cd client
   ```

2. 安装依赖：

   ```bash
   npm install
   ```

3. 构建前端：
   ```bash
   npm run build
   ```

### 后端部署

1. 进入后端目录：

   ```bash
   cd server
   ```

2. 安装 Go 依赖：

   ```bash
   go mod tidy
   ```

3. 配置文件（可选）：
   在 `server/configs/` 目录下创建 `config.yaml` 文件：

   ```yaml
   server:
     port: "4443"
   jwt:
     secretKey: "your-secret-key"
     expiresIn: "24h"
   database:
     path: "./data/cpmail.db"
   aliyun:
     accessKeyId: "your-access-key-id"
     accessKeySecret: "your-access-key-secret"
     bucket: "your-bucket-name"
     region: "oss-cn-chengdu"
   ```

4. 启动服务：
   ```bash
   go run cmd/server/main.go
   ```

服务默认运行在 `4443` 端口。

## API 接口说明

所有需要认证的接口都需要在请求头中包含 JWT Token：

```
Authorization: Bearer <token>
```

### 1. 用户认证

#### 登录

**POST** `/api/login/verify`

验证邮箱账号密码并获取 JWT Token。

**请求参数：**

```json
{
  "email": "your-email@example.com",
  "password": "your-app-password"
}
```

#### 登出

**POST** `/api/login/logout`

用户登出（客户端清除 Token 即可）。

### 2. 邮件管理

#### 发送邮件

**POST** `/api/email/sendemail` 🔒

批量发送邮件。

**请求参数：**

```json
{
  "email": "your-email@example.com",
  "password": "your-app-password",
  "subject": "邮件主题",
  "receiverItemsArray": [
    [
      "recipient1@example.com",
      "时间戳",
      { "变量名1": "变量值1", "变量名2": "变量值2" }
    ]
  ],
  "content": "邮件内容，可使用 {变量名1} 的形式插入变量"
}
```

#### 获取邮件统计

**GET** `/api/email/stats` 🔒

获取邮件发送统计信息。

### 3. 模板管理

#### 下载模板

**GET** `/api/template/xlsx` 🔒

下载收件人信息 Excel 模板文件。

#### 上传模板

**POST** `/api/template/upload` 🔒

上传 Excel 模板文件。

### 4. 文件上传

**POST** `/api/uploads/file` 🔒

上传文件。

🔒 表示需要 JWT Token 认证的接口

## 使用示例

### 1. 登录获取 Token

```javascript
// 登录获取 JWT Token
const loginResponse = await axios.post(
  "http://127.0.0.1:4443/api/login/verify",
  {
    email: "your-email@example.com",
    password: "your-app-password",
  }
)

const token = loginResponse.data.token
```

### 2. 发送邮件

```javascript
// 使用 Token 发送邮件
await axios.post(
  "http://127.0.0.1:4443/api/email/sendemail",
  {
    email: "your-email@example.com",
    password: "your-app-password",
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
  },
  {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }
)
```

### 3. 获取统计信息

```javascript
// 获取邮件发送统计
const stats = await axios.get("http://127.0.0.1:4443/api/email/stats", {
  headers: {
    Authorization: `Bearer ${token}`,
  },
})
```

## 支持的邮箱服务商

目前支持以下邮箱服务商：

- QQ 邮箱 (smtp.qq.com)
- 163 邮箱 (smtp.163.com)

如需添加其他邮箱服务商，可在 `server/internal/utils/smtp/smtp.go` 文件中添加相应配置。

## 注意事项

1. **邮箱密码**：使用邮箱的授权码（应用密码），而非登录密码
2. **JWT Token**：Token 默认有效期为 24 小时
3. **阿里云 OSS**：如需使用文件上传功能，需要正确配置阿里云相关参数
4. **数据库**：项目使用 SQLite 数据库，数据文件默认存储在 `server/data/cpmail.db`
5. **端口配置**：默认端口为 4443，可通过配置文件修改

## 许可证

本项目采用 MIT 许可证，详情请见 [LICENSE](LICENSE) 文件。
