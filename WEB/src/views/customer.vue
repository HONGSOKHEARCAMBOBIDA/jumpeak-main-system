<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getcustomer,
  addcustomer,
  updatecustomer,
} from "../api/services.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppForm from "../../components/form/AppForm.vue";

const notify = useNotification();
const userDataStore = useUserDataStore();

const customers = ref([]);
const loading = ref(false);
const useloading = useLoading();
const saving = ref(false);

const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogVisible = ref(false);
const isEdit = ref(false);
const editId = ref(null);
const formRef = ref();

// filters
const filters = reactive({
  name: "",
});

const StatusOption = [
  { label: "ACTIVE", value: "ACTIVE" },
  { label: "INACTIVE", value: "INACTIVE" },
  { label: "BLACKLISTED", value: "BLACKLISTED" },
];

const form = reactive({
  name: "",
  phone: "",
  address: "",
  notes: "",
  credit_limit: 0,
  credit_limit_enforced: false,
  status: "ACTIVE",
});

const rules = {
  name: [{ required: true, message: "Customer name is required" }],
  credit_limit: [{ required: true, message: "Credit limit is required" }],
};

const canAddCustomer = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Customer"),
);
const canEditCustomer = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "update.Customer"),
);

async function fetchCustomers() {
  loading.value = true;
  try {
    const res = await getcustomer({
      page: page.value,
      page_size: pageSize.value,
      name: filters.name || undefined,
    });
    customers.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function handleSearch() {
  page.value = 1;
  fetchCustomers();
}

function openCreate() {
  isEdit.value = false;
  editId.value = null;
  form.name = "";
  form.phone = "";
  form.address = "";
  form.notes = "";
  form.credit_limit = 0;
  form.credit_limit_enforced = false;
  form.status = "ACTIVE";
  dialogVisible.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, {
    name: row.name || "",
    phone: row.phone || "",
    address: row.address || "",
    notes: row.notes || "",
    credit_limit: row.credit_limit ?? 0,
    credit_limit_enforced: !!row.credit_limit_enforced,
    status: row.status || "ACTIVE",
  });
  dialogVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    if (isEdit.value) {
      // matches CustomerRequestUpdate
      const payload = {
        name: form.name,
        phone: form.phone || null,
        address: form.address || null,
        notes: form.notes || null,
        credit_limit: form.credit_limit,
        credit_limit_enforced: form.credit_limit_enforced,
        status: form.status,
      };
      await updatecustomer(editId.value, payload);
      notify.success("កែប្រែអតិថិជនបានជោគជ័យ");
    } else {
      // matches CustomerRequestCreate — no status/enforced field on create
      const payload = {
        name: form.name,
        phone: form.phone || null,
        address: form.address || null,
        notes: form.notes || null,
        credit_limit: form.credit_limit,
      };
      await addcustomer(payload);
      notify.success("បង្កើតអតិថិជនបានជោគជ័យ");
    }
    dialogVisible.value = false;
    fetchCustomers();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

onMounted(() => {
  fetchCustomers();
});
</script>

<template>
  <div>
    <div class="page-header">
      <AppInput
        v-model="filters.name"
        placeholder="ស្វែងរកតាមឈ្មោះ"
        clearable
        style="width: 240px; margin-right: 12px"
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      />
      <AppButton size="default" @click="handleSearch">ស្វែងរក</AppButton>
      <AppButton
        v-if="canAddCustomer"
        type="primary"
        @click="openCreate"
        :block="false"
        size="default"
      >
        បន្ថែមអតិថិជន
      </AppButton>
    </div>

    <el-card class="table-card">
      <AppTable
        :data="customers"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchCustomers"
        :columns="[
          { slot: 'name', label: 'ឈ្មោះ', minWidth: 150 },
          { prop: 'phone', label: 'ទូរស័ព្ទ', width: 130 },
          { prop: 'address', label: 'អាស័យដ្ឋាន', width: 130 },
          { prop: 'company_name', label: 'ក្រុមហ៊ុន', minWidth: 120 },
          { slot: 'branch_name', label: 'សាខា', minWidth: 120 },
          { slot: 'credit_limit', label: 'កម្រិតឥណទាន', width: 130 },
          { slot: 'current_outstanding', label: 'ជំពាក់បច្ចុប្បន្ន', width: 140 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
          { label: 'ចំណាំ', prop: 'notes', width: 120 },
        ]"
      >
      <template #name="{row}">
        <el-text>{{ row.name }} | <el-text size="small" type="primary">{{ row.customer_code }}</el-text></el-text>
      </template>
      <template #branch_name="{row}">
        <el-text>{{ row.branch_name }} | <el-text size="small" type="primary">{{ row.branch_code }}</el-text></el-text>
      </template>
      <template #credit_limit="{row}">
        <el-text>{{ row.credit_limit }} <el-text size="small" type="primary">{{ row.company_currency }}</el-text></el-text>
      </template>
      <template #current_outstanding="{row}">
        <el-text>{{ row.current_outstanding }} <el-text size="small" type="primary">{{ row.company_currency }}</el-text></el-text>
      </template>
        <template #status="{ row }">
          <el-text
            :type="
              row.status === 'ACTIVE'
                ? 'success'
                : row.status === 'BLACKLISTED'
                ? 'danger'
                : 'info'
            "
            size="small"
          >
            {{ row.status }}
          </el-text>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
              v-if="canEditCustomer"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="openEdit(row)"
            />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="dialogVisible"
      :title="isEdit ? 'កែប្រែអតិថិជន' : 'បន្ថែមអតិថិជន'"
      width="500px"
      :showDefaultFooter="false"
    >
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :show-actions="true"
        @submit="handleSave"
        submitText="រក្សាទុក"
      >
        <AppInput
          v-model="form.name"
          label="ឈ្មោះអតិថិជន"
          prop="name"
          placeholder="បញ្ចូលឈ្មោះអតិថិជន"
        />

        <AppInput
          v-model="form.phone"
          label="ទូរស័ព្ទ"
          prop="phone"
          placeholder="បញ្ចូលលេខទូរស័ព្ទ"
        />

        <AppInput
          v-model="form.address"
          label="អាសយដ្ឋាន"
          prop="address"
          placeholder="បញ្ចូលអាសយដ្ឋាន"
        />

        <AppInput
          v-model="form.notes"
          label="កំណត់ចំណាំ"
          prop="notes"
          type="textarea"
          placeholder="បញ្ចូលកំណត់ចំណាំ"
        />

        <AppInput
          v-model.number="form.credit_limit"
          label="កម្រិតឥណទាន"
          prop="credit_limit"
          type="number"
          placeholder="បញ្ចូលកម្រិតឥណទាន"
        />

        <template v-if="isEdit">
          <AppSelect
            v-model="form.credit_limit_enforced"
            label="តម្រូវឲ្យអនុវត្តកម្រិតឥណទាន"
            :options="[
              { label: 'YES', value: true },
              { label: 'NO', value: false },
            ]"
          />

          <AppSelect
            v-model="form.status"
            :options="StatusOption"
            label="ស្ថានភាព"
            placeholder="ស្ថានភាព"
            clearable
          />
        </template>
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 6px;
}
</style>