const express = require('express')
const router = express.Router()
const { handleTemplateDownload } = require('../handle/templateHandle')

router.get('/xlsx', handleTemplateDownload)

module.exports = router
