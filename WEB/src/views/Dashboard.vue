<template>
  <div class="dashboard">
    <div class="dash-grid">
      <el-card class="dash-card">
        <template #header
          ><span class="card-title">សកម្មភាពរហ័ស</span></template
        >
        <div class="action-grid">
          <div class="action-btn" @click="$router.push('/customer')">
            <el-icon size="26" color="#409eff">
              <UserFilled />
            </el-icon>
            <span>អតិថិជន</span>
          </div>
          <div class="action-btn" @click="$router.push('/product')">
            <el-icon size="26" color="#67c23a">
              <Clock />
            </el-icon>
            <span>ឥវ៉ាន់ទាំងអស់</span>
          </div>
          <div class="action-btn" @click="$router.push('/invoice')">
            <el-icon size="26" color="#8E44AD">
              <List />
            </el-icon>
            <span>កត់អ្នកជំពាក់</span>
          </div>
          <div class="action-btn" @click="$router.push('/leaverequest')">
            <el-icon size="26" color="#67C23A">
              <Promotion />
            </el-icon>
            <span>សុំច្បាប់</span>
          </div>
          <div class="action-btn" @click="$router.push('/company')">
            <el-icon size="26" color="#e6a23c">
              <School />
            </el-icon>
            <span>ក្រុមហ៑ុន</span>
          </div>
          <div class="action-btn" @click="$router.push('/company')">
            <el-icon size="26" color="#f56c6c">
              <OfficeBuilding />
            </el-icon>
            <span>សាខាតាមក្រុមហ៑ុន</span>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from "vue";
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
