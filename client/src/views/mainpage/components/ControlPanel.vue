<script setup lang="ts">
import { ref } from 'vue'
import { UploadFilled, Promotion } from '@element-plus/icons-vue'

defineProps({
  disableSend: Boolean
})

const emit = defineEmits([
  'file-change',
  'clear-data',
  'send-emails',
  'download-template'
])

const uploadKey = ref(0)
// 删除数据的弹窗
const clearDialogVisible = ref(false)
// 选择删除的数据
const selectedClearOptions = ref<string[]>([])
// 删除数据的选项
const clearOptions = [
  { value: 'subject', label: '邮件主题' },
  { value: 'content', label: '邮件内容' },
  { value: 'excel', label: '收件人数据' }
]

const openClearDialog = () => {
  selectedClearOptions.value = ['subject', 'content', 'excel'] // 默认全选
  clearDialogVisible.value = true
}

const confirmClear = () => {
  emit('clear-data', selectedClearOptions.value)
  clearDialogVisible.value = false
}

const handleFileChange = (file: any) => {
  emit('file-change', file)
  uploadKey.value++  // 强制刷新上传组件
}
</script>

<template>
  <div class="main">
    <div class="left-buttons">
      <el-upload :key="uploadKey" :limit="1" :show-file-list="false" accept=".xlsx, .xls" :on-change="handleFileChange"
                 :auto-upload="false" :http-request="() => { }">
        <el-button type="primary" color="#3370FF" class="upload_bt">
          <el-icon style="margin-right: 6px">
            <UploadFilled />
          </el-icon>
          上传Excel
        </el-button>
      </el-upload>
      <el-button class="delete_bt" @click="openClearDialog">清空数据</el-button>
      <el-dialog v-model="clearDialogVisible" title="选择要清除的内容" width="400px">
        <el-select v-model="selectedClearOptions" multiple placeholder="请选择要清除的内容" style="width: 100%">
          <el-option v-for="item in clearOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <template #footer>
          <el-button @click="clearDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmClear">确定清除</el-button>
        </template>
      </el-dialog>
    </div>

    <div class="right-buttons">
      <el-button type="primary" class="send_bt" color="#3370FF" :disabled="disableSend" @click="emit('send-emails')">
        <el-icon style="margin-right: 6px">
          <Promotion />
        </el-icon>
        一键发送
      </el-button>
      <el-button class="delete_bt" @click="emit('download-template')">
        下载模板
      </el-button>
    </div>
  </div>
</template>