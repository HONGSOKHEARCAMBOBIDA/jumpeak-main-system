<template>
  <el-form-item :label="label" :prop="prop">
    <el-input
      :model-value="modelValue"
      @update:model-value="handleInput"
      @keydown.enter="$emit('enter')"
      :type="type"
      :placeholder="placeholder"
      :prefix-icon="prefixIcon"
      :show-password="type === 'password'"
      size="large"
      :disabled="disabled"
      :clearable="clearable"
      :autofocus="autofocus"
      :readonly="readonly"
    />
  </el-form-item>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  prop: { type: String, default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  prefixIcon: { type: [String, Object], default: null },
  disabled: { type: Boolean, default: false },
  clearable: { type: Boolean, default: false },
  autofocus: { type: Boolean, default: false },
  readonly: { type: Boolean, default: false },
  trim: { type: Boolean, default: false } // true = trim spaces, false = keep as typed
})

const emit = defineEmits(['update:modelValue', 'enter'])

function handleInput(value) {
  emit('update:modelValue', props.trim ? value.trim() : value)
}
</script>