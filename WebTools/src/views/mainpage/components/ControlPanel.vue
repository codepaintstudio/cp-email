<script setup>
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
</script>

<template>
  <div class="main">
    <div class="left-buttons">
      <el-upload
        :limit="1"
        :show-file-list="false"
        accept=".xlsx, .xls"
        :on-change="(file) => emit('file-change', file)"
      >
        <el-button type="primary" color="#3370FF" class="upload_bt">
          <el-icon style="margin-right: 6px"><UploadFilled /></el-icon>
          上传Excel
        </el-button>
      </el-upload>
      <el-button class="delete_bt" @click="emit('clear-data')"
        >清空数据</el-button
      >
    </div>

    <div class="right-buttons">
      <el-button
        type="primary"
        class="send_bt"
        color="#3370FF"
        :disabled="disableSend"
        @click="emit('send-emails')"
      >
        <el-icon style="margin-right: 6px"><Promotion /></el-icon>
        一键发送
      </el-button>
      <el-button class="delete_bt" @click="emit('download-template')">
        下载模板
      </el-button>
    </div>
  </div>
</template>
