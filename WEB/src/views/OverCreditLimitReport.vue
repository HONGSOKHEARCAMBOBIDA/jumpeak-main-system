<script setup>
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getovercreditlitmireport,
  getcompanynopagitaion,
  getbranchnopagination,
} from "../api/services.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppForm from "../../components/form/AppForm.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
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
  company_id: null,
  branch_id: null,
});

const companyoptions = ref([]);
const branchoptions = ref([]);
async function fetchCompanyOptions() {
  try {
    const res = await getcompanynopagitaion();
    companyoptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }));
  } catch (e) {
    notify.error(
      e?.response?.data?.message || e.message || "Failed to load company",
    );
  }
}

async function loadBranchOption(companyID) {
  if (!companyID) return [];
  try {
    const res = await getbranchnopagination(companyID);
    return (res.data.data || []).map((f) => ({
      label: f.name,
      value: f.id,
    }));
  } catch (e) {
    notify.error(
      e?.response?.data?.message || e.message || "Failed to load branch",
    );
    return [];
  }
}

const StatusOption = [
  { label: "ACTIVE", value: "ACTIVE" },
  { label: "INACTIVE", value: "INACTIVE" },
  { label: "BLACKLISTED", value: "BLACKLISTED" },
];

async function fetchOverCreditlimti() {
  loading.value = true;
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    };
    if (filters.name) params.name = filters.name;
    if (filters.company_id) params.company_id = filters.company_id;
    if (filters.branch_id) params.branch_id = filters.branch_id;
    const res = await getovercreditlitmireport(params);
    customers.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

watch(
  () => filters.company_id,
  async (newVal) => {
    filters.branch_id = null;
    branchoptions.value = [];
    branchoptions.value = await loadBranchOption(newVal);
    page.value = 1;
    fetchOverCreditlimti();
  },
);

watch(
  () => filters.branch_id,
  async () => {
    page.value = 1;
    fetchOverCreditlimti();
  },
);

watch(
  () => filters.name,
  async () => {
    page.value = 1;
    fetchOverCreditlimti();
  },
);

onMounted(() => {
  fetchOverCreditlimti();
  fetchCompanyOptions();
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 4 },
        { slot: 'company', span: 4 },
        { slot: 'branch', span: 4 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #name>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះ"
          clearable
          size="small"
        />
      </template>
      <template #company>
        <AppSelect
          v-model="filters.company_id"
          :options="companyoptions"
          size="large"
          placeholder="ក្រុមហ៑ុន"
          clearable
        />
      </template>
      <template #branch>
        <AppSelect
          v-model="filters.branch_id"
          :options="branchoptions"
          size="large"
          placeholder="សាខា"
          clearable
        />
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="customers"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchOverCreditlimti"
        :columns="[
          { slot: 'name', label: 'ឈ្មោះ', minWidth: 150 },
          { prop: 'company_name', label: 'ក្រុមហ៊ុន', minWidth: 120 },
          { slot: 'branch_name', label: 'សាខា', minWidth: 120 },
          { prop: 'phone', label: 'ទូរស័ព្ទ', width: 130 },
          { prop: 'address', label: 'អាស័យដ្ឋាន', width: 130 },
          
          
          { slot: 'credit_limit', label: 'កម្រិតឥណទាន', width: 130 },
          {
            slot: 'current_outstanding',
            label: 'ជំពាក់បច្ចុប្បន្ន',
            width: 140,
          },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
          { label: 'ចំណាំ', prop: 'notes', width: 120 },
        ]"
      >
        <template #name="{ row }">
          <el-text
            >{{ row.name }} |
            <el-text size="small" type="primary">{{
              row.customer_code
            }}</el-text></el-text
          >
        </template>
        <template #branch_name="{ row }">
          <el-text
            >{{ row.branch_name }} |
            <el-text size="small" type="primary">{{
              row.branch_code
            }}</el-text></el-text
          >
        </template>
        <template #credit_limit="{ row }">
          <el-text tag="b" type="danger"
            >{{ row.credit_limit }}
            <el-text size="small" type="danger">{{
              row.company_currency
            }}</el-text></el-text
          >
        </template>
        <template #current_outstanding="{ row }">
          <el-text type="danger"
            >{{ row.current_outstanding }}
            <el-text size="small" type="danger">{{
              row.company_currency
            }}</el-text></el-text
          >
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
      </AppTable>
    </el-card>
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
