const express = require('express')
const router = express.Router()
const { handleAliyunUpload } = require('../handle/uploadHandle')

/**
 * 阿里云 OSS 客户端上传
 * POST /uploads/aliyun
 */
router.post('/aliyun', handleAliyunUpload)

module.exports = router
