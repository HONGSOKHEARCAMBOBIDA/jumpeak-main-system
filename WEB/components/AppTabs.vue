<template>
  <el-tabs
    v-model="activeTab"
    :type="type"
    :tab-position="tabPosition"
    :stretch="stretch"
    @tab-change="onTabChange"
  >
    <el-tab-pane
      v-for="tab in tabs"
      :key="tab.name"
      :label="tab.label"
      :name="tab.name"
      :disabled="tab.disabled"
      :lazy="lazy"
    >
      <slot :name="tab.name">
        {{ tab.content }}
      </slot>
    </el-tab-pane>
  </el-tabs>
</template>

<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  modelValue: { type: String, default: "" },
  tabs: { type: Array, required: true },
  type: { type: String, default: "" },
  tabPosition: { type: String, default: "top" },
  stretch: { type: Boolean, default: false },
  lazy: { type: Boolean, default: true },
});

const emit = defineEmits(["update:modelValue", "change"]);

const activeTab = ref(props.modelValue || props.tabs?.[0]?.name || "");

watch(
  () => props.modelValue,
  (val) => {
    if (val !== undefined && val !== activeTab.value) {
      activeTab.value = val;
    }
  },
);

watch(activeTab, (val) => {
  emit("update:modelValue", val);
});

function onTabChange(name) {
  emit("change", name);
}
</script>

