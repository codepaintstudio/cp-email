<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  Warning,
  Download,
  User,
  Upload,
  Edit,
  Promotion,
  QuestionFilled
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: true
  }
})

interface Step {
  icon: any
  title: string
  content: string
  subSteps?: string[]
  demo?: string
}

interface ImportantNote {
  type: 'error' | 'warning' | 'success' | 'info'
  title: string
  details: string[]
}

const emit = defineEmits(['update:modelValue'])

// 响应式状态
const hideForever = ref<boolean>(localStorage.getItem('hideIntro') === 'true')
const visible = ref<boolean>(props.modelValue)

// 步骤配置数据
const steps = ref<Step[]>([
  {
    icon: Download,
    title: '第一步：下载模板',
    content: '从后端获取标准Excel模板文件,确保数据格式正确',
    subSteps: ['点击右上角 <b>下载模板</b> 按钮', '模板第一列必须为 <em>email</em>'],
    demo: `   email
   user1@demo.com
   user2@demo.com`
  },
  {
    icon: User,
    title: '第二步：邮箱登录',
    content: '配置发送邮箱的SMTP授权信息',
    subSteps: [
      '进入邮箱设置 → 账户 → 开启SMTP服务',
      '获取16位授权码（非邮箱密码）',
      '在登录界面输入邮箱地址和授权码'
    ],
    demo: '邮箱地址：yourname@example.com\n授权码：ABC123DEF456GHIJ'
  },
  {
    icon: Upload,
    title: '第三步：上传数据',
    content: '上传包含收件人信息的Excel文件',
    subSteps: ['支持 .xlsx 和 .xls 格式', '最多同时处理100条数据', '自动识别首列作为收件地址']
  },
  {
    icon: Edit,
    title: '第四步：撰写邮件',
    content: '使用富文本编辑器编写邮件内容',
    subSteps: ['在"邮件主题"输入框中填写主题', '支持图片和链接发送'],
    demo: `尊敬的xx:
  您已成功申请xx职位,请查收面试通知。`
  },
  {
    icon: Promotion,
    title: '第五步：批量发送',
    content: '安全高效地完成邮件投递',
    subSteps: ['点击"一键发送"按钮开始发送', '显示发送进度和状态'],
    demo: `发送状态：
✅ 成功
⏳ 未发送
❌ 失败`
  }
])

// 注意事项配置
const importantNotes = ref<ImportantNote[]>([
  {
    type: 'error',
    title: '安全须知',
    details: ['授权码请妥善保管', '建议使用专用邮箱进行发送']
  },
  {
    type: 'warning',
    title: '发送策略',
    details: ['单次发送建议不超过100封邮件', '避免重复发送相同内容']
  },
  {
    type: 'success',
    title: '最佳实践',
    details: ['正式发送前进行测试邮件验证', '适当清空发送数据']
  },
  {
    type: 'info',
    title: '技术支持',
    details: [
      '如有问题，请联系技术支持团队',
      '邮箱:wujieruanchuang@163.com',
      'GitHub:https://github.com/CodePaintStudio/CP-EmailTools'
    ]
  }
])

// 监听prop变化
watch(
  () => props.modelValue,
  (newVal) => {
    visible.value = newVal
  }
)

// 处理关闭事件
const handleClose = (): void => {
  localStorage.setItem('hideIntro', hideForever.value.toString())
  visible.value = false
  emit('update:modelValue', false)
}
</script>

<template>
  <el-drawer
    v-model="visible"
    title="📧 CP-EmailTools 使用指南"
    size="100%"
    direction="ttb"
    :before-close="handleClose"
    :close-on-click-modal="false"
    class="intro-drawer"
  >
    <div class="intro-content">
      <!-- 可视化演示区域 -->
      <!-- <div class="demo-section">
         <img src="@/assets/email-demo.gif" alt="操作演示" class="demo-gif" /> -->
      <!-- <p class="demo-tip">📌 实时预览 - 从上传到发送的全流程演示</p> -->
      <!-- </div>  -->

      <!-- 步骤时间线 -->
      <div class="steps-section">
        <el-timeline>
          <el-timeline-item v-for="(step, index) in steps" :key="index" placement="top">
            <h3 class="step-title">{{ step.title }}</h3>
            <div class="step-content">
              <p class="step-description">{{ step.content }}</p>
              <ol class="sub-steps" v-if="step.subSteps">
                <li
                  v-for="(subStep, subIndex) in step.subSteps"
                  :key="subIndex"
                  v-html="subStep"
                ></li>
              </ol>
              <div v-if="step.demo" class="code-demo">
                <pre>{{ step.demo }}</pre>
              </div>
            </div>
          </el-timeline-item>
        </el-timeline>
      </div>

      <!-- 注意事项 -->
      <div class="notice-section">
        <h3 class="notice-title">
          <el-icon><Warning /></el-icon> 重要提醒
        </h3>
        <div class="notice-grid">
          <el-alert
            v-for="(note, index) in importantNotes"
            :key="index"
            :title="note.title"
            :type="note.type"
            :closable="false"
            show-icon
            class="notice-card"
          >
            <ul class="notice-list">
              <li v-for="(detail, detailIndex) in note.details" :key="detailIndex">
                {{ detail }}
              </li>
            </ul>
          </el-alert>
        </div>
      </div>
    </div>

    <!-- 底部操作 -->
    <template #footer>
      <div class="dialog-footer">
        <el-checkbox v-model="hideForever" class="never-show-check">
          <span>不再显示本指南</span>
          <el-tooltip content="可在页面右上角「使用说明」重新打开" placement="top">
            <el-icon class="tip-icon"><QuestionFilled /></el-icon>
          </el-tooltip>
        </el-checkbox>
        <div class="action-buttons">
          <el-button type="primary" @click="handleClose" class="start-button">
            <el-icon><Promotion /></el-icon>
            立即开始
          </el-button>
          <el-button plain @click="handleClose" class="later-button"> 稍后再说 </el-button>
        </div>
      </div>
    </template>
  </el-drawer>
</template>

<style lang="scss" scoped>
.intro-drawer {
  --step-icon-size: 28px;
  --step-color: #409eff;
  --notice-bg: #f8f9fa;

  .intro-content {
    padding: 20px;
    max-height: 80vh;
    overflow-y: auto;
  }

  .demo-section {
    margin-bottom: 30px;
    text-align: center;

    .demo-gif {
      width: 100%;
      max-width: 800px;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    .demo-tip {
      color: #666;
      font-size: 0.9em;
      margin-top: 10px;
    }
  }

  .steps-section {
    margin: 20px 0;

    :deep(.el-timeline) {
      padding-left: 40px;
    }

    .step-icon {
      width: var(--step-icon-size);
      height: var(--step-icon-size);
      color: white;
      background: var(--step-color);
      border-radius: 50%;
      padding: 6px;
    }

    .step-title {
      color: var(--step-color);
      font-size: 1.3em;
      margin: 10px 0;
    }

    .sub-steps {
      padding-left: 25px;
      color: #666;

      li {
        margin: 8px 0;
        line-height: 1.6;

        &::marker {
          color: var(--step-color);
        }
      }
    }

    .code-demo {
      background: var(--notice-bg);
      border-radius: 6px;
      padding: 12px;
      margin: 15px 0;
      pre {
        margin: 0;
        white-space: pre-wrap;
        font-family: Consolas, monospace;
      }
    }
  }

  .notice-section {
    margin-top: 30px;

    .notice-title {
      color: #e6a23c;
      font-size: 1.2em;
      margin-bottom: 15px;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .notice-grid {
      display: grid;
      gap: 15px;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    }

    .notice-card {
      transition: transform 0.3s;

      &:hover {
        transform: translateY(-2px);
      }
    }

    .notice-list {
      padding-left: 20px;
      li {
        list-style: disc;
        margin: 6px 0;
      }
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 20px;
    background: var(--notice-bg);
    border-top: 1px solid #eee;

    .never-show-check {
      .tip-icon {
        margin-left: 8px;
        color: #909399;
        cursor: help;
      }
    }

    .action-buttons {
      display: flex;
      gap: 12px;

      .start-button {
        padding: 10px 25px;
      }

      .later-button {
        &:hover {
          color: var(--step-color);
          border-color: var(--step-color);
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .intro-drawer {
    .intro-content {
      padding: 10px;
      max-height: 70vh;
    }

    .demo-gif {
      border-radius: 4px;
    }

    .steps-section {
      :deep(.el-timeline) {
        padding-left: 20px;
      }
    }

    .dialog-footer {
      flex-direction: column;
      gap: 15px;

      .action-buttons {
        width: 100%;
        button {
          flex: 1;
        }
      }
    }
  }
}
</style>
