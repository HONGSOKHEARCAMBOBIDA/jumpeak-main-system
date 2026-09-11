<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { getcustomerledger, getcustomer } from "../api/services.js";
import AppTable from "../../components/AppTable.vue";
import { useNotification } from "../../composables/useNotification.js";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";

const notify = useNotification();

const entries = ref([]);
const loading = ref(false);

const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

// filters
const filters = reactive({
  customer_id: null,
  reference_type: null,
  date_from: "",
  date_to: "",
});

const ReferenceTypeOption = [
  { label: "វិក្កយបត្រ", value: "INVOICE" },
  { label: "ការសងប្រាក់", value: "PAYMENT" },
  { label: "បង់ប្រាក់ទៅអតិថិជនវិញ", value: "REFUND" },
  { label: "កែបំណុល", value: "ADJUSTMENT" },
];

const ReferenceTypeLabel = (status) => {
  return ReferenceTypeOption.find((item) => item.value === status)?.label || status;
};
// --- customer remote search ---
const customerOptions = ref([]);
const customerSearching = ref(false);
async function searchCustomers(query) {
  customerSearching.value = true;
  try {
    const res = await getcustomer({ page: 1, page_size: 20, name: query || "" });
    customerOptions.value = (res.data.data || []).map((c) => ({
      label: `${c.name} (${c.customer_code})`,
      value: c.id,
      raw: c,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to search customers");
  } finally {
    customerSearching.value = false;
  }
}

async function fetchEntries() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.customer_id) params.customer_id = filters.customer_id;
    if (filters.reference_type) params.reference_type = filters.reference_type;
    if (filters.date_from) params.date_from = filters.date_from;
    if (filters.date_to) params.date_to = filters.date_to;
    const res = await getcustomerledger(params);
    entries.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function onFilterChange() {
  page.value = 1;
  fetchEntries();
}

const currentBalance = computed(() => {
  if (entries.value.length === 0) return null;
  return entries.value[entries.value.length - 1].running_balance;
});

onMounted(() => {
  fetchEntries();
  searchCustomers("");
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'customer', span: 6 },
        { slot: 'reference_type', span: 5 },
        { slot: 'date_from', span: 4 },
        { slot: 'date_to', span: 4 },
      ]"
    >
      <template #customer>
        <AppSelect
          v-model="filters.customer_id"
          :options="customerOptions"
          label="អតិថិជន"
          placeholder="ជ្រើសរើសអតិថិជន"
          size="large"
          filterable
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          clearable
          @change="onFilterChange"
        />
      </template>
      <template #reference_type>
        <AppSelect
          v-model="filters.reference_type"
          :options="ReferenceTypeOption"
          label="ប្រភេទ"
          size="large"
          placeholder="ប្រភេទប្រតិបត្តិការ"
          clearable
          @change="onFilterChange"
        />
      </template>
      <template #date_from>
        <AppInput v-model="filters.date_from" type="date" placeholder="ពីថ្ងៃទី" @change="onFilterChange" />
      </template>
      <template #date_to>
        <AppInput v-model="filters.date_to" type="date" placeholder="ដល់ថ្ងៃទី" @change="onFilterChange" />
      </template>
    </AppFilterBar>

    <el-card v-if="filters.customer_id && currentBalance !== null" class="balance-card">
      <el-text>ប្រាក់ជំពាក់នៅសល់ : </el-text>
      <el-text tag="b" size="large" :type="currentBalance > 0 ? 'danger' : 'success'">
        {{ currentBalance }}
      </el-text>
    </el-card>

    <el-card class="table-card">
      <AppTable
        :data="entries"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchEntries"
        :columns="[
          { prop: 'entry_date', label: 'កាលបរិច្ឆេទ', width: 100 },
          { prop: 'customer_name', label: 'អតិថិជន', minWidth: 80 },
          { prop: 'company_Name', label: 'ក្រុមហ៑ុន', Width: 150 },
          { slot: 'branch_Name', label: 'សាខា', Width: 200 },
          { label: 'ប្រភេទ', slot: 'reference_type', width: 140 },
          { prop: 'description', label: 'បរិយាយ', minWidth: 200 },
          { slot: 'debit', label: 'ជំពាក់កើន (Debit)', width: 150 },
          { slot: 'credit', label: 'បង់កើន (Credit)', width: 150 },
          { slot: 'running_balance', label: 'ប្រាក់ជំពាក់នៅសល់', width: 140 },
        ]"
      >
      <template #branch_Name="{row}">
        <el-text>{{ row.branch_Name }} | <el-text size="small" type="primary">{{ row.branch_Code }}</el-text></el-text>
      </template>
        <template #reference_type="{ row }">
          <el-text size="small" type="primary">{{ ReferenceTypeLabel(row.reference_type) }}</el-text>
        </template>
        <template #debit="{ row }">
          <el-text v-if="row.debit > 0" type="danger">{{ row.debit }}</el-text>
          <el-text v-else type="info">—</el-text>
        </template>
        <template #credit="{ row }">
          <el-text v-if="row.credit > 0" type="success">{{ row.credit }}</el-text>
          <el-text v-else type="info">—</el-text>
        </template>
        <template #running_balance="{ row }">
          <el-text tag="b" style="color: black;">{{ row.running_balance }}</el-text>
        </template>
      </AppTable>
    </el-card>
  </div>
</template>

<style scoped>
.table-card {
  border-radius: 6px;
}

.balance-card {
  margin-bottom: 12px;
  border-radius: 6px;
}
</style>