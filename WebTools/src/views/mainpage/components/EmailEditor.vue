<script setup>
import { defineProps, defineEmits } from 'vue'
import Tiptap from '@/components/Tiptap.vue'

defineProps({
  subject: String,
  content: String
})

const emit = defineEmits(['update:subject', 'update:content'])
</script>

<template>
  <div class="input">
    <p class="label">主题：</p>
    <el-input
      class="subject"
      :modelValue="subject"
      @update:modelValue="(val) => emit('update:subject', val)"
      placeholder="请输入邮件主题"
    />
  </div>
  <div class="show">
    <Tiptap
      :modelValue="content"
      @update:modelValue="(val) => emit('update:content', val)"
    />
  </div>
</template>

<style lang="scss" scoped>
:deep(.el-input__wrapper) {
  box-shadow: none;
}
</style>
