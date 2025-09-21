<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  data: Record<string, any>[]
  columns: string[]
  loading: boolean
}>()

const emit = defineEmits(['delete-row', 'add-row'])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(30)

// 计算当前页的数据
const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return props.data.slice(start, end)
})

// 处理页码改变
const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

watch(
  () => props.data,
  (newData) => {
    const totalPages = Math.ceil(newData.length / pageSize.value)
    if (currentPage.value > totalPages) {
      currentPage.value = 1
    }
  },
)
</script>

<template>
  <div class="table-container">
  <el-table class="table" v-loading="loading" border :data="paginatedData" height="500">
    <el-table-column label="状态" width="85">
      <template v-slot:default="{ row }">
        <el-button
          :type="
            row.state === 0 ? 'danger' : row.state === 1 ? 'warning' : 'success'
          "
          size="small"
          class="state_bt"
          disabled
        >
          {{ row.state === 0 ? '失败' : row.state === 1 ? '未发送' : '成功' }}
        </el-button>
      </template>
    </el-table-column>

    <el-table-column
      v-for="key in columns"
      :key="key"
      :prop="key"
      :label="key"
      min-width="200"
    >
      <template v-slot="{ row }">
        <el-input v-model="row[key]"></el-input>
      </template>
    </el-table-column>

    <el-table-column fixed="right" label="操作" width="100">
      <template #default="{ $index }">
        <div class="operation">
          <el-button
            type="danger"
            size="default"
            circle
            :icon="Delete"
            plain
            @click.prevent="emit('delete-row', $index)"
          />
          <el-button
            type="primary"
            size="default"
            circle
            class="add-row-btn"
            @click="emit('add-row')"
          >
            +
          </el-button>
        </div>
      </template>
    </el-table-column>
  </el-table>
  <!-- 添加分页组件 -->
    <div class="pagination-container">
      <el-pagination
        v-if="data.length > pageSize"
        background
        layout="prev, pager, next"
        :total="data.length"
        :page-size="pageSize"
        :current-page="currentPage"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
:deep(.el-table__header-wrapper .el-table_1_column_1) {
  text-align: center;
}

.table-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}
</style>