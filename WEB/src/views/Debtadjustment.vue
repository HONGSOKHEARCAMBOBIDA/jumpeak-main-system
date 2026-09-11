<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getdebtadjustment,
  adddebtadjustment,
  getcustomer,
  getinvoice,
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

const adjustments = ref([]);
const loading = ref(false);
const useloading = useLoading();
const saving = ref(false);

const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogVisible = ref(false);
const formRef = ref();

// filters
const filters = reactive({
  customer_id: null,
  type: null,
});

const TypeOption = [
  { label: "លុបបំណុល", value: "WRITE_OFF" },
  { label: "កែតម្រូវ", value: "CORRECTION" },
  { label: "បញ្ចុះតម្លៃ", value: "DISCOUNT" },
];

const getTypeLabel = (status) => {
  return TypeOption.find((item) => item.value === status)?.label || status;
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

// --- optional invoice targeting for the picked customer ---
const openInvoiceOptions = ref([]);
const invoicesLoading = ref(false);
async function loadOpenInvoices(customerId) {
  openInvoiceOptions.value = [];
  if (!customerId) return;
  invoicesLoading.value = true;
  try {
    const res = await getinvoice({ customer_id: customerId, page: 1, page_size: 100 });
    openInvoiceOptions.value = (res.data.data || [])
      .filter((i) => i.outstanding_amount > 0)
      .map((i) => ({
        label: `${i.invoice_number} — នៅជំពាក់ ${i.outstanding_amount}`,
        value: i.id,
        raw: i,
      }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to load invoices");
  } finally {
    invoicesLoading.value = false;
  }
}

// --- create form ---
const form = reactive({
  customer_id: null,
  invoice_id: null,
  type: "WRITE_OFF",
  amount: 0,
  reason: "",
});

const rules = {
  customer_id: [{ required: true, message: "Customer is required" }],
  type: [{ required: true, message: "Type is required" }],
  amount: [{ required: true, message: "Amount is required" }],
  reason: [{ required: true, message: "Reason is required" }],
};

async function onCustomerPickedInForm(customerId) {
  form.invoice_id = null;
  await loadOpenInvoices(customerId);
}

const invoiceOutstanding = computed(
  () => openInvoiceOptions.value.find((o) => o.value === form.invoice_id)?.raw?.outstanding_amount ?? null,
);

const canAddAdjustment = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.DebAdjustment"),
);

async function fetchAdjustments() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.customer_id) params.customer_id = filters.customer_id;
    if (filters.type) params.type = filters.type;
    const res = await getdebtadjustment(params);
    adjustments.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  form.customer_id = null;
  form.invoice_id = null;
  form.type = "WRITE_OFF";
  form.amount = 0;
  form.reason = "";
  openInvoiceOptions.value = [];
  searchCustomers("");
  dialogVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  if (form.invoice_id && invoiceOutstanding.value !== null && Number(form.amount) > invoiceOutstanding.value) {
    notify.error("ចំនួនកែសម្រួលលើសពីសមតុល្យវិក័យបត្រ");
    return;
  }

  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    const payload = {
      customer_id: form.customer_id,
      invoice_id: form.invoice_id || null,
      type: form.type,
      amount: Number(form.amount) || 0,
      reason: form.reason,
    };
    await adddebtadjustment(payload);
    notify.success("បង្កើតការកែសម្រួលបំណុលបានជោគជ័យ");
    dialogVisible.value = false;
    fetchAdjustments();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

onMounted(() => {
  fetchAdjustments();
  searchCustomers("");
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'customer', span: 7 },
        { slot: 'type', span: 6 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #customer>
        <AppSelect
          v-model="filters.customer_id"
          :options="customerOptions"
          label="អតិថិជន"
          size="large"
          placeholder="អតិថិជន"
          filterable
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          clearable
          @change="fetchAdjustments"
        />
      </template>
      <template #type>
        <AppSelect
          v-model="filters.type"
          :options="TypeOption"
          label="ប្រភេទ"
           size="large"
          placeholder="ប្រភេទ"
          clearable
          @change="fetchAdjustments"
        />
      </template>
      <template #create>
        <AppButton
          v-if="canAddAdjustment"
          type="primary"
          @click="openCreate"
          :block="false"
           size="large"
        >
          បង្កើតការកែសម្រួលបំណុល
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="adjustments"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchAdjustments"
        :columns="[
          { prop: 'customer_name', label: 'អតិថិជន', minWidth: 150 },
          { prop: 'invoice_number', label: 'លេខវិក័យបត្រ', width: 140 },
           { label: 'ក្រុមហ៑ុន', prop: 'company_Name', width: 180 },
            { label: 'សាខា', slot: 'branch_Name', width: 200 },
          { label: 'ប្រភេទ', slot: 'type', width: 120 },
          { slot: 'amount', label: 'ចំនួន', width: 120 },
          { prop: 'reason', label: 'មូលហេតុ', minWidth: 160 },
          { prop: 'created_at', label: 'កាលបរិច្ឆេទ', width: 150 },
          { prop: 'approved_by', label: 'កែបំណុលដោយ', width: 200 },
        ]"
      >
      <template #branch_Name="{row}">
        <el-text>{{ row.branch_Name }} | <el-text size="small" type="primary">{{ row.branch_Code }}</el-text></el-text>
      </template>
        <template #type="{ row }">
          <el-text
            :type="row.type === 'WRITE_OFF' ? 'danger' : row.type === 'DISCOUNT' ? 'warning' : 'info'"
            size="small"
          >
            {{ getTypeLabel(row.type) }}
          </el-text>
        </template>
        <template #amount="{ row }">
          <el-text tag="b" type="danger">{{ row.amount }} <el-text size="small" type="primary">{{ row.currency }}</el-text></el-text>
        </template>
      </AppTable>
    </el-card>

    <AppDialog v-model="dialogVisible" title="បង្កើតការកែសម្រួលបំណុល" width="40%" :showDefaultFooter="false">
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :show-actions="true"
        @submit="handleSave"
        submitText="រក្សាទុក"
      >
        <AppSelect
          v-model="form.customer_id"
          :options="customerOptions"
          label="អតិថិជន"
          prop="customer_id"
          placeholder="ជ្រើសរើសអតិថិជន"
          filterable
          size="large"
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          @change="onCustomerPickedInForm"
        />

        <AppSelect
          v-model="form.invoice_id"
          :options="openInvoiceOptions"
          label="វិក័យបត្រ (ស្រេចចិត្ត)"
          size="large"
          placeholder="ជ្រើសរើសវិក័យបត្រ បើមាន"
          :loading="invoicesLoading"
          clearable
        />

        <el-row :gutter="16">
          <el-col :span="12">
            <AppSelect v-model="form.type" :options="TypeOption" size="large" label="ប្រភេទ" prop="type" />
          </el-col>
          <el-col :span="12">
            <AppInput v-model.number="form.amount" label="ចំនួន" size="large" prop="amount" type="number" />
          </el-col>
        </el-row>

        <AppInput v-model="form.reason" label="មូលហេតុ" prop="reason" size="large" type="textarea" placeholder="បញ្ចូលមូលហេតុ" />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.table-card {
  border-radius: 6px;
}
</style>