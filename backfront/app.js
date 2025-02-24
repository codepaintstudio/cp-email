const express = require('express')
const { sendRes, RES_CODE } = require('./handle/errorHandle')

//创建服务器实例
const server = express()
const expressJWT = require('express-jwt')
const config = require('./config')

//解决跨域
const cors = require('cors')
server.use(cors())

//配置jwt中间件
server.use(
  expressJWT
    .expressjwt({
      secret: config.jwtSecretKey,
      algorithms: ['HS256'],
    })
    .unless({
      path: ['/api/login/verify'],
    })
)

//添加错误处理中间件
server.use((err, req, res, next) => {
  // 设置CORS头
  res.header('Access-Control-Allow-Origin', 'http://localhost:5173')
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
  res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization')
  if (err.name === 'UnauthorizedError') {
    console.log(err.name)
    return res.status(401).json({
      code: 401,
      msg: 'token无效或已过期',
      data: null,
    })
  }
  next()
})

//配置实例解析项
server.use(express.urlencoded({ extends: false }))
server.use(express.json())

// 创建登录路由
const loginRouter = require('./router/verifyLogin')
// 注册登录路由
server.use('/api/login', loginRouter)

//创建邮箱路由
const sendEmail = require('./router/emailServer')
//注册邮箱路由
server.use('/api/email', sendEmail)

// 创建上传路由
const uploadsRouter = require('./router/uploads')
// 注册上传路由
server.use('/api/uploads', uploadsRouter)

// 创建下载模板路由
const templateDownload = require('./router/templateDownload')
// 注册下载模板路由
server.use('/api/template', templateDownload)

//启动服务
server.listen(4443, function () {
  const port = this.address().port // 获取被分配的端口号
  console.log(`Express server listening on port ${port}`)
})
