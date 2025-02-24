import axios, {
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig
} from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user' // 假设你使用 Pinia 或 Vuex 存储 Token
import router from '../router'
import { nextTick } from 'vue'

// 定义基本的接口响应格式
interface ApiResponse<T = unknown> {
  code: number
  data: T
  msg?: string
}

const baseURL = 'http://127.0.0.1:4443/api'

const instance = axios.create({
  baseURL
})

// 请求拦截器
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const UserStore = useUserStore()
    const token = UserStore.Token
    // 检验token是否过期
    if (token && config.headers) {
      config.headers['Authorization'] = token // 在请求头中加上Token
    }
    // 确保返回 config，避免 undefined
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  (res: AxiosResponse<ApiResponse>) => {
    if (res.data.code === 0 || res.status === 200) {
      return res
    }
    console.log(res.data)
    ElMessage.error(res.data.msg || '服务异常')
    return Promise.reject(res.data)
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      nextTick(() => {
        router.push('/')
      })
    }
    return Promise.reject(error)
  }
)

export default instance
export { baseURL }
