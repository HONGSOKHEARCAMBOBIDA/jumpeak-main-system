<template>
  <div
    class="app-table"
    :style="{ '--brand-color': companyStore.color || '#4589ce' }"
  >
    <el-table
      v-if="!isMobile"
      ref="tableRef"
      :data="data"
      v-loading="loading"
      stripe
      border
      v-bind="tableProps"
      @selection-change="onSelectionChange"
    >
      <el-table-column v-if="selectable" type="selection" width="50" />

      <el-table-column
        v-if="showIndex"
        type="index"
        :label="indexLabel"
        width="70"
      />

      <el-table-column
        v-for="col in columns"
        :key="col.prop || col.label"
        v-bind="col"
      >
        <template v-if="col.slot" #default="scope">
          <slot :name="col.slot" v-bind="scope" />
        </template>
      </el-table-column>

      <el-table-column
        v-if="$slots.actions"
        :label="actionsLabel"
        :width="actionsWidth"
      >
        <template #default="scope">
          <slot name="actions" v-bind="scope" />
        </template>
      </el-table-column>
    </el-table>
    <div v-else class="app-table-cards" v-loading="loading">
      <el-empty v-if="!data.length && !loading" description="គ្មានទិន្នន័យ" />
      <div
        v-for="(row, idx) in data"
        :key="idx"
        class="app-table-card"
        :class="{ 'is-selected': selectable && isSelected(row) }"
      >
        <div v-if="selectable" class="app-table-card-header">
          <el-checkbox
            :model-value="isSelected(row)"
            @change="(val) => toggleMobileSelection(row, val)"
          />
        </div>

        <div
          v-for="col in columns"
          :key="col.prop || col.label"
          class="app-table-card-row"
        >
          <span class="app-table-card-label">{{ col.label }}</span>
          <span class="app-table-card-value">
            <slot v-if="col.slot" :name="col.slot" :row="row" :$index="idx" />
            <template v-else>{{ row[col.prop] }}</template>
          </span>
        </div>
        <div v-if="$slots.actions" class="app-table-card-actions">
          <slot name="actions" :row="row" :$index="idx" />
        </div>
      </div>
    </div>

    <div v-if="showPagination" class="app-table-pagination">
      <el-pagination
        v-model:current-page="currentPageProxy"
        v-model:page-size="pageSizeProxy"
        :total="total"
        :layout="isMobile ? 'total, prev, next' : 'total, prev, pager, next'"
        @change="$emit('page-change')"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
// import { useCompanyStore } from '../stores/company'
import { useCompanyStore } from "../src/stores/company";
const companyStore = useCompanyStore();
const props = defineProps({
  data: { type: Array, default: () => [] },
  columns: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  showIndex: { type: Boolean, default: false },
  indexLabel: { type: String, default: "ល.រ" },
  actionsLabel: { type: String, default: "សកម្មភាព" },
  actionsWidth: { type: [String, Number], default: 190 },
  mobileBreakpoint: { type: Number, default: 768 },
  showPagination: { type: Boolean, default: true },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  total: { type: Number, default: 0 },
  tableProps: { type: Object, default: () => ({}) },
  selectable: { type: Boolean, default: false },
});

const emit = defineEmits([
  "update:currentPage",
  "update:pageSize",
  "page-change",
  "selection-change",
]);

const currentPageProxy = computed({
  get: () => props.currentPage,
  set: (val) => emit("update:currentPage", val),
});
const pageSizeProxy = computed({
  get: () => props.pageSize,
  set: (val) => emit("update:pageSize", val),
});

const windowWidth = ref(window.innerWidth);
function onResize() {
  windowWidth.value = window.innerWidth;
}
onMounted(() => {
  window.addEventListener("resize", onResize);
  companyStore.fetchColor();
});
onUnmounted(() => window.removeEventListener("resize", onResize));

const isMobile = computed(() => windowWidth.value <= props.mobileBreakpoint);
const selectedRows = ref([]);
function onSelectionChange(rows) {
  // fired by el-table on desktop
  selectedRows.value = rows;
  emit("selection-change", rows);
}
function isSelected(row) {
  return selectedRows.value.includes(row);
}
function toggleMobileSelection(row, checked) {
  if (checked) {
    selectedRows.value.push(row);
  } else {
    selectedRows.value = selectedRows.value.filter((r) => r !== row);
  }
  emit("selection-change", selectedRows.value);
}
watch(
  () => props.data,
  () => {
    selectedRows.value = [];
    emit("selection-change", []);
  },
);
</script>

<style scoped>
.app-table-cards {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

:deep(.el-table__header-wrapper th) {
  /* background-color: var(--brand-color,#4589ce); */
  color: #030303 !important;
}

.app-table-card {
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  padding: 12px 14px;
  position: relative;
  background: var(--el-bg-color);
}

.app-table-card-index {
  position: absolute;
  top: 10px;
  right: 12px;
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

.app-table-card-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  padding: 6px 0;
  border-bottom: 1px dashed var(--el-border-color-lighter);
  font-size: 13px;
}

.app-table-card-row:last-of-type {
  border-bottom: none;
}

.app-table-card-label {
  color: var(--el-text-color-secondary);
  flex-shrink: 0;
  font-weight: 500;
}

.app-table-card-value {
  text-align: right;
  word-break: break-word;
}

.app-table-card-actions {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: flex-end;
}

.app-table-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.app-table-card-header {
  margin-bottom: 8px;
}
.app-table-card.is-selected {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}
@media (max-width: 768px) {
  .app-table-pagination {
    justify-content: center;
  }
}
</style>
