// 此文件实现图片阿里云OSS上传
const { BadRequest } = require('http-errors')
const { sendRes, RES_CODE } = require('./errorHandle')
const { singleFileUpload, client } = require('../utils/aliyun')
const crypto = require('crypto')
const path = require('path')

// 计算文件哈希值的函数
const calculateBufferHash = (buffer) => {
  const hash = crypto.createHash('md5')
  hash.update(buffer)
  return hash.digest('hex')
}

const handleAliyunUpload = (req, res) => {
  try {
    singleFileUpload(req, res, async function (error) {
      if (error) {
        return sendRes(res, error)
      }

      if (!req.file) {
        return sendRes(res, new BadRequest('请选择要上传的文件。'))
      }
      try {
        const fileBuffer = req.file.buffer
        const fileHash = calculateBufferHash(fileBuffer)
        const fileExt = path.extname(req.file.originalname)
        // 生成文件名
        console.log('fileHash', fileHash+fileExt)
        const ossPath = `uploads/${fileHash}${fileExt}`
        try {
          // 尝试文件上传
          const result = await client.put(ossPath, fileBuffer, {
            headers: {
              'Content-Type': req.file.mimetype,
              'x-oss-object-acl': 'public-read',
              // 只有在文件不存在时才上传
              'x-oss-forbid-overwrite': 'true',
            },
          })
          //成功 说明第一次上传
          sendRes(res, RES_CODE.SUCCESS, {
            file: result.url,
            isExisting: false,
          })
        } catch (error) {
          // 如果文件已存在（上传被禁止），获取现有文件的URL
          if (error.code === 'FileAlreadyExists') {
            console.log('文件已存在，获取现有文件的URL')
            const existingUrl = client.signatureUrl(ossPath)
            return sendRes(res, RES_CODE.SUCCESS, {
              file: existingUrl,
              isExisting: true,
            })
          }
          throw error
        }
      } catch (error) {
        console.error('文件处理失败:', error)
        return sendRes(res, RES_CODE.ERROR, '文件上传失败')
      }
    })
  } catch (error) {
    sendRes(res, error)
  }
}

module.exports = {
  handleAliyunUpload,
}
