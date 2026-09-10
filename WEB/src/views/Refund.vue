<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getrefund,
  addrefund,
  getcustomer,
  getpayment,
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

const refunds = ref([]);
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
});

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

// --- payments for the picked customer (only completed ones can be refunded) ---
const paymentOptions = ref([]);
const paymentsLoading = ref(false);
async function loadCustomerPayments(customerId) {
  paymentOptions.value = [];
  if (!customerId) return;
  paymentsLoading.value = true;
  try {
    const res = await getpayment({ customer_id: customerId, status: "COMPLETED", page: 1, page_size: 50 });
    paymentOptions.value = (res.data.data || []).map((p) => ({
      label: `${p.payment_number} — ${p.amount} ${p.currency_code}`,
      value: p.id,
      raw: p,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to load payments");
  } finally {
    paymentsLoading.value = false;
  }
}

// --- invoices the customer has actually paid something on (refund targets) ---
const paidInvoiceOptions = ref([]);
const invoicesLoading = ref(false);
async function loadPaidInvoices(customerId) {
  paidInvoiceOptions.value = [];
  if (!customerId) return;
  invoicesLoading.value = true;
  try {
    const res = await getinvoice({ customer_id: customerId, page: 1, page_size: 100 });
    paidInvoiceOptions.value = (res.data.data || [])
      .filter((i) => i.paid_amount > 0)
      .map((i) => ({
        label: `${i.invoice_number} — បានបង់ ${i.paid_amount}`,
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
  payment_id: null,
  amount: 0,
  reason: "",
  refunded_at: "",
});

const allocations = ref([blankAllocation()]);

const rules = {
  customer_id: [{ required: true, message: "Customer is required" }],
  payment_id: [{ required: true, message: "Payment is required" }],
  amount: [{ required: true, message: "Amount is required" }],
  reason: [{ required: true, message: "Reason is required" }],
  refunded_at: [{ required: true, message: "Refund date is required" }],
};

function addAllocationRow() {
  allocations.value.push(blankAllocation());
}
function removeAllocationRow(key) {
  if (allocations.value.length === 1) return;
  allocations.value = allocations.value.filter((r) => r.key !== key);
}
const allocatedTotal = computed(() =>
  allocations.value.reduce((sum, row) => sum + (Number(row.amount) || 0), 0),
);
const allocationDiff = computed(() => (Number(form.amount) || 0) - allocatedTotal.value);

async function onCustomerPickedInForm(customerId) {
  form.payment_id = null;
  allocations.value = [blankAllocation()];
  await Promise.all([loadCustomerPayments(customerId), loadPaidInvoices(customerId)]);
}

const canAddRefund = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Refund"),
);

async function fetchRefunds() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.customer_id) params.customer_id = filters.customer_id;
    const res = await getrefund(params);
    refunds.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  form.customer_id = null;
  form.payment_id = null;
  form.amount = 0;
  form.reason = "";
  form.refunded_at = "";
  allocations.value = [blankAllocation()];
  paymentOptions.value = [];
  paidInvoiceOptions.value = [];
  searchCustomers("");
  dialogVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  const cleanAllocations = allocations.value
    .filter((r) => r.invoice_id && Number(r.amount) > 0)
    .map((r) => ({ invoice_id: r.invoice_id, amount: Number(r.amount) }));

  if (cleanAllocations.length === 0) {
    notify.error("សូមបែងចែកទឹកប្រាក់ត្រឡប់ទៅលើវិក័យបត្រយ៉ាងតិចមួយ");
    return;
  }
  if (allocationDiff.value !== 0) {
    notify.error("ការបែងចែកត្រូវតែស្មើនឹងចំនួនទឹកប្រាក់ត្រឡប់សរុប");
    return;
  }

  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    const payload = {
      payment_id: form.payment_id,
      amount: Number(form.amount) || 0,
      reason: form.reason,
      refunded_at: form.refunded_at,
      allocations: cleanAllocations,
    };
    await addrefund(payload);
    notify.success("បង្កើតការត្រឡប់ប្រាក់បានជោគជ័យ");
    dialogVisible.value = false;
    fetchRefunds();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

onMounted(() => {
  fetchRefunds();
  searchCustomers("");
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'customer', span: 8 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #customer>
        <AppSelect
          v-model="filters.customer_id"
          :options="customerOptions"
          label="អតិថិជន"
          placeholder="អតិថិជន"
          size="large"
          filterable
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          clearable
          @change="fetchRefunds"
        />
      </template>
      <template #create>
        <AppButton
          v-if="canAddRefund"
          type="primary"
          @click="openCreate"
          :block="false"
          size="large"
        >
          បង្កើតការត្រឡប់ប្រាក់
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="refunds"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchRefunds"
        :columns="[
          { prop: 'payment_number', label: 'ការទូទាត់ដើម', minWidth: 140 },
          { prop: 'customer_name', label: 'អតិថិជន', minWidth: 150 },
          { slot: 'amount', label: 'ចំនួនត្រឡប់', width: 130 },
          { prop: 'reason', label: 'មូលហេតុ', minWidth: 160 },
          { prop: 'refunded_at', label: 'កាលបរិច្ឆេទ', width: 120 },
        ]"
      >
        <template #amount="{ row }">
          <el-text tag="b" type="danger">{{ row.amount }}</el-text>
        </template>
      </AppTable>
    </el-card>

    <AppDialog v-model="dialogVisible" title="បង្កើតការត្រឡប់ប្រាក់" width="45%" :showDefaultFooter="false">
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
          size="large"
          prop="customer_id"
          placeholder="ជ្រើសរើសអតិថិជន"
          filterable
          remote
          :remote-method="searchCustomers"
          :loading="customerSearching"
          @change="onCustomerPickedInForm"
        />

        <AppSelect
          v-model="form.payment_id"
          :options="paymentOptions"
          label="ការទូទាត់ដើម"
          size="large"
          prop="payment_id"
          placeholder="ជ្រើសរើសការទូទាត់"
          :loading="paymentsLoading"
        />

        <el-row :gutter="16">
          <el-col :span="12">
            <AppInput v-model.number="form.amount" label="ចំនួនត្រឡប់" prop="amount" type="number" />
          </el-col>
          <el-col :span="12">
            <AppInput v-model="form.refunded_at" label="កាលបរិច្ឆេទត្រឡប់" prop="refunded_at" type="date" />
          </el-col>
        </el-row>

        <AppInput v-model="form.reason" label="មូលហេតុ" prop="reason" type="textarea" placeholder="បញ្ចូលមូលហេតុ" />

        <el-divider content-position="left">បែងចែកចេញពីវិក័យបត្រ</el-divider>

        <div class="item-rows">
          <div v-for="row in allocations" :key="row.key" class="item-row">
            <el-row :gutter="10">
              <el-col :span="16">
                <AppSelect
                  v-model="row.invoice_id"
                  :options="paidInvoiceOptions"
                   size="large"
                  placeholder="ជ្រើសរើសវិក័យបត្រ"
                  :loading="invoicesLoading"
                  clearable
                />
              </el-col>
              <el-col :span="6">
                <AppInput v-model.number="row.amount" type="number" placeholder="ចំនួន" />
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
          <AppButton size="default" type="default" icon="Plus" @click="addAllocationRow">
            បន្ថែមការបែងចែក
          </AppButton>
          <el-text tag="b" :type="allocationDiff !== 0 ? 'danger' : 'success'">
            ភាពខុសគ្នា: {{ allocationDiff.toFixed(2) }}
          </el-text>
        </div>
      </AppForm>
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

.item-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  margin-bottom: 10px;
}
</style>