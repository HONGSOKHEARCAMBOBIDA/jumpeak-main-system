<template>
  <router-view />
</template>

<script setup>
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Siemreap&display=swap');
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Siemreap', sans-serif;
  background: #f0f2f5;
}

.el-button,
.el-input,
.el-select,
.el-dialog,
.el-table,
.el-menu,
.el-form,
.el-card,
.el-dropdown,
.el-date-picker,
.el-input__inner::placeholder {
  font-family: 'Siemreap', sans-serif;
}
.el-pagination {
  font-family: 'Siemreap', sans-serif !important;
}
.el-loading-mask {
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

:root {
  --primary: #083767;
  --success: #67c23a;
  --warning: #e6a23c;
  --danger: #f56c6c;
  
}
</style>
