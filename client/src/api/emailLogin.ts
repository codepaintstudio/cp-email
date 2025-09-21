import request from '@/utils/request'
// 登录api

interface LoginRequestData {
  email: string
  password: string
}
export const LoginVerify = (data: LoginRequestData) =>
  request.post('/login/verify', data)
