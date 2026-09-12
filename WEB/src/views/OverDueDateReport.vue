<script setup>
import { ref, reactive, onMounted, computed, watch, nextTick } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
getoverdatereport,
  getproduct,
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
import PrintableTable from "../../components/PrintableTable.vue"; // adjust path
const notify = useNotification();
const userDataStore = useUserDataStore();

const invoices = ref([]);
const loading = ref(false);
const useloading = useLoading();
const saving = ref(false);


const showDetail = ref(false);
const selectedInvoice = ref(null);

function openDetail(row) {
  selectedInvoice.value = row;
  showDetail.value = true;
}

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
  { label: "កំពុងជំពាក់", value: "OPEN" },
  { label: "បានបង់មួយផ្នែក", value: "PARTIALLY_PAID" },
  { label: "បានបង់រួច", value: "PAID" },
  { label: "បានលុបចោល", value: "CANCELLED" },
  { label: "បានលុបបំណុល", value: "WRITTEN_OFF" },
];

const invoiceitemcolumns = [
  { prop: "product_Name", label: "មុខទំនិញ", minwidth: 120 },
  { prop: "quantity", label: "ចំនួន", minwidth: 120 },
  { slot: "unit_price", label: "តម្លៃរាយ", minwidth: 120 },
  { slot: "discount_amount", label: "បញ្ចុះតម្លៃ", minwidth: 120 },
  { slot: "subtotal", label: "តម្លៃសរុប", minwidth: 120 },
  { prop: "description", label: "ផ្សេងៗ", minwidth: 120 },
];

const getStatusLabel = (status) => {
  return StatusOption.find((item) => item.value === status)?.label || status;
};



async function fetchOverDueDate() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.invoice_number) params.invoice_number = filters.invoice_number;
    if (filters.customer_id) params.customer_id = filters.customer_id;
    if (filters.status) params.status = filters.status;
    const res = await getoverdatereport(params);
    invoices.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}



watch(
  () => filters.invoice_number,
  async () => {
    page.value = 1;
    fetchOverDueDate();
  },
);

onMounted(() => {
  fetchOverDueDate();
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
        />
      </template>
      <template #customer>
        <AppSelect
          v-model="filters.customer_id"
          :options="customerOptions"
          placeholder="អតិថិជន"
          size="large"
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
          size="large"
          placeholder="ស្ថានភាព"
          clearable
          @change="fetchInvoices"
        />
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        expandable
        :data="invoices"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchInvoices"
        actions-width="120px"
        :columns="[
          { prop: 'invoice_number', label: 'លេខវិក័យបត្រ', minWidth: 140 },
          { prop: 'customer_name', label: 'អតិថិជន', minWidth: 150 },
          { prop: 'company_Name', label: 'ក្រុមហ៑ុន', width: 150 },
          { prop: 'branch_Name', label: 'សាខា', width: 150 },
          { prop: 'invoice_date', label: 'កាលបរិច្ឆេទ', width: 120 },
          { prop: 'due_date', label: 'ថ្ងៃកំណត់', width: 120 },
          { slot: 'late_count', label: 'ចំនួនថ្ងៃផុតកំណត់', width: 130,align:'center' },
          { slot: 'total_amount', label: 'សរុប', width: 120 },
          { slot: 'paid_amount', label: 'បានបង់', width: 120 },
          { slot: 'outstanding_amount', label: 'នៅជំពាក់', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 130 },
          { label: 'ផ្សេងៗ', prop: 'cancel_reason', minWidth: 150 },
        ]"
      >
      <template #late_count="{row}">
        <el-text tag="b" style="color: crimson;">{{ row.late_count }} ថ្ងៃ</el-text>
      </template>
        <template #total_amount="{ row }">
          <el-text tag="b" type="primary"
            >{{ row.total_amount }}
            <el-text size="small" type="primary">{{
              row.currency_code
            }}</el-text></el-text
          >
        </template>
        <template #paid_amount="{ row }">
          <el-text tag="b" type="success"
            >{{ row.paid_amount }}
            <el-text size="small" type="success">{{
              row.currency_code
            }}</el-text></el-text
          >
        </template>
        <template #outstanding_amount="{ row }">
          <el-text
            tag="b"
            :type="row.outstanding_amount > 0 ? 'danger' : 'info'"
            >{{ row.outstanding_amount }}
            <el-text size="small" type="danger">{{
              row.currency_code
            }}</el-text></el-text
          >
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
            size="defualt"
          >
            {{ getStatusLabel(row.status) }}
          </el-text>
        </template>

        <template #actions="{row}">
          <el-tooltip content="មើលលំអិត" placement="top">
            <AppButton
             
              size="small"
              icon="View"
              type="success"
              circle
              @click="openDetail(row)"
            />
          </el-tooltip>
        </template>


        <template #expand="{ row: item }">
          <el-divider content-position="left">
            <el-text> លំអិតវិក័យបត្រ </el-text>
          </el-divider>
          <AppTable
            expandable
            :data="item.invoice_item"
            :columns="invoiceitemcolumns"
            :show-pagination="false"
          >
          <template #unit_price="{row:item}">
            <el-text>{{ item.unit_price }}{{ item.currency_code }}</el-text>
          </template>
          <template #discount_amount="{row:item}">
            <el-text>{{ item.discount_amount }}{{ item.currency_code }}</el-text>
          </template>
          <template #subtotal="{row:item}">
            <el-text>{{ item.subtotal }}{{ item.currency_code }}</el-text>
          </template>
          </AppTable>
        </template>
      </AppTable>
    </el-card>

   
  </div>

    <AppDialog
      v-model="showDetail"
      title="លំអិតវិក័យបត្រ"
      width="40%"
      :showDefaultFooter="false"
    >
      <AppTable
        :data="selectedInvoice?.invoice_item || []"
        :columns="invoiceitemcolumns"
        :show-pagination="false"
      >
        <template #unit_price="{ row: item }">
          {{ item.unit_price }} {{ item.currency_code }}
        </template>

        <template #discount_amount="{ row: item }">
          {{ item.discount_amount }} {{ item.currency_code }}
        </template>

        <template #subtotal="{ row: item }">
          {{ item.subtotal }} {{ item.currency_code }}
        </template>
      </AppTable>
    </AppDialog>
</template>

<style scoped>
.print-only-wrapper {
  display: none;
}
@media print {
  .print-only-wrapper {
    display: block;
  }
}
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
