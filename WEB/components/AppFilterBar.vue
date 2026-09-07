<template>
  <el-card class="filter-card" shadow="never">
    <el-row :gutter="12" class="filter-row">
      <el-col
        v-for="(field, i) in normalizedFields"
        :key="i"
        :xs="24"
        :sm="12"
        :md="field.span || 6"
        class="filter-col"
      >
        <slot :name="field.slot" />
      </el-col>
      <el-col :xs="24" :sm="24" :md="actionSpan" class="filter-actions">
        <slot name="actions" />
      </el-col>
    </el-row>
  </el-card>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  fields: { type: Array, default: () => [] }, 
  actionSpan: { type: Number, default: 3 },
})

const normalizedFields = computed(() => props.fields)
</script>

<style scoped>
.filter-card {
  margin-bottom: 16px;
}

.filter-row {
  row-gap: 12px;
}

.filter-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .filter-actions {
    justify-content: stretch;
  }
  .filter-actions :deep(.el-button) {
    width: 100%;
  }
}
</style>