<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import { useUserStore } from '../../stores/user'
import { useRouter } from 'vue-router'
import { LoginVerify } from '../../api/emailLogin'

interface LoginData {
  email: string
  password: string
}

const router = useRouter()
const UserStore = useUserStore()
const formRef = ref<any>(null)


const emailData = ref<LoginData>({
  email: '',
  password: ''
})


const rules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    {
      type: 'email',
      message: '请输入正确的邮箱地址',
      trigger: ['blur', 'change']
    }
  ],
  password: [
    { required: true, message: '请输入邮箱授权码', trigger: 'blur' },
    { min: 6, message: '授权码至少为6个字符', trigger: 'blur' },
    { max: 16, message: '授权码不能超过16个字符', trigger: 'blur' }
  ]
}

onMounted(() => {
  emailData.value = {
    email: UserStore.email,
    password: UserStore.password
  }
})

async function onSubmit() {
  formRef.value.validate(async (valid: boolean) => {
    if (!valid) {
      ElMessage.error('请填写完整的登录信息')
      return false
    }

    const loading = ElLoading.service({
      lock: true,
      text: '登录中...',
      background: 'rgba(0, 0, 0, 0.7)'
    })

    try {
      const response = await LoginVerify({
        email: emailData.value.email,
        password: emailData.value.password
      })

      if (response.data.code === 0) {
        ElMessage({
          message: '登录成功',
          type: 'success'
        })
        loading.close()

        // 存储用户信息
        UserStore.setToken(response.data.data.token)
        UserStore.setEmail(emailData.value.email)
        UserStore.setPassword(emailData.value.password)

        // 跳转到主页面
        router.push('/main')
      } else {
        loading.close()
        ElMessage.error(response.data.msg || '登录失败，请检查邮箱和授权码')
      }
    } catch (error) {
      loading.close()
      ElMessage.error('登录请求失败，请稍后再试')
    }
  })
}

function onReset() {
  formRef.value.resetFields()
}
</script>

<template>
  <div class="login-container">
    <!-- 标题和logo -->
    <div class="title">
      <div class="box"><img src="@/assets/logo.png" alt="Logo" /></div>
      <div class="text">CP-EmailTools</div>
    </div>
    <!-- 表单 -->
    <el-form ref="formRef" :model="emailData" label-width="auto" :rules="rules" class="el-form">
      <div class="text">登录</div>
      <!-- 邮箱 -->
      <el-form-item prop="email" class="el-form-item">
        <el-input v-model="emailData.email" autocomplete="off" placeholder="输入邮箱地址" class="input-field"></el-input>
      </el-form-item>
      <!-- 邮箱授权码 -->
      <el-form-item prop="password" class="el-form-item">
        <el-input type="password" v-model="emailData.password" placeholder="请输入邮箱授权码" autocomplete="off"
                  class="input-field"></el-input>
      </el-form-item>
      <!-- 确定取消按钮 -->
      <div class="bt">
        <el-button color="#3370FF" type="primary" @click="onSubmit" class="inner_bt">确认</el-button>
        <el-button plain="true" @click="onReset" class="inner_bt">清空</el-button>
      </div>
    </el-form>
  </div>
</template>

<style scoped lang="scss">
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  margin: -8px;
  background: conic-gradient(from 109.82431090096661deg at 74.47646856307983% -86.3436758518219%,
      #c1d2fa 1%,
      #f8f9fa 37%,
      #ffffff 49%,
      #feffff 55%);

  .title {
    display: flex;
    align-items: center;
    justify-content: center;
    transform: translateX(-15%);

    .box {
      width: 8rem;

      img {
        width: 10rem;
        height: 100px;
        object-fit: contain;
      }
    }

    .text {
      font-size: 2rem;
      font-weight: bold;
      color: #000;
    }
  }

  .el-form {
    background: #fff;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    width: 400px;
    display: flex;
    flex-direction: column;
    align-items: center;

    @media screen and (max-width: 768px) {
      width: 100%;
      max-width: 350px;
      padding: 1.5rem;
      box-sizing: border-box;
    }

    .text {
      font-size: 1.5rem;
      font-weight: bold;
      font-family: 'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;
      color: #333;
      margin-bottom: 2rem;
      width: 90%;
      justify-content: flex-start;
    }
  }

  .el-form-item {
    width: 80%;
    margin-bottom: 1.5rem;

    @media screen and (max-width: 768px) {
      width: 85%;
      margin-bottom: 1rem;
    }
  }

  .el-form-item label {
    font-size: 1rem;
    color: #333;
  }

  .input-field {
    width: 100%;

    :deep(.el-input__inner) {
      height: 40px;
      line-height: 40px;
    }
  }

  .bt {
    width: 80%;
    display: flex;
    justify-content: center;
    gap: 2rem;
    margin-top: 1rem;

    @media screen and (max-width: 768px) {
      width: 85%;
      flex-direction: column;
      gap: 1rem;
    }

    .inner_bt {
      min-width: 120px;

      @media screen and (max-width: 768px) {
        width: 100%;
        margin: 0;
        height: 40px;
      }
    }
  }
}

:deep(.el-input__inner) {
  padding-left: 10px;
}
</style>
