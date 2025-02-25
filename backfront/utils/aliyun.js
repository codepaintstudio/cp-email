const multer = require('multer')
const MAO = require('multer-aliyun-oss')
const OSS = require('ali-oss')
const { BadRequest } = require('http-errors')
require('dotenv').config()

// 阿里云配置信息
const config = {
  region: process.env.ALIYUN_REGION,
  accessKeyId: process.env.ALIYUN_ACCESS_KEY_ID,
  accessKeySecret: process.env.ALIYUN_ACCESS_KEY_SECRET,
  bucket: process.env.ALIYUN_BUCKET,
}

const client = new OSS(config)

// multer 配置信息
const upload = multer({
  storage: multer.memoryStorage(),  // 使用内存存储
  limits: {
    fileSize: 5 * 1024 * 1024,
  },
  fileFilter: function (req, file, cb) {
    const fileType = file.mimetype.split('/')[0]
    const isImage = fileType === 'image'

    if (!isImage) {
      return cb(new BadRequest('只允许上传图片。'))
    }
    cb(null, true)
  },
})

// 单文件上传，指定表单字段名为 file
const singleFileUpload = upload.single('file')

module.exports = {
  config,
  client,
  singleFileUpload,
}
