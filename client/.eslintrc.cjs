require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  extends: [
    'plugin:vue/vue3-essential', // Vue 3 基础规则
    'eslint:recommended', // ESLint 推荐规则
    'plugin:@typescript-eslint/recommended', // TypeScript 推荐规则
    '@vue/eslint-config-prettier/skip-formatting', // 集成 Prettier
    '@vue/typescript/recommended' // Vue 官方推荐的 TypeScript 规则
  ],

  parserOptions: {
    ecmaVersion: 'latest', // 使用最新的 ECMAScript 版本
    sourceType: 'module', // 使用 ES 模块
    parser: '@typescript-eslint/parser' // 指定解析器
  },
  plugins: ['@typescript-eslint'], // 启用 TypeScript 插件
  rules: {
    'prettier/prettier': [
      'warn',
      {
        singleQuote: true, // 单引号
        semi: false, // 无分号
        printWidth: 80, // 每行宽度至多80字符
        trailingComma: 'none', // 不加对象|数组最后逗号
        endOfLine: 'auto' // 换行符号不限制（win mac 不一致）
      }
    ],
    'vue/multi-word-component-names': [
      'warn',
      {
        ignores: ['index'] // vue组件名称多单词组成（忽略index.vue）
      }
    ],
    'vue/no-setup-props-destructure': ['off'], // 关闭 props 解构的校验
    'no-undef': 'error', // 添加未定义变量错误提示
    '@typescript-eslint/no-unused-vars': 'warn', // TypeScript 未使用变量警告
    '@typescript-eslint/no-explicit-any': 'off', // 允许使用 any 类型
    '@typescript-eslint/explicit-module-boundary-types': 'off' // 不强制导出函数的返回类型
  },
  globals: {
    ElMessage: 'readonly',
    ElMessageBox: 'readonly',
    ElLoading: 'readonly'
  }
}
