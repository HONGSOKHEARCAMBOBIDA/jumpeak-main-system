<script setup>
defineProps({
  description: { type: String, default: 'No data available' },
  showAction: { type: Boolean, default: false },
  actionText: { type: String, default: 'Refresh' },
})
defineEmits(['action'])
</script>

<template>
  <div class="flex flex-col items-center justify-center py-14 px-4 text-center">
    <el-icon :size="52" class="mb-3" style="color: var(--color-text-placeholder)">
      <DocumentRemove />
    </el-icon>
    <p class="text-sm" style="color: var(--color-text-secondary)">{{ description }}</p>
    <el-button v-if="showAction" class="mt-4" size="small" @click="$emit('action')">
      {{ actionText }}
    </el-button>
  </div>
</template>
