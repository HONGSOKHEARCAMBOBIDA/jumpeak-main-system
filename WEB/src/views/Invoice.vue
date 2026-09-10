<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getinvoice,
  addinvoice,
  cancelinvoice,
  getcustomer,
  getproduct,
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

const invoices = ref([]);
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
  invoice_number: "",
  customer_id: null,
  status: null,
});

const StatusOption = [
  { label: "OPEN", value: "OPEN" },
  { label: "PARTIALLY_PAID", value: "PARTIALLY_PAID" },
  { label: "PAID", value: "PAID" },
  { label: "CANCELLED", value: "CANCELLED" },
  { label: "WRITTEN_OFF", value: "WRITTEN_OFF" },
];

// --- customer remote search (used by the filter bar and the create form) ---
const customerOptions = ref([]);
const customerSearching = ref(false);
async function searchCustomers(query) {
  customerSearching.value = true;
  try {
    const res = await getcustomer({ page: 1, page_size: 20, name: query || "" });
    customerOptions.value = (res.data.data || []).map((c) => ({
      label: `${c.name} (${c.customer_code}) | ${c.status}`,
      value: c.id,
      raw: c,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to search customers");
  } finally {
    customerSearching.value = false;
  }
}

// --- product remote search (used by the item rows) ---
const productOptions = ref([]);
const productSearching = ref(false);
async function searchProducts(query) {
  productSearching.value = true;
  try {
    const res = await getproduct({ page: 1, page_size: 20, name: query || "" });
    productOptions.value = (res.data.data || []).map((p) => ({
      label: p.name,
      value: p.id,
      raw: p,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to search products");
  } finally {
    productSearching.value = false;
  }
}

const branchOptions = ref([]);
async function loadBranchOptionsForCompany(companyId) {
  if (!companyId) {
    branchOptions.value = [];
    return;
  }
  try {
    const res = await getbranchnopagination(companyId);
    branchOptions.value = (res.data.data || []).map((b) => ({
      label: b.name,
      value: b.id,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to load branch");
  }
}

// --- create form ---
function blankItem() {
  return {
    key: Date.now() + Math.random(),
    product_id: null,
    description: "",
    quantity: 1,
    unit_price: 0,
    discount_amount: 0,
  };
}

const form = reactive({
  customer_id: null,
  branch_id: null,
  invoice_date: "",
  due_date: "",
  currency_code: "KHR",
  exchange_rate_to_base: 1,
});

const items = ref([blankItem()]);

const rules = {
  customer_id: [{ required: true, message: "Customer is required" }],
  branch_id: [{ required: true, message: "Branch is required" }],
  invoice_date: [{ required: true, message: "Invoice date is required" }],
  due_date: [{ required: true, message: "Due date is required" }],
  currency_code: [{ required: true, message: "Currency is required" }],
};

function addItemRow() {
  items.value.push(blankItem());
}

function removeItemRow(key) {
  if (items.value.length === 1) return;
  items.value = items.value.filter((r) => r.key !== key);
}

function onProductPicked(row) {
  const opt = productOptions.value.find((o) => o.value === row.product_id);
  row.description = ""
  if (opt && !row.description) row.description = opt.label;
}

function itemSubtotal(row) {
  return (Number(row.quantity) || 0) * (Number(row.unit_price) || 0) - (Number(row.discount_amount) || 0);
}

const invoiceTotal = computed(() =>
  items.value.reduce((sum, row) => sum + itemSubtotal(row), 0),
);

async function onCustomerPickedInForm(customerId) {
  form.branch_id = null;
  const opt = customerOptions.value.find((o) => o.value === customerId);
  await loadBranchOptionsForCompany(opt?.raw?.company_id || null);
}

const canAddInvoice = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Invoice"),
);
const canCancelInvoice = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "Cancel.Invoice"),
);

async function fetchInvoices() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.invoice_number) params.invoice_number = filters.invoice_number;
    if (filters.customer_id) params.customer_id = filters.customer_id;
    if (filters.status) params.status = filters.status;
    const res = await getinvoice(params);
    invoices.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  form.customer_id = null;
  form.branch_id = null;
  form.invoice_date = "";
  form.due_date = "";
  form.currency_code = "KHR";
  form.exchange_rate_to_base = 1;
  items.value = [blankItem()];
  branchOptions.value = [];
  searchCustomers("");
  searchProducts("");
  dialogVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  const cleanItems = items.value
    .filter((r) => r.product_id || r.description)
    .map((r) => ({
      product_id: r.product_id || null,
      description: r.description,
      quantity: Number(r.quantity) || 0,
      unit_price: Number(r.unit_price) || 0,
      discount_amount: Number(r.discount_amount) || 0,
    }));
  if (cleanItems.length === 0) {
    notify.error("សូមបញ្ចូលមុខទំនិញយ៉ាងតិចមួយ");
    return;
  }

  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    const payload = {
      customer_id: form.customer_id,
      branch_id: form.branch_id,
      invoice_date: form.invoice_date,
      due_date: form.due_date,
      currency_code: form.currency_code,
      exchange_rate_to_base: Number(form.exchange_rate_to_base) || 1,
      items: cleanItems,
    };
    await addinvoice(payload);
    notify.success("បង្កើតវិក័យបត្របានជោគជ័យ");
    dialogVisible.value = false;
    fetchInvoices();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

// --- cancel ---
const cancelDialogVisible = ref(false);
const cancelTargetId = ref(null);
const cancelReason = ref("");

function openCancel(row) {
  cancelTargetId.value = row.id;
  cancelReason.value = "";
  cancelDialogVisible.value = true;
}

async function confirmCancel() {
  if (!cancelReason.value.trim()) {
    notify.error("សូមបញ្ចូលមូលហេតុ");
    return;
  }
  useloading.show({ text: "កំពុងលុបចោល..." });
  try {
    await cancelinvoice(cancelTargetId.value, { reason: cancelReason.value });
    notify.success("បានលុបចោលវិក័យបត្រ");
    cancelDialogVisible.value = false;
    fetchInvoices();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    useloading.hide();
  }
}

function canCancelRow(row) {
  return row.status === "OPEN" || row.status === "PARTIALLY_PAID";
}

onMounted(() => {
  fetchInvoices();
  searchCustomers("");
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'invoice_number', span: 5 },
        { slot: 'customer', span: 6 },
        { slot: 'status', span: 5 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #invoice_number>
        <AppInput
          v-model="filters.invoice_number"
          placeholder="ស្វែងរកលេខវិក័យបត្រ"
          clearable
          size="small"
          @change="fetchInvoices"
        />
      </template>
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
          @change="fetchInvoices"
        />
      </template>
      <template #status>
        <AppSelect
          v-model="filters.status"
          :options="StatusOption"
          label="ស្ថានភាព"
          placeholder="ស្ថានភាព"
          clearable
          @change="fetchInvoices"
        />
      </template>
      <template #create>
        <AppButton
          v-if="canAddInvoice"
          type="primary"
          @click="openCreate"
          :block="false"
          size="default"
        >
          បង្កើតវិក័យបត្រ
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="invoices"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchInvoices"
        :columns="[
          { prop: 'invoice_number', label: 'លេខវិក័យបត្រ', minWidth: 140 },
          { prop: 'customer_name', label: 'អតិថិជន', minWidth: 150 },
          { prop: 'invoice_date', label: 'កាលបរិច្ឆេទ', width: 120 },
          { prop: 'due_date', label: 'ថ្ងៃកំណត់', width: 120 },
          { slot: 'total_amount', label: 'សរុប', width: 120 },
          { slot: 'paid_amount', label: 'បានបង់', width: 120 },
          { slot: 'outstanding_amount', label: 'នៅជំពាក់', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 130 },
        ]"
      >
        <template #total_amount="{ row }">
          <el-text tag="b" type="primary">{{ row.total_amount }} <el-text size="small" type="primary">{{ row.currency_code }}</el-text></el-text>
        </template>
        <template #paid_amount="{ row }">
          <el-text tag="b" type="success">{{ row.paid_amount }} <el-text size="small" type="success">{{ row.currency_code }}</el-text></el-text>
        </template>
        <template #outstanding_amount="{ row }">
          <el-text tag="b" :type="row.outstanding_amount > 0 ? 'danger' : 'info'">{{ row.outstanding_amount }} <el-text size="small" type="danger">{{ row.currency_code }}</el-text></el-text>
        </template>
        <template #status="{ row }">
          <el-text
            :type="
              row.status === 'PAID'
                ? 'success'
                : row.status === 'CANCELLED' || row.status === 'WRITTEN_OFF'
                ? 'info'
                : row.status === 'PARTIALLY_PAID'
                ? 'warning'
                : 'primary'
            "
            size="small"
          >
            {{ row.status }}
          </el-text>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="លុបចោល" placement="top">
            <AppButton
              v-if="canCancelInvoice && canCancelRow(row)"
              size="small"
              icon="CircleClose"
              type="danger"
              circle
              @click="openCancel(row)"
            />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <AppDialog v-model="dialogVisible" title="បង្កើតវិក័យបត្រ" width="65%" :showDefaultFooter="false">
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :show-actions="true"
        @submit="handleSave"
        submitText="រក្សាទុក"
      >
        <el-row :gutter="16">
          <el-col :span="12">
            <AppSelect
              v-model="form.customer_id"
              :options="customerOptions"
              label="អតិថិជន"
              prop="customer_id"
              size="large"
              placeholder="ជ្រើសរើសអតិថិជន"
              filterable
              remote
              :remote-method="searchCustomers"
              :loading="customerSearching"
              @change="onCustomerPickedInForm"
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="form.branch_id"
              :options="branchOptions"
              label="សាខា"
              prop="branch_id"
              size="large"
              placeholder="ជ្រើសរើសសាខា"
            />
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <AppInput
              v-model="form.invoice_date"
              label="កាលបរិច្ឆេទវិក័យបត្រ"
              prop="invoice_date"
              type="date"
            />
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model="form.due_date"
              label="ថ្ងៃកំណត់បង់"
              prop="due_date"
              type="date"
            />
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <AppInput
              v-model="form.currency_code"
              label="រូបិយប័ណ្ណ"
              prop="currency_code"
              placeholder="រូបិយប័ណ្ណ"
            />
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model.number="form.exchange_rate_to_base"
              label="អត្រាប្តូរប្រាក់"
              prop="exchange_rate_to_base"
              type="number"
            />
          </el-col>
        </el-row>

        <el-divider content-position="left">មុខទំនិញ</el-divider>

        <div class="item-rows">
          <div v-for="row in items" :key="row.key" class="item-row">
            <el-row :gutter="10">
              <el-col :span="6">
                <AppSelect
                  v-model="row.product_id"
                  :options="productOptions"
                  placeholder="ឥវ៉ាន់"
                  size="large"
                  label="ឥវ៉ាន់"
                  filterable
                  remote
                  :remote-method="searchProducts"
                  :loading="productSearching"
                  clearable
                  @change="onProductPicked(row)"
                />
              </el-col>
              <el-col :span="6">
                <AppInput v-model="row.description" placeholder="បរិយាយ" label="បរិយាយ"/>
              </el-col>
              <el-col :span="3">
                <AppInput v-model.number="row.quantity" type="number" placeholder="ចំនួន" label="ចំនួន"/>
              </el-col>
              <el-col :span="4">
                <AppInput v-model.number="row.unit_price" type="number" placeholder="តម្លៃឯកតា" label="តម្លៃឯកតា"/>
              </el-col>
              <el-col :span="3">
                <AppInput v-model.number="row.discount_amount" type="number" placeholder="បញ្ចុះតម្លៃ" label="បញ្ចុះតម្លៃ"/>
              </el-col>
              <el-col :span="2" class="item-subtotal">
                {{ itemSubtotal(row).toFixed(2) }} {{ form.currency_code }}
              </el-col>
              <el-col :span="2">
                <AppButton
                  v-if="items.length > 1"
                  size="small"
                  icon="Delete"
                  type="danger"
                  circle
                  @click="removeItemRow(row.key)"
                />
              </el-col>
            </el-row>
          </div>
        </div>

        <div class="item-actions">
          <AppButton size="default" type="default" icon="Plus" @click="addItemRow">
            បន្ថែមមុខទំនិញ
          </AppButton>
          <el-text tag="b" size="large">សរុប: {{ invoiceTotal.toFixed(2) }} {{ form.currency_code }}</el-text>
        </div>
      </AppForm>
    </AppDialog>

    <AppDialog v-model="cancelDialogVisible" title="លុបចោលវិក័យបត្រ" width="420px" :showDefaultFooter="false">
      <AppInput
        v-model="cancelReason"
        label="មូលហេតុ"
        type="textarea"
        placeholder="បញ្ចូលមូលហេតុនៃការលុបចោល"
      />
      <div class="dialog-actions">
        <AppButton type="default" @click="cancelDialogVisible = false">បោះបង់</AppButton>
        <AppButton type="danger" @click="confirmCancel">បញ្ជាក់ការលុបចោល</AppButton>
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
  justify-content: flex-end;
  font-weight: 600;
}

.item-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  margin-top: 20px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}
</style>