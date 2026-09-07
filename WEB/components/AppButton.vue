<template>
  <el-button
    v-bind="$attrs"
    :type="type"
    :size="size"
    :loading="loading"
    :disabled="disabled || loading"
    :native-type="nativeType"
    :circle="circle"
    class="app-button"
    :class="{ 'app-button--block': block, 'app-button--circle': circle }"
    @click="$emit('click')"
  >
    <el-icon v-if="icon"><component :is="icon" /></el-icon>
    <slot v-if="!circle" />
  </el-button>
</template>

<script setup>
defineProps({
  type: { type: String, default: 'primary' },
  size: { type: String, default: 'large' },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  nativeType: { type: String, default: 'button' },
  block: { type: Boolean, default: false },
  icon: { type: [String, Object], default: null },
  circle: { type: Boolean, default: false }, 
})
defineEmits(['click'])
</script>

<style scoped>
.app-button {
  height: 44px;
  font-size: 15px;
  font-weight: 600;
}
.app-button--block {
  width: 100%;
}
.app-button--circle.el-button--small {
  width: 32px;
  height: 32px;
}
.app-button--circle.el-button--large {
  width: 50px;
  height: 50px;
}
.app-button--circle {
  width: 44px;
  height: 44px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>