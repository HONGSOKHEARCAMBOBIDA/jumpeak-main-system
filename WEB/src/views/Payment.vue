<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getpayment,
  addpayment,
  voidpayment,
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

const payments = ref([]);
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
  status: null,
  method: null,
});

const StatusOption = [
  { label: "COMPLETED", value: "COMPLETED" },
  { label: "VOIDED", value: "VOIDED" },
];

const MethodOption = [
  { label: "CASH", value: "CASH" },
  { label: "BANK", value: "BANK" },
  { label: "ABA", value: "ABA" },
  { label: "ACLEDA", value: "ACLEDA" },
  { label: "OTHER", value: "OTHER" },
];

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

// --- outstanding invoices for the picked customer (allocation targets) ---
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
        label: `${i.invoice_number} — នៅជំពាក់ ${i.outstanding_amount}${i.currency_code}`,
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
function blankAllocation() {
  return { key: Date.now() + Math.random(), invoice_id: null, amount: 0 };
}

const form = reactive({
  customer_id: null,
  payment_date: "",
  currency_code: "KHR",
  exchange_rate_to_base: 1,
  amount: 0,
  method: "CASH",
  reference_number: "",
  note: "",
});

const allocations = ref([blankAllocation()]);

const rules = {
  customer_id: [{ required: true, message: "Customer is required" }],
  payment_date: [{ required: true, message: "Payment date is required" }],
  currency_code: [{ required: true, message: "Currency is required" }],
  amount: [{ required: true, message: "Amount is required" }],
  method: [{ required: true, message: "Method is required" }],
};

function addAllocationRow() {
  allocations.value.push(blankAllocation());
}
function removeAllocationRow(key) {
  if (allocations.value.length === 1) return;
  allocations.value = allocations.value.filter((r) => r.key !== key);
}
function invoiceOutstanding(invoiceId) {
  return openInvoiceOptions.value.find((o) => o.value === invoiceId)?.raw?.outstanding_amount ?? 0;
}
const allocatedTotal = computed(() =>
  allocations.value.reduce((sum, row) => sum + (Number(row.amount) || 0), 0),
);
const unallocated = computed(() => (Number(form.amount) || 0) - allocatedTotal.value);

async function onCustomerPickedInForm(customerId) {
  allocations.value = [blankAllocation()];
  await loadOpenInvoices(customerId);
}

const canAddPayment = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Payment"),
);
const canVoidPayment = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "Void.Payment"),
);

async function fetchPayments() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.customer_id) params.customer_id = filters.customer_id;
    if (filters.status) params.status = filters.status;
    if (filters.method) params.method = filters.method;
    const res = await getpayment(params);
    payments.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  form.customer_id = null;
  form.payment_date = "";
  form.currency_code = "KHR";
  form.exchange_rate_to_base = 1;
  form.amount = 0;
  form.method = "CASH";
  form.reference_number = "";
  form.note = "";
  allocations.value = [blankAllocation()];
  openInvoiceOptions.value = [];
  searchCustomers("");
  dialogVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  const cleanAllocations = allocations.value
    .filter((r) => r.invoice_id && Number(r.amount) > 0)
    .map((r) => ({ invoice_id: r.invoice_id, amount: Number(r.amount) }));

  if (unallocated.value < 0) {
    notify.error("ចំនួនបែងចែកលើសពីចំនួនទឹកប្រាក់បង់");
    return;
  }

  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    const payload = {
      customer_id: form.customer_id,
      payment_date: form.payment_date,
      currency_code: form.currency_code,
      exchange_rate_to_base: Number(form.exchange_rate_to_base) || 1,
      amount: Number(form.amount) || 0,
      method: form.method,
      reference_number: form.reference_number || null,
      note: form.note || null,
      allocations: cleanAllocations,
    };
    await addpayment(payload);
    notify.success("បង្កើតការទូទាត់បានជោគជ័យ");
    dialogVisible.value = false;
    fetchPayments();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

// --- void ---
const voidDialogVisible = ref(false);
const voidTargetId = ref(null);
const voidReason = ref("");

function openVoid(row) {
  voidTargetId.value = row.id;
  voidReason.value = "";
  voidDialogVisible.value = true;
}

async function confirmVoid() {
  if (!voidReason.value.trim()) {
    notify.error("សូមបញ្ចូលមូលហេតុ");
    return;
  }
  useloading.show({ text: "កំពុងលុបចោល..." });
  try {
    await voidpayment(voidTargetId.value, { reason: voidReason.value });
    notify.success("បានលុបចោលការទូទាត់");
    voidDialogVisible.value = false;
    fetchPayments();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    useloading.hide();
  }
}

onMounted(() => {
  fetchPayments();
  searchCustomers("");
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'customer', span: 6 },
        { slot: 'status', span: 5 },
        { slot: 'method', span: 5 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #customer>
        <AppSelect
          v-model="filters.customer_id"
          :options="customerOptions"
          label="អតិថិជន"
          placeholder="អតិថិជន"
          filterable
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          clearable
          @change="fetchPayments"
        />
      </template>
      <template #status>
        <AppSelect
          v-model="filters.status"
          :options="StatusOption"
          label="ស្ថានភាព"
          placeholder="ស្ថានភាព"
          clearable
          @change="fetchPayments"
        />
      </template>
      <template #method>
        <AppSelect
          v-model="filters.method"
          :options="MethodOption"
          label="មធ្យោបាយបង់ប្រាក់"
          placeholder="មធ្យោបាយបង់ប្រាក់"
          clearable
          @change="fetchPayments"
        />
      </template>
      <template #create>
        <AppButton
          v-if="canAddPayment"
          type="primary"
          @click="openCreate"
          :block="false"
          size="default"
        >
          បង្កើតការទូទាត់
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="payments"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchPayments"
        :columns="[
          { prop: 'payment_number', label: 'លេខបង់ប្រាក់', minWidth: 140 },
          { prop: 'customer_name', label: 'អតិថិជន', minWidth: 150 },
          { prop: 'payment_date', label: 'កាលបរិច្ឆេទ', width: 120 },
          { slot: 'amount', label: 'ចំនួនទឹកប្រាក់', width: 130 },
          { prop: 'method', label: 'មធ្យោបាយ', width: 110 },
          { prop: 'reference_number', label: 'លេខយោង', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 110 },
        ]"
      >
        <template #amount="{ row }">
          <el-text tag="b">{{ row.amount }} <el-text size="small" type="primary">{{ row.currency_code }}</el-text></el-text>
        </template>
        <template #status="{ row }">
          <el-text :type="row.status === 'COMPLETED' ? 'success' : 'info'" size="small">
            {{ row.status }}
          </el-text>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="លុបចោល" placement="top">
            <AppButton
              v-if="canVoidPayment && row.status === 'COMPLETED'"
              size="small"
              icon="CircleClose"
              type="danger"
              circle
              @click="openVoid(row)"
            />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <AppDialog v-model="dialogVisible" title="បង្កើតការទូទាត់" width="45%" :showDefaultFooter="false">
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
          size="large"
          filterable
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          @change="onCustomerPickedInForm"
        />

        <el-row :gutter="16">
          <el-col :span="12">
            <AppInput v-model="form.payment_date" label="កាលបរិច្ឆេទបង់ប្រាក់" prop="payment_date" type="date" />
          </el-col>
          <el-col :span="12">
            <AppSelect v-model="form.method" :options="MethodOption" size="large" label="មធ្យោបាយបង់ប្រាក់" prop="method" />
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="8">
            <AppInput v-model="form.currency_code" label="រូបិយប័ណ្ណ" prop="currency_code" placeholder="USD" />
          </el-col>
          <el-col :span="8">
            <AppInput v-model.number="form.exchange_rate_to_base" label="អត្រាប្តូរប្រាក់" type="number" />
          </el-col>
          <el-col :span="8">
            <AppInput v-model.number="form.amount" label="ចំនួនទឹកប្រាក់" prop="amount" type="number" />
          </el-col>
        </el-row>

        <AppInput v-model="form.reference_number" label="លេខយោង" placeholder="លេខយោង (ស្រេចចិត្ត)" />
        <AppInput v-model="form.note" label="ចំណាំ" type="textarea" placeholder="ចំណាំ (ស្រេចចិត្ត)" />

        <el-divider content-position="left">បែងចែកទៅលើវិក័យបត្រ</el-divider>

        <div class="item-rows">
          <div v-for="row in allocations" :key="row.key" class="item-row">
            <el-row :gutter="10">
              <el-col :span="14">
                <AppSelect
                  v-model="row.invoice_id"
                  :options="openInvoiceOptions"
                  label="លេខវិក័យបត្រ"
                  placeholder="ជ្រើសរើសវិក័យបត្រ"
                  size="large"
                  :loading="invoicesLoading"
                  clearable
                />
              </el-col>
              <el-col :span="6">
                <AppInput v-model.number="row.amount" type="number" placeholder="ចំនួនប្រាក់សង" label="ចំនួនប្រាក់សង" />
              </el-col>
              <el-col :span="2" class="item-subtotal">
                {{ invoiceOutstanding(row.invoice_id) }} 
              </el-col>
              <el-col :span="2">
                <AppButton
                  v-if="allocations.length > 1"
                  size="small"
                  icon="Delete"
                  type="danger"
                  circle
                  @click="removeAllocationRow(row.key)"
                />
              </el-col>
            </el-row>
          </div>
        </div>

       
                  <div class="item-actions">
          <AppButton size="default" type="primary" plain icon="Plus" @click="addAllocationRow">
            បន្ថែមការបែងចែក
          </AppButton>
          <el-text tag="b" :type="unallocated < 0 ? 'danger' : 'info'">
            មិនទាន់បែងចែក: {{ unallocated.toFixed(2) }}
          </el-text>
        </div>
    
      </AppForm>
    </AppDialog>

    <AppDialog v-model="voidDialogVisible" title="លុបចោលការទូទាត់" width="420px" :showDefaultFooter="false">
      <AppInput
        v-model="voidReason"
        label="មូលហេតុ"
        type="textarea"
        placeholder="បញ្ចូលមូលហេតុនៃការលុបចោល"
      />
      <div class="dialog-actions">
        <AppButton type="default" @click="voidDialogVisible = false">បោះបង់</AppButton>
        <AppButton type="danger" @click="confirmVoid">បញ្ជាក់ការលុបចោល</AppButton>
      </div>
    </AppDialog>
  </div>
</template>

<style scoped>
.table-card {
  border-radius: 6px;
}

.item-rows {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-subtotal {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.item-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  margin-bottom: 10px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}
</style>