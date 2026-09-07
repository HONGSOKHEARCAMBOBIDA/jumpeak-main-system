<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import {
  getrole,
  getrolehaspermission,
  deleterolehaspermission,
  addrolehaspermission,
  editrole,
} from "../api/services.js";

const role = ref([]);
const rolehaspermission = ref([]);
const notify = useNotification();
const loading = ref(false);
const submitting = ref(false);
const roleFormRef = ref();
const dialogVisible = ref(false);
const roledialogvisible = ref(false);
const isAddMode = ref(false);
const selectroleID = ref(null);
const total = ref(0);
const searchKeyword = ref("");

// reactive object សម្រាប់ form (មិនមែន factory function ទេ)
const roleForm = reactive({
  id: null,
  name: "",
  display_name: "",
});

// Validation rules
const roleFormRules = {
  name: [{ required: true, message: "សូមបញ្ចូលឈ្មោះតួនាទី", trigger: "blur" }],
  display_name: [
    { required: true, message: "សូមបញ្ចូលឈ្មោះបង្ហាញ", trigger: "blur" },
  ],
};

// ត្រងបញ្ជី role តាមអក្សរដែល User វាយស្វែងរក
const filteredRoles = computed(() => {
  if (!searchKeyword.value) return role.value;
  const kw = searchKeyword.value.toLowerCase();
  return role.value.filter(
    (r) =>
      r.name?.toLowerCase().includes(kw) ||
      r.display_name?.toLowerCase().includes(kw)
  );
});

async function fetchrole() {
  loading.value = true;
  try {
    const res = await getrole();
    role.value = res.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការទាញយកទិន្នន័យ");
  } finally {
    loading.value = false;
  }
}

async function fetchrolehaspermission() {
  loading.value = true;
  try {
    const res = await getrolehaspermission(selectroleID.value);
    rolehaspermission.value = res.data.data || [];
    total.value = rolehaspermission.value.length;
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការទាញយកទិន្នន័យ");
  } finally {
    loading.value = false;
  }
}

// ចុចលើប៊ូតុងគ្រប់គ្រងសិទ្ធិ -> បើក Dialog ហើយទាញ permission របស់ role នោះ
function openPermissionDialog(row) {
  selectroleID.value = row.id;
  isAddMode.value = true;
  dialogVisible.value = true;
  fetchrolehaspermission();
}

// ចុចប៊ូតុងកែប្រែ -> បំពេញទិន្នន័យចាស់ចូល Form
function openRoleDialog(row) {
  selectroleID.value = row.id;
  roleForm.id = row.id;
  roleForm.name = row.name;
  roleForm.display_name = row.display_name;
  roledialogvisible.value = true;
}

function closeRoleDialog() {
  roledialogvisible.value = false;
  roleFormRef.value?.clearValidate();
}

// នៅពេលបិទ/បើក switch លើ permission ណាមួយ
async function handleToggleAssigned(row) {
  const payload = {
    role_id: selectroleID.value,
    permission_ids: [row.id],
  };

  submitting.value = true;
  try {
    if (row.assigned) {
      await addrolehaspermission(payload);
      notify.success("បន្ថែមសិទ្ធិដោយជោគជ័យ");
    } else {
      await deleterolehaspermission(payload);
      notify.success("ដកសិទ្ធិដោយជោគជ័យ");
    }
  } catch (e) {
    row.assigned = !row.assigned;
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការធ្វើបច្ចុប្បន្នភាព");
  } finally {
    submitting.value = false;
  }
}

async function updaterole() {
  const valid = await roleFormRef.value.validate().catch(() => false);
  if (!valid) return;

  const payload = {
    name: roleForm.name,
    display_name: roleForm.display_name,
  };

  submitting.value = true;
  try {
    await editrole(roleForm.id, payload);
    notify.success("កែប្រែដោយជោគជ័យ");
    roledialogvisible.value = false;
    fetchrole(); // ទាញទិន្នន័យតារាងឡើងវិញ ដើម្បីបង្ហាញការផ្លាស់ប្តូរថ្មី
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការធ្វើបច្ចុប្បន្នភាព");
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  fetchrole();
});
</script>

<template>
  <div class="role-page">
    <!-- ក្បាលទំព័រ + កម្មវិធីស្វែងរក -->
    <div class="page-header">
      <h2 class="page-title">គ្រប់គ្រងតួនាទី</h2>
      <AppInput
        v-model="searchKeyword"
        placeholder="ស្វែងរកតួនាទី..."
        clearable
        prefix-icon="Search"
        style="width: 260px"
      />
    </div>

    <AppTable
      :data="filteredRoles"
      :loading="loading"
      show-index
      :columns="[
        { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
        { prop: 'display_name', label: 'ឈ្មោះបង្ហាញ', minWidth: 120 },
      ]"
    >
      <template #actions="{ row }">
        <el-tooltip content="កែប្រែតួនាទី" placement="top">
          <AppButton
            size="small"
            type="warning"
            icon="Edit"
            circle
            @click="openRoleDialog(row)"
          />
        </el-tooltip>
        <el-tooltip content="គ្រប់គ្រងសិទ្ធិ" placement="top">
          <AppButton
            size="small"
            type="info"
            icon="Key"
            circle
            @click="openPermissionDialog(row)"
          />
        </el-tooltip>
      </template>
    </AppTable>

    <!-- Dialog គ្រប់គ្រងសិទ្ធិ -->
    <AppDialog
      v-model="dialogVisible"
      title="គ្រប់គ្រងសិទ្ធិ"
      width="700px"
    >
      <p class="dialog-subtitle">
        បិទ/បើក ដើម្បីផ្តល់ ឬដកសិទ្ធិចេញពីតួនាទីនេះ
      </p>
      <AppTable
        :data="rolehaspermission"
        :loading="loading"
        show-index
        :total="total"
        empty-text="មិនមានសិទ្ធិណាមួយទេ"
        :columns="[
          { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
          { prop: 'display_name', label: 'ឈ្មោះបង្ហាញ', minWidth: 120 },
          { label: 'សិទ្ធិ', slot: 'assigned', minWidth: 80, align: 'center' },
        ]"
      >
        <template #assigned="{ row }">
          <el-switch
            v-model="row.assigned"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #dcdfe6"
            :loading="submitting"
            @change="handleToggleAssigned(row)"
          />
        </template>
      </AppTable>

      <template #footer>
        <AppButton @click="dialogVisible = false">បិទ</AppButton>
      </template>
    </AppDialog>

    <!-- Dialog កែប្រែតួនាទី -->
    <AppDialog
      v-model="roledialogvisible"
      title="កែប្រែតួនាទី"
      width="480px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleFormRules"
        label-position="top"
        @submit.prevent="updaterole"
      >
          <AppInput
            label="ឈ្មោះ"
            prop="name"
            v-model="roleForm.name"
            placeholder="បញ្ចូលឈ្មោះតួនាទី"
            clearable
            :disabled="submitting"
            @enter="updaterole"
          />

          <AppInput
            label="ឈ្មោះបង្ហាញ"
            prop="display_name"
            v-model="roleForm.display_name"
            placeholder="បញ្ចូលឈ្មោះបង្ហាញ"
            clearable
            @enter="updaterole"
            :disabled="submitting"
          />
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <AppButton @click="closeRoleDialog" :disabled="submitting" plain>
            បោះបង់
          </AppButton>
          <AppButton @click="updaterole" :loading="submitting" type="primary">
            រក្សាទុក
          </AppButton>
        </div>
      </template>
    </AppDialog>
  </div>
</template>

<style scoped>
.role-page {
  padding: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.dialog-subtitle {
  color: #909399;
  font-size: 13px;
  margin: 0 0 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>