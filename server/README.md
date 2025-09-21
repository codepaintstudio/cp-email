# CP Mail Server

CP Mail 是一个使用 Go 语言编写的邮件发送服务，支持批量发送邮件、模板下载、文件上传等功能。

## 功能特性

1. 用户认证（JWT）
2. 邮件批量发送
3. 邮件发送统计（使用 SQLite）
4. 模板文件下载
5. 文件上传
6. 跨域支持

## 技术栈

- Go 1.24.5
- Gin Web 框架
- GORM + SQLite
- JWT 认证

## 目录结构

```
server/
├── cmd/
│   └── server/
│       └── main.go          # 程序入口
├── configs/
│   └── config.yaml          # 配置文件
├── data/                    # SQLite 数据库文件
├── internal/
│   ├── app/                 # 应用初始化
│   ├── config/              # 配置管理
│   ├── handler/             # HTTP 处理程序
│   ├── middleware/          # 中间件
│   ├── service/             # 业务逻辑
│   └── utils/               # 工具类
│       ├── db/              # 数据库工具
│       ├── response/        # 响应处理
│       └── smtp/            # SMTP 工具
├── public/
│   ├── template/            # 模板文件
│   └── uploads/             # 上传文件
└── go.mod                   # Go 模块定义
```

## 安装和运行

1. 确保已安装 Go 1.24.5 或更高版本
2. 克隆项目代码
3. 进入项目目录
4. 运行以下命令安装依赖：
   ```
   go mod tidy
   ```
5. 编译程序：
   ```
   go build -o cpmail cmd/server/main.go
   ```
6. 运行服务：
   ```
   ./cpmail
   ```

## API 接口

### 认证接口

- `POST /api/login/verify` - 用户登录
- `POST /api/login/logout` - 用户登出

### 邮件接口

- `POST /api/email/sendemail` - 发送邮件（需要认证）
- `GET /api/email/stats` - 获取邮件发送统计（需要认证）

### 模板接口

- `GET /api/template/xlsx` - 下载模板文件（需要认证）
- `POST /api/template/upload` - 上传模板文件（需要认证）

### 上传接口

- `POST /api/uploads/file` - 上传文件（需要认证）

## 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
server:
  port: "4443"              # 服务端口

jwt:
  secretKey: "mima"         # JWT 密钥
  expiresIn: "24h"          # JWT 过期时间

database:
  path: "./data/cpmail.db"  # SQLite 数据库路径
```

## 使用说明

1. 启动服务后，默认监听 4443 端口
2. 首先调用登录接口进行认证
3. 获取到 JWT token 后，在后续请求的 Authorization 头中携带 token
4. 可以通过上传接口上传模板文件到 `public/template/todaytemplate.xlsx`
5. 邮件发送统计信息保存在 SQLite 数据库中