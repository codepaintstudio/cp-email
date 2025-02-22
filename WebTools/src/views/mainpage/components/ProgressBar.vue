<script setup>
defineProps({
  percentage: Number,
  isActive: Boolean
})
</script>

<template>
  <el-affix :offset="0" class="fixed" v-show="isActive">
    <el-progress
      :percentage="percentage"
      status="success"
      :indeterminate="true"
      :duration="5"
    />
  </el-affix>
</template>
