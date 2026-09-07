<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="title"
    :width="isMobile ? '100%' : width"
    :fullscreen="isMobile"
    :class="['app-dialog', isMobile ? 'app-dialog--mobile' : 'app-dialog--desktop', customClass]"
    :destroy-on-close="destroyOnClose"
    :close-on-click-modal="!isMobile && closeOnClickModal"
    :show-close="!isMobile"
    :align-center="!isMobile"
    @closed="$emit('closed')"
  >
    <template v-if="isMobile" #header>
      <div class="app-dialog__mobile-header">
        <span class="app-dialog__mobile-title">{{ title }}</span>
        <el-icon class="app-dialog__mobile-close" @click="$emit('update:modelValue', false)">
          <Close />
        </el-icon>
      </div>
    </template>

    <div class="app-dialog__body">
      <slot />
    </div>

    <template v-if="$slots.footer" #footer>
      <div :class="['app-dialog__footer', isMobile ? 'app-dialog__footer--mobile' : '']">
        <slot name="footer" />
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { Close } from "@element-plus/icons-vue";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: "" },
  width: { type: String, default: "640px" },
  destroyOnClose: { type: Boolean, default: true },
  closeOnClickModal: { type: Boolean, default: true },
  customClass: { type: String, default: "" },
  breakpoint: { type: Number, default: 768 },
});

defineEmits(["update:modelValue", "closed"]);

const isMobile = ref(false);

function checkIsMobile() {
  isMobile.value = window.innerWidth < props.breakpoint;
}

onMounted(() => {
  checkIsMobile();
  window.addEventListener("resize", checkIsMobile);
});

onUnmounted(() => {
  window.removeEventListener("resize", checkIsMobile);
});
</script>

<style scoped>
/* Desktop */
.app-dialog--desktop :deep(.el-dialog__header) {
  margin-right: 0;
  padding-bottom: 12px;
}

.app-dialog__body {
  max-height: 70vh;
  overflow-y: auto;
  padding-right: 4px;
}

.app-dialog__footer {
  display: flex;
  justify-content: flex-end;

}

/* Mobile: full screen, slide up, sticky footer */
.app-dialog--mobile :deep(.el-dialog) {
  display: flex;
  flex-direction: column;
  margin: 0;
  border-radius: 0;
}


.app-dialog--mobile :deep(.el-dialog__body) {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.app-dialog--mobile :deep(.el-dialog__footer) {
  padding: 12px 16px;
  border-top: 1px solid var(--el-border-color-lighter);
  position: sticky;
  bottom: 0;
  background: var(--el-bg-color);
}

.app-dialog__mobile-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.app-dialog__mobile-title {
  font-size: 16px;
  font-weight: 600;
}

.app-dialog__mobile-close {
  font-size: 18px;
  cursor: pointer;
  color: var(--el-text-color-secondary);
}

.app-dialog--mobile .app-dialog__body {
  max-height: none;
}

</style>