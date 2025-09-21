<script setup>
import { onMounted, onBeforeUnmount, ref, watch, computed,defineProps,defineEmits  } from 'vue'
import { Color } from '@tiptap/extension-color'
import ListItem from '@tiptap/extension-list-item'
import TextStyle from '@tiptap/extension-text-style'
import Underline from '@tiptap/extension-underline'
import { CharacterCount } from '@tiptap/extension-character-count'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import { Editor, EditorContent } from '@tiptap/vue-3'
import SvgIcon from './Svglcon.vue'
import { uploadPhotoService } from '@/api/postUpload'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

// 创建编辑器实例
const editor = ref(null)
const limit = 2000 // 字数限制
const linkDialogVisible = ref(false) // 控制对话框显示
const linkUrl = ref('') // 存储链接地址

const percentage = computed(() => {
  return Math.round(
    (100 / limit) * editor.value.storage.characterCount.characters()
  )
})

onMounted(() => {
  editor.value = new Editor({
    extensions: [
      Color.configure({ types: [TextStyle.name, ListItem.name] }),
      TextStyle.configure({ types: [ListItem.name] }),
      StarterKit,
      Underline,
      Link.configure({
        openOnClick: false, // 禁用点击链接时打开
        autolink: true // 自动检测链接
      }), // 添加链接扩展
      CharacterCount.configure({
        limit: 2000 // 可选：设置字数限制
      }),
      Image.configure({
        inline: true, // 图片作为内联元素
        allowBase64: false // 不允许 base64 图片
      })
    ],
    content: props.modelValue,
    onUpdate: () => {
      // 当编辑器内容更新时，触发事件
      emit('update:modelValue', editor.value.getHTML())
    }
  })
})

// 监听父组件传递的值变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue !== editor.value.getHTML()) {
      editor.value.commands.setContent(newValue, false) // 更新编辑器内容
    }
  }
)

onBeforeUnmount(() => {
  if (editor.value) {
    editor.value.destroy()
  }
})

// 处理标题下拉列表的选项点击
const handleHeadingCommand = (level) => {
  editor.value.chain().focus().toggleHeading({ level }).run()
}

// 打开链接对话框
const openLinkDialog = () => {
  // 获取当前选中的链接地址（如果有）
  linkUrl.value = editor.value.getAttributes('link').href || ''
  linkDialogVisible.value = true
}

// 确认链接
const confirmLink = () => {
  if (linkUrl.value === '') {
    // 如果链接为空，移除链接
    editor.value.chain().focus().extendMarkRange('link').unsetLink().run()
  } else {
    // 更新或添加链接
    editor.value
      .chain()
      .focus()
      .extendMarkRange('link')
      .setLink({ href: linkUrl.value })
      .run()
  }
  linkDialogVisible.value = false
}

// 处理图片上传
const handleImageUpload = async (options) => {
  const { file } = options // 从 options 中解构出 file
  if (!file) return

  try {
    const formData = new FormData()
    formData.append('file', file)

    const response = await uploadPhotoService(formData)
    console.log('上传响应:', response) // 打印完整的响应

    // 检查 response 的结构
    if (
      !response ||
      !response.data ||
      !response.data.file
    ) {
      throw new Error('上传响应格式不正确')
    }

    const imageUrl = response.data.file
    console.log('图片 URL:', imageUrl)

    if (editor.value) {
      editor.value.chain().focus().setImage({ src: imageUrl }).run()
    } else {
      console.error('编辑器未初始化')
    }
  } catch (error) {
    console.error('图片上传失败:', error.message || error)
  }
}
</script>

<template>
  <div v-if="editor" class="container">
    <div class="control-group">
      <div class="button-group">
        <el-button class="colorinput">
          <input
            type="color"
            @input="editor.chain().focus().setColor($event.target.value).run()"
            :value="editor.getAttributes('textStyle').color || '#000000'"
          />
        </el-button>
        <el-tooltip content="加粗" placement="bottom">
          <el-button
            @click="editor.chain().focus().toggleBold().run()"
            :disabled="!editor.can().chain().focus().toggleBold().run()"
            :class="{ 'is-active': editor.isActive('bold') }"
          >
            <SvgIcon name="bold"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="斜体" placement="bottom">
          <el-button
            @click="editor.chain().focus().toggleItalic().run()"
            :disabled="!editor.can().chain().focus().toggleItalic().run()"
            :class="{ 'is-active': editor.isActive('italic') }"
          >
            <SvgIcon name="Text_Italic"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="标题" placement="top">
          <el-dropdown
            @command="handleHeadingCommand"
            :class="{ 'is-active': editor.isActive('heading') }"
          >
            <el-button>
              <SvgIcon name="H" color="#000"></SvgIcon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  divided
                  :command="1"
                  :class="{
                    'is-active': editor.isActive('heading', { level: 1 })
                  }"
                >
                  <SvgIcon name="h1" color="#000"></SvgIcon>
                </el-dropdown-item>
                <el-dropdown-item
                  :command="2"
                  divided
                  :class="{
                    'is-active': editor.isActive('heading', { level: 2 })
                  }"
                >
                  <SvgIcon name="h2" color="#000"></SvgIcon>
                </el-dropdown-item>
                <el-dropdown-item
                  :command="3"
                  divided
                  :class="{
                    'is-active': editor.isActive('heading', { level: 3 })
                  }"
                >
                  <SvgIcon name="h3" color="#000"></SvgIcon>
                </el-dropdown-item>
                <el-dropdown-item
                  :command="4"
                  divided
                  :class="{
                    'is-active': editor.isActive('heading', { level: 4 })
                  }"
                >
                  <SvgIcon name="h4" color="#000"></SvgIcon>
                </el-dropdown-item>
                <el-dropdown-item
                  :command="5"
                  divided
                  :class="{
                    'is-active': editor.isActive('heading', { level: 5 })
                  }"
                >
                  <SvgIcon name="h5" color="#000"></SvgIcon>
                </el-dropdown-item>
                <el-dropdown-item
                  :command="6"
                  divided
                  :class="{
                    'is-active': editor.isActive('heading', { level: 6 })
                  }"
                >
                  <SvgIcon name="h6" color="#000"></SvgIcon>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-tooltip>
        <el-tooltip content="无序列表" placement="bottom">
          <el-button
            @click="editor.chain().focus().toggleBulletList().run()"
            :class="{ 'is-active': editor.isActive('bulletList') }"
          >
            <SvgIcon name="Unorderedlistings"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="有序列表" placement="bottom">
          <el-button
            @click="editor.chain().focus().toggleOrderedList().run()"
            :class="{ 'is-active': editor.isActive('orderedList') }"
          >
            <SvgIcon name="Orderedlists" color="#000"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="分界线" placement="bottom">
          <el-button @click="editor.chain().focus().setHorizontalRule().run()">
            <SvgIcon name="boundary" color="#000"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="下划线" placement="bottom">
          <el-button
            @click="editor.chain().focus().toggleUnderline().run()"
            :disabled="!editor.can().chain().focus().toggleUnderline().run()"
            :class="{ 'is-active': editor.isActive('underline') }"
          >
            <SvgIcon name="underline" color="#000"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="图片上传" placement="bottom">
          <el-upload
            action="#"
            :show-file-list="false"
            :before-upload="handleImageUpload"
            :http-request="handleImageUpload"
          >
            <el-button>
              <SvgIcon name="image" color="#000"></SvgIcon>
            </el-button>
          </el-upload>
        </el-tooltip>
        <el-tooltip content="添加链接" placement="bottom">
          <el-button
            @click="openLinkDialog"
            :class="{ 'is-active': editor.isActive('link') }"
          >
            <SvgIcon name="link" color="#000"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="撤消" placement="bottom">
          <el-button
            @click="editor.chain().focus().undo().run()"
            :disabled="!editor.can().chain().focus().undo().run()"
          >
            <SvgIcon name="undo" color="#000"></SvgIcon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="重做" placement="bottom">
          <el-button
            @click="editor.chain().focus().redo().run()"
            :disabled="!editor.can().chain().focus().redo().run()"
          >
            <SvgIcon name="redo" color="#000"></SvgIcon>
          </el-button>
        </el-tooltip>
      </div>
    </div>
    <editor-content :editor="editor" />
    <div
      :class="{
        'character-count': true,
        'character-count--warning':
          editor.storage.characterCount.characters() === limit
      }"
    >
      <svg height="20" width="20" viewBox="0 0 20 20">
        <circle r="10" cx="10" cy="10" fill="#e9ecef" />
        <circle
          r="5"
          cx="10"
          cy="10"
          fill="transparent"
          stroke="currentColor"
          stroke-width="10"
          :stroke-dasharray="`calc(${percentage} * 31.4 / 100) 31.4`"
          transform="rotate(-90) translate(-20)"
        />
        <circle r="6" cx="10" cy="10" fill="white" />
      </svg>
      {{ editor.storage.characterCount.characters() }} / {{ limit }} 个字符
    </div>
    <!-- 改tiptap原本的添加链接框 -->
    <el-dialog
      v-if="linkDialogVisible"
      v-model="linkDialogVisible"
      title="添加链接"
      width="30%"
    >
      <el-input v-model="linkUrl" placeholder="请输入链接地址"></el-input>
      <template #footer>
        <el-button @click="linkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmLink">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss">
/* 容器样式 */
.container {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  background: #ffffff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

/* 工具栏样式 */
.control-group {
  background: #f8f9fa;
  padding: 8px;
  border-bottom: 1px solid #e0e0e0;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

/* 按钮样式 */
button {
  background: #ffffff;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition:
    background 0.2s,
    border-color 0.2s,
    color 0.2s;

  &:hover {
    background: #f3f4f6;
    border-color: #9ca3af;
  }

  &:disabled {
    background: #f3f4f6;
    color: #9ca3af;
    cursor: not-allowed;
  }

  &.is-active {
    background: #3b82f6;
    border-color: #3b82f6;
    color: #ffffff;

    &:hover {
      background: #2563eb;
      border-color: #2563eb;
    }
  }
}

/* 编辑器内容区域样式 */
.tiptap {
  padding: 16px;
  min-height: 400px;
  outline: none;

  /* 段落样式 */
  p {
    margin: 0 0 1rem 0;
    line-height: 1.6;
    color: #374151;
  }

  /* 标题样式 */
  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    margin: 1.5rem 0 1rem 0;
    line-height: 1.3;
    color: #111827;
  }

  h1 {
    font-size: 2rem;
    font-weight: 700;
  }

  h2 {
    font-size: 1.75rem;
    font-weight: 600;
  }

  h3 {
    font-size: 1.5rem;
    font-weight: 600;
  }

  h4 {
    font-size: 1.25rem;
    font-weight: 600;
  }

  h5 {
    font-size: 1rem;
    font-weight: 600;
  }

  h6 {
    font-size: 0.875rem;
    font-weight: 600;
  }

  /* 列表样式 */
  ul,
  ol {
    padding-left: 1.5rem;
    margin: 1rem 0;

    li {
      margin: 0.5rem 0;
    }
  }

  ul {
    list-style-type: disc;
  }

  ol {
    list-style-type: decimal;
  }

  /* 代码块样式 */
  pre {
    background: #1e293b;
    color: #f8fafc;
    padding: 1rem;
    border-radius: 6px;
    overflow-x: auto;
    margin: 1rem 0;

    code {
      background: none;
      color: inherit;
      font-size: 0.875rem;
      font-family: 'JetBrains Mono', monospace;
    }
  }

  /* 行内代码样式 */
  code {
    background: #f3f4f6;
    color: #dc2626;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.875rem;
    font-family: 'JetBrains Mono', monospace;
  }

  /* 引用样式 */
  blockquote {
    border-left: 4px solid #3b82f6;
    margin: 1rem 0;
    padding: 0.5rem 1rem;
    background: #f8f9fa;
    color: #374151;
    font-style: italic;

    p {
      margin: 0;
    }
  }

  /* 水平分割线样式 */
  hr {
    border: none;
    border-top: 1px solid #e0e0e0;
    margin: 1.5rem 0;
  }

  /* 图片样式 */
  img {
    max-width: 100%;
    height: auto;
    border-radius: 6px;
    margin: 1rem 0;
  }
}

/* 激活标题按钮时的样式 */
.active-heading {
  color: #3b82f6; /* 蓝色 */
  font-weight: 600; /* 加粗 */
}

.el-dropdown-menu {
  .el-dropdown-menu__item--divided {
    margin: 0 !important;
  }
}

.button-group {
  .el-button {
    margin: 0 !important;
  }
}

/* 确保下划线可见 */
.tiptap u {
  text-decoration: underline;
}

.tiptap [style*='text-decoration: underline'] {
  text-decoration: underline;
}

/* 字数统计样式 */
.character-count {
  align-items: center;
  color: #666; // 默认文字颜色
  display: flex;
  font-size: 0.75rem;
  gap: 0.5rem;
  margin: 1.5rem;

  svg {
    color: #6f42c1; // 进度条默认颜色（紫色）
  }

  &--warning,
  &--warning svg {
    color: #dc3545; // 超出字数限制时的颜色（红色）
  }
}

/* 进度条样式 */
svg {
  height: 20px;
  width: 20px;
  viewbox: 0 0 20 20;

  circle {
    r: 10;
    cx: 10;
    cy: 10;

    /* 背景圆 */
    &:first-child {
      fill: #e9ecef; // 背景颜色
    }

    /* 进度条圆 */
    &:nth-child(2) {
      fill: transparent;
      stroke: currentColor; // 进度条颜色
      stroke-width: 10;
      transform: rotate(-90) translate(-20);
    }

    /* 中心圆 */
    &:last-child {
      r: 6;
      fill: white; // 中心填充颜色
    }
  }
}
.colorinput.el-button {
  padding: 0 !important;
}
</style>
