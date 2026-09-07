<script setup>
defineProps({
  modelValue: { type: [String, Number, Array], default: '' },
  label: { type: String, default: '' },
  prop: { type: String, default: '' },
  placeholder: { type: String, default: 'Select' },
  options: { type: Array, default: () => [] },
  multiple: { type: Boolean, default: false },
  clearable: { type: Boolean, default: true },
  filterable: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue', 'change'])

function normalize(opt) {
  return typeof opt === 'object' ? opt : { label: opt, value: opt }
}
</script>

<template>
  <el-form-item :prop="prop" :label="label" class="app-select mb-3">
    <el-select
      :model-value="modelValue"
      :placeholder="placeholder"
      :multiple="multiple"
      :clearable="clearable"
      :filterable="filterable"
      :disabled="disabled"
      class="w-full"
      @update:model-value="$emit('update:modelValue', $event)"
      @change="$emit('change', $event)"
    >
      <el-option
        v-for="opt in options.map(normalize)"
        :key="opt.value"
        :label="opt.label"
        :value="opt.value"
      />
    </el-select>
  </el-form-item>
</template>