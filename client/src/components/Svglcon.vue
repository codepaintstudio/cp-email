<template>
  <svg :style="{ width, height }">
    <use :href="prefix + name" :fill="color"></use>
  </svg>
</template>

<script setup lang="ts">
// 接收父组件传递过来的参数
defineProps({
  prefix: {
    type: String,
    default: '#icon-'
  },
  // svg名称
  name: String,
  // svg颜色
  color: {
    type: String,
    default: '#ffffff'
  },
  // svg宽度
  width: {
    type: String,
    default: '16px'
  },
  // svg高度
  height: {
    type: String,
    default: '16px'
  }
})
</script>
