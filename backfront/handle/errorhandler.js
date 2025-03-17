// 定义错误中间件
function errorHandler(err, req, res, next) {

  if (err.name === 'UnauthorizedError') {
    console.log(err.name)
    return res.status(401).json({
      code: 401,
      msg: 'token无效或已过期',
      data: null,
    })
  }
  next()
}

module.exports = errorHandler