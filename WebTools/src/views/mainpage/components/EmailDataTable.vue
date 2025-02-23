<script setup  lang="ts">
import { Delete } from '@element-plus/icons-vue'

defineProps<{
  data: Record<string, any>[]
  columns: string[]
  loading: boolean
}>()

const emit = defineEmits(['delete-row', 'add-row'])
</script>

<template>
  <el-table class="table" v-loading="loading" border :data="data" height="500">
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
            size="normal"
            circle
            :icon="Delete"
            plain
            @click.prevent="emit('delete-row', $index)"
          />
          <el-button
            type="primary"
            size="normal"
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
</template>

<style lang="scss" scoped>
:deep(.el-table__header-wrapper .el-table_1_column_1) {
  text-align: center;
}
</style>
