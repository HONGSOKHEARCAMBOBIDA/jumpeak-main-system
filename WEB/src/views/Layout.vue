<template>
  <div class="layout" :style="{ '--brand-color': companyStore.color || '#1a1a2e' }">
    <aside class="sidebar" :class="{ collapsed }">
      <div class="sidebar-header">
        <el-icon size="26" color="#fff"><UserFilled /></el-icon>
        <span v-if="!collapsed" class="brand">ប្រព័ន្ធគ្រប់គ្រងធនធានមនុស្ស</span>
      </div>
      <el-menu :default-active="$route.path" router :collapse="collapsed"
        background-color="var(--brand-color)" text-color="#a0a8c0" active-text-color="#fff">
        <el-menu-item v-for="n in nav" :key="n.path" :index="n.path">
          <el-icon><component :is="n.icon" /></el-icon>
          <template #title>{{ n.label }}</template>
        </el-menu-item>
      </el-menu>
    </aside>

    <div class="wrap">
      <header class="topbar">
        <div class="left">
          <el-button class="d-only" :icon="collapsed ? 'Expand' : 'Fold'" circle plain @click="collapsed = !collapsed" />
          <b>{{ pageTitle }}</b>
        </div>
        <el-dropdown @command="handleCommand">
          <div class="avatar-row">
            <el-avatar :size="32" style="background:#409eff">{{ userName[0]?.toUpperCase() }}</el-avatar>
            <span class="d-only" style="font-size:14px">{{ userName }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">ប្រវត្តរូប</el-dropdown-item>
              <el-dropdown-item command="logout" divided>ចាកចេញ</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </header>
      <main class="content"><router-view /></main>
    </div>

    <!-- Mobile bottom nav -->
    <nav class="bnav">
      <router-link v-for="n in nav" :key="n.path" :to="n.path" class="bitem" :class="{ active: $route.path === n.path }">
        <el-icon size="20"><component :is="n.icon" /></el-icon>
        <span>{{ n.short }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCompanyStore } from '../stores/company'
import { useUserDataStore } from '../stores/user_data'
import { ElMessage } from 'element-plus'
const companyStore = useCompanyStore()
const loading = ref(false)
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const userdataStore = useUserDataStore()
const userName = computed(() => userdataStore.name || 'User')
const collapsed = ref(false)
const nav = computed(() => {
  return router.getRoutes()
    .filter(r => r.meta?.showInNav)
    .filter(r => !r.meta?.permission || userdataStore.permissions?.some(p => p.name === r.meta.permission))
    .map(r => ({
      path: r.path, // getRoutes() returns full resolved paths already
      icon: r.meta.icon,
      label: r.meta.title,
      short: r.meta.short
    }))
})
// const pageTitle = computed(() => route.meta.title(n => n.path === route.path)?.label || 'HR')
const pageTitle = computed(() => route.meta.title || 'HR')
function handleCommand(cmd) {
  cmd === 'logout' ? (auth.logout(), router.push('/login')) : router.push('/profile')
}
onMounted(() => {
  companyStore.fetchColor()
  userdataStore.getuserdata()
})
</script>

<style scoped>
.layout { height: 100vh; display: flex; overflow: hidden; }


/* Sidebar */
.sidebar { width: 260px; background: var(--brand-color, #1a1a2e); display: flex; flex-direction: column; transition: width .3s; overflow: hidden; flex-shrink: 0; }
.sidebar.collapsed { width: 64px; }
.sidebar-header { display: flex; align-items: center; gap: 10px; padding: 18px 14px; border-bottom: 1px solid rgba(255,255,255,.08); }
.brand { color: #fff; font-size: 15px; font-weight: 700; white-space: nowrap; }
:deep(.el-menu) { border: none; }
:deep(.el-menu-item) { margin: 2px 8px; border-radius: 6px; }
:deep(.el-menu-item.is-active) { background: rgba(64,158,255,.2) !important; color: #fff !important; }

/* Main */
.wrap { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.topbar { height: 54px; background: #fff; border-bottom: 1px solid #ebeef5; display: flex; align-items: center; justify-content: space-between; padding: 0 20px; flex-shrink: 0; box-shadow: 0 1px 4px rgba(0,0,0,.06); }
.left { display: flex; align-items: center; gap: 12px; }
.avatar-row { display: flex; align-items: center; gap: 8px; cursor: pointer; padding: 4px 8px; border-radius: 8px; }
.avatar-row:hover { background: #f5f7fa; }
.content { flex: 1; overflow-y: auto; padding: 20px; background: #f0f2f5; }

/* Bottom nav — hidden on desktop */
.bnav { display: none; }

/* Mobile */
@media (max-width: 768px) {
  .sidebar, .d-only { display: none !important; }

  .content { padding: 12px; padding-bottom: 70px; }
  .topbar { padding: 0 14px; height: 50px; }

  .bnav {
    display: flex;
    position: fixed; bottom: 0; left: 0; right: 0;
    height: 58px; background: #fff;
    border-top: 1px solid #ebeef5;
    box-shadow: 0 -2px 8px rgba(0,0,0,.07);
    z-index: 999;
  }
  .bitem {
    flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center;
    gap: 2px; font-size: 10px; color: #aaa; text-decoration: none;
  }
  .bitem.active { color: #409eff; }
}
</style>
