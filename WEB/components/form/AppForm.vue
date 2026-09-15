<script setup>
import { ref } from 'vue'
import AppButton from '../AppButton.vue'
const props = defineProps({
  model: { type: Object, required: true },
  rules: { type: Object, default: () => ({}) },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  labelWidth: { type: [String, Number], default: '120px' },
  labelPosition: { type: String, default: 'top' },
  submitText: { type: String, default: 'រក្សាទុក' },
  resetText: { type: String, default: 'ចាកចេញ' },
  showActions: { type: Boolean, default: true },
})
const emit = defineEmits(['submit', 'reset'])

const formRef = ref(null)

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    emit('submit', props.model)
  } catch {
    // validation failed — el-form already highlights the offending fields
  }
}

function handleReset() {
  formRef.value?.resetFields()
  emit('reset')
}

defineExpose({ validate: () => formRef.value?.validate(), resetFields: handleReset, formRef })
</script>

<template>
  <el-form
    ref="formRef"
    :model="model"
    :rules="rules"
    :label-width="labelWidth"
    :label-position="labelPosition"
    :disabled="disabled || loading"
    @submit.prevent
  >
    <slot />

    <div v-if="showActions" class="flex items-center gap-2 mt-2">
      <app-button  type="primary" :loading="loading" @click="handleSubmit" size="large">{{ submitText }}</app-button>
       <app-button  type="warning" :disabled="loading" @click="handleReset" size="large">{{ resetText }}</app-button>
    </div>
  </el-form>
</template>
