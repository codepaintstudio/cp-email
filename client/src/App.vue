<script setup lang="ts">
import IntroDialog from './components/IntroDialog.vue'
import { ref } from 'vue';

const showIntro = ref<boolean>(localStorage.getItem('hideIntro') !== 'true');
</script>

<template>
  <div>
    <!-- 页面路由内容 -->
    <router-view />
    
    <!-- 使用介绍 -->
    <IntroDialog v-model="showIntro" />

    <!-- 介绍按钮弹窗 -->
    <div class="explain" v-show="!showIntro">
      <el-button plain @click="showIntro = true" class="delete_bt explain_button" type="info">
        ?
      </el-button>
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
