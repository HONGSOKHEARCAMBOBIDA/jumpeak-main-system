<template>
  <div class="bk-wrap">
    <div class="bk-header">
      <div>
        <h1 class="bk-title">Database backup</h1>
      </div>
      <el-button v-if="canBackup" type="primary" :loading="triggering" :icon="Upload" @click="handleTriggerBackup">
        Back up now
      </el-button>
    </div>

    <el-alert
      v-if="alert.message"
      :title="alert.message"
      :type="alert.type"
      show-icon
      closable
      class="bk-alert"
      @close="alert.message = ''"
    />

    <div class="bk-stats">
      <div class="bk-stat">
        <p class="bk-stat-label">Total backups</p>
        <p class="bk-stat-value">{{ backups.length }}</p>
        <p class="bk-stat-sub">Stored files</p>
      </div>
      <div class="bk-stat">
        <p class="bk-stat-label">Latest backup</p>
        <p class="bk-stat-value">{{ latestBackupLabel }}</p>
        <p class="bk-stat-sub">Last run</p>
      </div>
      <div class="bk-stat">
        <p class="bk-stat-label">Total size</p>
        <p class="bk-stat-value">{{ totalSize }}</p>
        <p class="bk-stat-sub">All files</p>
      </div>
    </div>

    <div class="bk-card">
      <div class="bk-card-header">
        <span class="bk-card-title">
          <span class="bk-status-dot" />
          Backup files
        </span>
        <el-button link type="primary" :icon="Refresh" @click="loadBackups">Refresh</el-button>
      </div>

      <el-table
        :data="backups"
        v-loading="loading"
        style="width: 100%"
        stripe
      >
        <el-table-column label="Filename" min-width="260">
          <template #default="{ row }">
            <span class="bk-filename">{{ row.filename }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Size" width="140">
          <template #default="{ row }">
            <el-tag type="warning" size="small" effect="plain">{{ formatSize(row.size_bytes) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="Created" width="160">
          <template #default="{ row }">
            <span class="bk-date">{{ formatRelative(row.created_at) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Action" width="140" align="right">
          <template #default="{ row }">
            <!-- <el-button
              link
              type="primary"
              size="large"
              :icon="Download"
              @click="handleDownload(row.filename)"
            >
              Download
            </el-button> -->
            <el-button v-if="canDownloadBackup" type="success" :icon="Download" circle @click="handleDownload(row.filename)" />
            <el-button v-if="canDeleteBackup" type="danger" :icon="Delete" circle @click="handleDelete(row.filename)"/>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Upload, Download, Refresh,Delete } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'
import { triggerBackup, listBackups, downloadBackup,deleteBackup } from '../api/services'
import { useUserDataStore } from '../stores/user_data'
const backups = ref([])
const loading = ref(false)
const triggering = ref(false)
const alert = ref({ type: 'success', message: '' })
const userdataStore = useUserDataStore()
const canBackup = computed(() => userdataStore.permissions?.some(p => p.name === 'add.backup'))
const canDownloadBackup = computed(() => userdataStore.permissions?.some(p => p.name === 'view.download.backup'))
const canDeleteBackup = computed(() => userdataStore.permissions?.some(p => p.name === 'delete.backup'))
const showAlert = (type, message) => {
  alert.value = { type, message }
  setTimeout(() => { alert.value.message = '' }, 4000)
}

const latestBackupLabel = computed(() => {
  if (!backups.value.length) return '—'
  return formatRelative(backups.value[0].created_at)
})

const totalSize = computed(() => {
  const total = backups.value.reduce((sum, b) => sum + (b.size_bytes ?? 0), 0)
  return formatSize(total)
})

const loadBackups = async () => {
  loading.value = true
  try {
    const res = await listBackups()
    backups.value = res.data.backups ?? []
  } catch {
    ElNotification({ title: 'Warning', message: 'ទាញទិន្ន័យមិនបាន', type: 'warning' })
  } finally {
    loading.value = false
  }
}

const handleTriggerBackup = async () => {
  triggering.value = true
  try {
    await triggerBackup()
    ElNotification({ title: 'Success', message: 'Back Up ទិន្ន័យបានជោគជ័យ', type: 'success' })
    await loadBackups()
  } catch {
    ElNotification({ title: 'Error', message: 'Back Up មិនជោគជ័យ', type: 'error' })
  } finally {
    triggering.value = false
  }
}


const handleDownload = async (filename) => {
  try {
    const res = await downloadBackup(filename)
    const url = URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    ElNotification({ title: 'Error', message: 'Download fail', type: 'error' })
  }
}

const handleDelete = async (filename) => {
  try {
    const res = await deleteBackup(filename)
     ElNotification({ title: 'Success', message: 'Delete Back Up Success', type: 'success' })
      await loadBackups()
  } catch {
    ElNotification({ title: 'Error', message: 'Download fail', type: 'error' })
  }
}


const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

const formatRelative = (iso) => {
  const diff = Math.floor((Date.now() - new Date(iso)) / 1000)
  if (diff < 60) return diff + 's ago'
  if (diff < 3600) return Math.floor(diff / 60) + 'm ago'
  if (diff < 86400) return Math.floor(diff / 3600) + 'h ago'
  return Math.floor(diff / 86400) + 'd ago'
}

onMounted(loadBackups)
</script>

<style scoped>
.bk-wrap {
  padding: 1.5rem;
}

.bk-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.bk-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 4px;
}

.bk-subtitle {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.bk-alert {
  margin-bottom: 1rem;
}

/* Stats */
.bk-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 1.25rem;
}

.bk-stat {
  border: 1px solid #409EFF;;
  border-radius: 6px;
  padding: 14px 16px;
}

.bk-stat-label {
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--el-text-color-secondary);
  margin: 0 0 4px;
}

.bk-stat-value {
  font-size: 22px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
  line-height: 1.2;
}

.bk-stat-sub {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin: 2px 0 0;
}

/* Card */
.bk-card {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  overflow: hidden;
}

.bk-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  background: var(--el-fill-color-lighter);
}

.bk-card-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.bk-status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--el-color-success);
}

/* Table overrides */
.bk-filename {
  font-family: monospace;
  font-size: 16px;
  color: var(--el-text-color-primary);
}

.bk-date {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

:deep(.el-table) {
  border-radius: 0;
}

:deep(.el-table th.el-table__cell) {
  background: var(--el-fill-color-lighter);
  font-size: 12px;
  font-weight: 500;
  color: var(--el-text-color-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

:deep(.el-table td.el-table__cell) {
  padding: 10px 12px;
}
</style>