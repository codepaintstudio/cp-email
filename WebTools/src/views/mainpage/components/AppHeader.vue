<script setup>
import { Message, ArrowDown } from '@element-plus/icons-vue'
defineProps({
  userEmail: String,
  isLoggedIn: Boolean
})

const emit = defineEmits(['logout'])
</script>

<template>
  <div class="header">
    <div class="left">
      <div class="logo">
        <img src="@/assets/logo.png" alt="" />
      </div>
      <div class="text">
        <h2>CP-EmailTools</h2>
        <p>
          Excel转邮件群发工具 |
          <span>请务必将Excel表格第一列设置为邮箱地址</span>
        </p>
      </div>
    </div>
    <div class="right">
      <el-icon color="#000"><Message /></el-icon>
      <span>{{ userEmail }}</span>
      <el-dropdown>
        <span class="el-dropdown-link">
          <el-icon class="el-icon--right">
            <arrow-down />
          </el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu class="menu-button">
            <el-button
              v-if="isLoggedIn"
              class="delete_bt"
              type="warning"
              @click="emit('logout')"
            >
              退出登录
            </el-button>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>
