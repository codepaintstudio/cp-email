import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore(
  'user',
  () => {
    const Token = ref<string>('')
    const email = ref<string>('')
    const password = ref<string>('')

    function setToken(e: string) {
      Token.value = e
    }

    function setEmail(e: string) {
      email.value = e
    }

    function setPassword(e: string) {
      password.value = e
    }

    function clearToken() {
      Token.value = ''
      email.value = ''
      password.value = ''
    }

    const hasToken = () => !!Token.value

    return {
      Token,
      setToken,
      clearToken,
      hasToken,
      email,
      password,
      setEmail,
      setPassword
    }
  },
  {
    persist: true
  }
)
