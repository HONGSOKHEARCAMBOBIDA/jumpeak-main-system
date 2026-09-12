<script setup>
import { ref, reactive, onMounted, watch } from "vue";
import {
  getcustomeroutstandingreport,
  getcompanynopagitaion,
  getbranchnopagination,
} from "../api/services.js";
import AppTable from "../../components/AppTable.vue";
import AppSelect from "../../components/AppSelect.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import { useNotification } from "../../composables/useNotification.js";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";

const notify = useNotification();

const report = ref([]);
const loading = ref(false);

const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const showDetail = ref(false);
const selected = ref(null);

function openDetail(row) {
  selected.value = row;
  showDetail.value = true;
}

const filters = reactive({
  company_id: null,
  branch_id: null,
});

const companyOptions = ref([]);
const branchOptions = ref([]);

const branchColumns = [
  { prop: "Branch_name", label: "សាខា", minWidth: 150 },
  { slot: "branch_total", label: "នៅជំពាក់", minWidth: 150 },
];

async function loadCompanyOptions() {
  try {
    const res = await getcompanynopagitaion();
    companyOptions.value = (res.data.data || []).map((c) => ({
      label: c.name,
      value: c.id,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.error || "Failed to load companies");
  }
}

async function loadBranchOption(companyID) {
  if(!companyID) return []
  try {

    const res = await getbranchnopagination(companyID)
    return (res.data.data || []).map((f) => ({
      label: f.name,
      value: f.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load branch')
    return []
  }
}

async function fetchReport() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.company_id) params.company_id = filters.company_id;
    if (filters.branch_id) params.branch_id = filters.branch_id;

    const res = await getcustomeroutstandingreport(params);
    report.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e?.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

watch(
  () => filters.company_id,
  async (newVal) => {
    filters.branch_id = null;
    branchOptions.value = [];
    branchOptions.value = await loadBranchOption(newVal);
    page.value = 1;
    fetchReport();
  }
);

watch(
  () => filters.branch_id,
  async () => {
    page.value = 1;
    fetchReport();
  }
);


onMounted(() => {
  loadCompanyOptions();
  fetchReport();
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'company', span: 8 },
        { slot: 'branch', span: 8 },
      ]"
    >
      <template #company>
        <AppSelect
          v-model="filters.company_id"
          :options="companyOptions"
          placeholder="ក្រុមហ៊ុន"
          size="large"
          filterable
          clearable
        />
      </template>
      <template #branch>
        <AppSelect
          v-model="filters.branch_id"
          :options="branchOptions"
          placeholder="សាខា"
          size="large"
          filterable
          clearable
        />
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        expandable
        :data="report"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchReport"
        :columns="[
          { prop: 'company_name', label: 'ក្រុមហ៊ុន', minWidth: 180 },
          { slot: 'total_amount', label: 'នៅជំពាក់សរុប', minWidth: 160 },
        ]"
      >
        <template #total_amount="{ row }">
          <el-text
            tag="b"
            :type="row.total_amount > 0 ? 'danger' : 'info'"
          >
            {{ row.total_amount }}
            <el-text size="small" type="danger">{{ row.currency }}</el-text>
          </el-text>
        </template>

        <template #actions="{row}">
        <el-tooltip content="មើលលំអិត" placement="top">
            <AppButton
              @click="openDetail(row)"
              size="small"
              icon="View"
              type="success"
              circle
             
            />
          </el-tooltip>
        </template>

        <template #expand="{ row }">
          <el-divider content-position="left">
            <el-text>លម្អិតតាមសាខា</el-text>
          </el-divider>
          <AppTable
            show-index
            :data="row.branch_outstanding"
            :columns="branchColumns"
            :show-pagination="false"
          >
            <template #branch_total="{ row: b }">
              <el-text
                tag="b"
                :type="b.total_amount > 0 ? 'danger' : 'info'"
              >
                {{ b.total_amount }}
                <el-text size="small" type="danger">{{ b.currency }}</el-text>
              </el-text>
            </template>
          </AppTable>
        </template>
      </AppTable>
    </el-card>
  </div>

      <AppDialog
      v-model="showDetail"
      title="លំអិតតាមសាខា"
      width="40%"
      :showDefaultFooter="false"
    >
          <AppTable
            show-index
            :data="selected?.branch_outstanding || []"
            :columns="branchColumns"
            :show-pagination="false"
          >
            <template #branch_total="{ row: b }">
              <el-text
                tag="b"
                :type="b.total_amount > 0 ? 'danger' : 'info'"
              >
                {{ b.total_amount }}
                <el-text size="small" type="danger">{{ b.currency }}</el-text>
              </el-text>
            </template>
          </AppTable>
    </AppDialog>
</template>

<style scoped>
.table-card {
  border-radius: 6px;
}
</style>