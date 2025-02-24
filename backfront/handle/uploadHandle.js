// 此文件实现图片阿里云OSS上传
const { BadRequest } = require('http-errors')
const { sendRes, RES_CODE } = require('./errorHandle')
const { singleFileUpload } = require('../utils/aliyun')

const handleAliyunUpload = (req, res) => {
  try {
    singleFileUpload(req, res, async function (error) {
      if (error) {
        return sendRes(res, error)
      }

      if (!req.file) {
        return sendRes(res, new BadRequest('请选择要上传的文件。'))
      }

      sendRes(res, RES_CODE.SUCCESS, { file: req.file.url })
    })
  } catch (error) {
    sendRes(res, error)
  }
}

module.exports = {
  handleAliyunUpload,
}
