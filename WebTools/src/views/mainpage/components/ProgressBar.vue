<script setup lang='ts'>
import { defineProps } from 'vue'
defineProps<{
  isActive: boolean,
  percentage: number
}>()
</script>

<template>
  <el-affix :offset="0" class="fixed" v-show="isActive">
    <el-progress :percentage="percentage" status="success" :indeterminate="true" :duration="5" />
  </el-affix>
</template>
