// 此文件实现模板下载功能
const path = require('path')
const fs = require('fs')
const { sendRes, RES_CODE } = require('./errorHandle')

const handleTemplateDownload = (req, res) => {
  const filePath = path.join(process.cwd(), 'public', 'template', 'todaytemplate.xlsx')

  if (!fs.existsSync(filePath)) {
    console.error('文件不存在:', filePath)
    return sendRes(res, RES_CODE.FIELD_IS_EMPTY, '文件不存在，请检查路径或文件是否已上传')
  }

  try {
    const fileStream = fs.createReadStream(filePath)
    res.setHeader('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')
    res.setHeader('Content-Disposition', 'attachment; filename="template.xlsx"')

    fileStream.pipe(res)

    fileStream.on('error', (error) => {
      console.error('文件读取失败:', error)
      sendRes(res, 500, '文件读取失败')
    })

    fileStream.on('end', () => {
      console.log('文件发送成功')
    })
  } catch (error) {
    console.error('文件处理失败:', error)
    sendRes(res, 500, '文件处理失败')
  }
}

module.exports = {
  handleTemplateDownload,
}
