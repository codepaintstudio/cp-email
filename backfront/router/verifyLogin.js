const express = require('express')
const router = express.Router()
const { verifyLogin } = require('../handle/loginHandle')

router.post('/verify', verifyLogin)

module.exports = router
