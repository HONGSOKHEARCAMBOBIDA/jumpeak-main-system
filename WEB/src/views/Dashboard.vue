<template>
  <div class="dashboard">
    <div class="dash-grid">
      <el-card class="dash-card">
        <template #header><span class="card-title">សកម្មភាពរហ័ស</span></template>
        <div class="action-grid">
          <div class="action-btn" @click="$router.push('/users')">
            <el-icon size="26" color="#409eff">
              <UserFilled />
            </el-icon>
            <span>បុគ្គលិក</span>
          </div>
          <div class="action-btn" @click="$router.push('/attendance')">
            <el-icon size="26" color="#67c23a">
              <Clock />
            </el-icon>
            <span>មេីលវត្តមាន</span>
          </div>
          <div class="action-btn" @click="$router.push('/leavetype')">
            <el-icon size="26" color="#8E44AD">
              <List />
            </el-icon>
            <span>ប្រភេទច្បាប់ក្នុងក្រុមហ៑ុន</span>
          </div>
          <div class="action-btn" @click="$router.push('/leaverequest')">
            <el-icon size="26" color="#67C23A">
              <Promotion />
            </el-icon>
            <span>សុំច្បាប់</span>
          </div>
          <div class="action-btn" @click="$router.push('/createattendance')">
            <el-icon size="26" color="#e6a23c">
              <Camera />
            </el-icon>
            <span>កត់ត្រាវត្តមាន</span>
          </div>
          <div class="action-btn" @click="$router.push('/company')">
            <el-icon size="26" color="#f56c6c">
              <OfficeBuilding />
            </el-icon>
            <span>ក្រុមហ៑ុន</span>
          </div>
        </div>
      </el-card>

      <el-card class="dash-card">
        <template #header>
          <div class="">
            <AppFilterBar :fields="[
              { slot: 'date', span: 10 },
              { slot: 'company', span: 10 }
            ]">
              <template #date>
                <el-date-picker v-model="filters.check_date" type="date" placeholder="ជ្រេីសរេីសថ្ងៃទី"
                  value-format="YYYY-MM-DD" clearable @change="loadAttendance" style="width: 100%" size="large" />
              </template>
              <template #company>
                <el-select v-model="filters.company_id" placeholder="ក្រុមហ៑ុន" clearable style="width: 100%"
                  size="large" @change="loadAttendance">
                  <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
                </el-select>
              </template>
            </AppFilterBar>
          </div>
        </template>
        <el-card>
          <AppTable :data="recentAttendance" :loading="loading" v-model:current-page="page" v-model:page-size="pageSize"
            :total="total" @page-change="loadAttendance()" :columns="[
              { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
              { prop: 'gender_string', label: 'ភេទ', minWidth: 90 },
              { prop: 'role_name', label: 'តួនាទី', minWidth: 110 },
              { prop: 'company_name', label: 'ក្រុមហ៑ុន', minWidth: 130 },
              { prop: 'check_date', label: 'ថ្ងៃស្កែន', minWidth: 120 },
              { label: 'ស្ថានភាព', slot: 'status', width: 100 },
            ]">
            <template #status="{ row }">
              <el-tag :type="row.status === 'COMPLETE' ? 'success' : 'warning'" size="small">
                {{ row.status }}
              </el-tag>
            </template>

            <template #actions="{ row }">
              <div class="detail-action">
                <el-badge :value="row.attendance_record?.length || 0"
                  :type="row.status === 'COMPLETE' ? 'success' : 'warning'">
                  <el-button size="large" icon="List" circle @click="viewRecords(row)" />
                </el-badge>
              </div>
            </template>
          </AppTable>
        </el-card>

        <AppDialog v-model="recordsDialog" :title="`លម្អិត — ${selectedRow?.name}`" width="60%">
          <AppTable :data="selectedRow?.attendance_record || []" :total="selectedRow?.attendance_record.length || []"
            :columns="[
              { prop: 'day_string', label: 'ថ្ងៃ', width: 100 },
              { prop: 'type_string', label: 'ប្រភេទ', width: 150 },
              { prop: 'check_time', label: 'ម៉ោងបានស្កែន', width: 150 },
              { prop: 'scheduled_time', label: 'ម៉ោងត្រូវស្កែន', width: 110 },
              { prop: 'time_diff', label: 'យឺត/មុនម៉ោង', width: 150 },
              { label: 'ស្ថានភាព', slot: 'attendance_type_name', width: 100 },
              { label: 'ទីតាំងស្កែន', slot: 'check_location' },
              { label: 'មេីលទីតាំង', slot: 'view_location' },
              { prop: 'reason', label: 'មូលហេតុ', width: 120 },
            ]">
            <template #attendance_type_name="{ row }">
              <el-tag size="small" :type="getAttendTypeTag(row.attendance_type)">
                {{ row.attendance_type_name || "—" }}
              </el-tag>
            </template>
            <template #check_location="{ row }">
              <el-icon size="large" :color="row.inzone ? '#67c23a' : '#f56c6c'">
                <component :is="row.inzone ? 'CircleCheckFilled' : 'CircleCloseFilled'" />
              </el-icon>
            </template>
            <template #view_location="{ row }">
              <el-icon size="large" :color="row.inzone ? '#67c23a' : '#f56c6c'" style="cursor: pointer"
                @click="openLocation(row)">
                <LocationInformation />
              </el-icon>
            </template>
          </AppTable>
        </AppDialog>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, } from "vue";
import { getAttendance, viewcompanyscan } from "../api/services";
import { ElNotification } from "element-plus";
import AppTable from "../../components/AppTable.vue";
import AppDialog from "../../components/AppDialog.vue";
import { LocationInformation } from "@element-plus/icons-vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
const loading = ref(false);
const recentAttendance = ref([]);
const recordsDialog = ref(false);
const selectedRow = ref(null);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const companys = ref([]);
const filters = reactive({
  check_date: null,
  company_id: null
});

async function loadAttendance() {
  loading.value = true;

  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
      check_date: filters.check_date,
      company_id: filters.company_id
    };
    const res = await getAttendance(params);
    recentAttendance.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    ElNotification({
      title: "មានបញ្ហា",
      message: e.response?.data?.error || e.message,
    });
  } finally {
    loading.value = false;
  }
}

async function fetchmanagecompany() {
  loading.value = true;
  try {
    const res = await viewcompanyscan()
    companys.value = res.data.data || [];
  } catch (e) {
    ElMessage.error("Failed to load managecompany");
  } finally {
    loading.value = false;
  }
}

function openLocation(row) {
  if (!row.latitdude || !row.longitude) {
    return ElMessage.warning("មិនមានទីតាំងស្កែន");
  }
  const url = `https://www.google.com/maps?q=${row.latitdude},${row.longitude}`;
  window.open(url, "_blank");
}
function getAttendTypeTag(type) {
  const map = { 1: "success", 2: "info", 3: "warning", 4: "danger" };
  return map[type] || "info";
}
function viewRecords(row) {
  selectedRow.value = row;
  recordsDialog.value = true;
}
onMounted(async () => {
  // loadAttendance();
  fetchmanagecompany()

});
</script>

<style scoped>
/* ── Stats grid: 4-col → 2-col ───────────────── */
.card-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
}

.stat-label {
  font-size: 12px;
  color: #888;
  margin-top: 2px;
}

/* ── Dashboard two-col grid ──────────────────── */
.dash-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 16px;
}

.dash-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.card-title {
  font-weight: 600;
  font-size: 14px;
}

/* ── Quick actions ────────────────────────────── */
.action-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 8px;
  border-radius: 10px;
  border: 1px solid #ebeef5;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 12px;
  color: #555;
  text-align: center;
}

.action-btn:hover {
  background: #f5f7fa;
  border-color: #409eff;
  color: #409eff;
  transform: translateY(-2px);
}

/* ── Attendance list ──────────────────────────── */
.attend-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 0;
  border-bottom: 1px solid #f5f5f5;
}

.attend-row:last-child {
  border-bottom: none;
}

.attend-info {
  flex: 1;
  min-width: 0;
}

.attend-name {
  font-weight: 500;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.attend-meta {
  font-size: 11px;
  color: #13097b;
}

.detail-action {
  padding: 10px 0;
}

/* ── Mobile overrides ────────────────────────── */
@media (max-width: 767px) {
  .card-header-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .card-header-row .el-date-picker {
    width: 100% !important;
  }

  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .stat-card {
    padding: 12px;
    gap: 10px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
  }

  .stat-value {
    font-size: 20px;
  }

  .stat-label {
    font-size: 11px;
  }

  .dash-grid {
    grid-template-columns: 1fr;
    /* stack vertically */
    gap: 12px;
  }
}

@media (max-width: 400px) {

  /* Very small screens: action buttons tighter */
  .action-btn {
    padding: 12px 6px;
    font-size: 11px;
  }
}
</style>
