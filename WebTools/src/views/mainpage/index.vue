<script setup lang='ts'>
import { ref, computed, onMounted } from 'vue'
import * as XLSX from 'xlsx'
import { useUserStore } from '../../stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTemplateDownload } from '../../api/gettemplate'
import { useRouter } from 'vue-router'
import ProgressBar from './components/ProgressBar.vue'
import AppHeader from './components/AppHeader.vue'
import ControlPanel from './components/ControlPanel.vue'
import EmailEditor from './components/EmailEditor.vue'
import EmailDataTable from './components/EmailDataTable.vue'
import { sendEmailService } from '@/api/postEmail'

const router = useRouter()
const UserStore = useUserStore()

interface ExcelRow {
  email: string
  state: number
}

// 状态管理
const isClick = ref<boolean>(false)
const excelData = ref<ExcelRow[]>([{ email: '', state: 1 }])
const emailContent = ref<string>('')
const subject = ref<string>('')
const loading = ref<boolean>(false)
const sendProcess = ref<number>(100)
const acceptedEmail = ref<string[]>([])
const MAX_ROWS = 10000
const userEemil = ref<string>('')
const isLogin = ref<boolean>(UserStore.hasToken())

// 处理列
const columns = computed<string[]>(() => {
  const cols = new Set<string>()
  excelData.value.forEach((row) => {
    Object.keys(row).forEach((key) => cols.add(key))
  })
  return Array.from(cols).filter((col) => col !== 'state')
})

onMounted(() => {
  if (UserStore.hasToken()) {
    userEemil.value = UserStore.email
  } else {
    router.replace({ name: 'Login' })
  }
})

// 文件处理方法
let debounceTimer: number | null = null
const handleFileChange = async (file: File) => {
  if (!file) return
  if (debounceTimer) clearTimeout(debounceTimer)

  debounceTimer = window.setTimeout(async () => {
    loading.value = true
    try {
      const data = await readExcelFile(file)
      const workbook = XLSX.read(data, { type: 'binary' })
      const worksheet = workbook.Sheets[workbook.SheetNames[0]]
      let jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 2 }) as ExcelRow[]

      if (jsonData.length > MAX_ROWS) {
        ElMessage.error(`最多允许 ${MAX_ROWS} 行数据`)
        excelData.value = [{ email: '', state: 1 }]
        return
      }

      jsonData = jsonData.map((entry) => ({ ...entry, state: 1 }))
      const emptyRowIndex = excelData.value.findIndex((row) => !row.email)

      if (emptyRowIndex > -1) {
        excelData.value.splice(emptyRowIndex, 1, ...jsonData)
      } else {
        excelData.value = [...excelData.value, ...jsonData]
      }
    } finally {
      loading.value = false
    }
  }, 100)
}

const readExcelFile = (File: any): Promise<string | ArrayBuffer | null> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      if (e.target?.result) {
        resolve(e.target.result)
      } else {
        reject(new Error('Failed to read file'))
      }
    }
    reader.onerror = (error: ProgressEvent<FileReader>) => reject(error)
    reader.readAsBinaryString(File.raw)
  })
}

// 其他操作方法
const deleteData = () => {
  loading.value = true
  excelData.value = [{ email: '', state: 1 }]
  emailContent.value = ''
  subject.value = ''
  isClick.value = false
  sendProcess.value = 0
  loading.value = false
  ElMessage.success('数据已清空')
}

const sendEmails = async () => {
  if (!validateBeforeSend()) return
  isClick.value = true

  try {
    const receiverItems = excelData.value.map((item) => [
      item.email,
      new Date().toLocaleString()
    ])

    const res = await sendEmailService({
      email: UserStore.email,
      password: UserStore.password,
      subject: subject.value,
      receiverItemsArray: receiverItems,
      content: emailContent.value
    })

    acceptedEmail.value = res.data.data.successList
    updateEmailStates()
    if(acceptedEmail.value.length > 0){
      ElMessage.success('邮件发送完成')
    }else{
      ElMessage.error('请检插邮件号是否符合要求')
    }
  } finally {
    isClick.value = false
  }
}

const validateBeforeSend = () => {
  if (excelData.value[0].email === '' && excelData.value.length <= 1) {
    ElMessage.warning('请先导入Excel数据')
    return false
  }
  if (!emailContent.value) {
    ElMessage.warning('请输入邮件内容')
    return false
  }
  if (!subject.value) {
    ElMessage.warning('请输入邮件主题')
    return false
  }
  return true
}

const updateEmailStates = () => {
  excelData.value.forEach((item) => {
    const email = item.email
    item.state = acceptedEmail.value.includes(email) ? 2 : 0
  })
}

const handleDeleteRow = (index: number) => {
  if (excelData.value.length === 1) {
    ElMessage.error('至少保留一行数据')
    return
  }

  ElMessageBox.confirm('确定删除该行？')
    .then(() => {
      excelData.value.splice(index, 1)
      ElMessage.success('删除成功')
    })
    .catch(() => ElMessage.info('取消删除'))
}

const handleAddRow = () => {
  if (excelData.value.length < 100) {
    excelData.value.push({ email: '', state: 1 })
  } else {
    ElMessage.warning('最多添加100行')
  }
}

const handleDownloadTemplate = async (): Promise<void> => {
  try {
    const res = await getTemplateDownload()
    const blob = new Blob([res.data], { type: res.headers['content-type'] })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'template.xlsx'
    link.click()
    URL.revokeObjectURL(link.href)
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

const handleLogout = () => {
  UserStore.clearToken()
  router.replace({ name: 'Login' })
  ElMessage.success('已退出登录')
}
</script>

<template>
  <div class="box">
    <!-- 进度条 -->
    <ProgressBar :percentage="sendProcess" :is-active="isClick" />

    <div class="banner">
      <!-- 头部logo和账户信息 -->
      <AppHeader
        :user-email="userEemil"
        :is-logged-in="isLogin"
        @logout="handleLogout"
      />

      <!-- 功能区域 -->
      <ControlPanel
        :disable-send="excelData.length === 0 || isClick"
        @file-change="handleFileChange"
        @clear-data="deleteData"
        @send-emails="sendEmails"
        @download-template="handleDownloadTemplate"
      />

      <!-- 邮件编写区域 -->
      <EmailEditor
        :subject="subject"
        :content="emailContent"
        @update:subject="(val) => (subject = val)"
        @update:content="(val) => (emailContent = val)"
      />

      <!-- 邮件发送表格信息 -->
      <EmailDataTable
        :data="excelData"
        :columns="columns"
        :loading="loading"
        @delete-row="handleDeleteRow"
        @add-row="handleAddRow"
      />
    </div>
  </div>
  <!-- 底部版权介绍 -->
  <div class="footer">
    <p>
      Copyright © 2024 EmailTools.Designed by
      <a href="https://github.com/CodePaintStudio">CodePaint</a>
    </p>
  </div>
</template>

<style lang="scss">
:deep(.el-input__wrapper) {
  box-shadow: none;
}
.footer {
  padding: 20px 0;
  margin: 2rem 0 -8px;
  background-color: #f8f8f8;
  text-align: center;
  font-size: 14px;
  color: #666;
  border-top: 1px solid #ddd;

  a {
    color: #3370ff;
    text-decoration: none;
    font-weight: bold;

    &:hover {
      text-decoration: underline;
    }
  }

  p {
    margin: 0;
    padding: 0;
  }
}
</style>
