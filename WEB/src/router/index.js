import { createRouter, createWebHistory } from "vue-router";
import { useUserDataStore } from '../stores/user_data'
import { useAuthStore } from "../stores/auth.js";
const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/Login.vue"),
    meta: { public: true },
  },
  {
    path: "/",
    component: () => import("../views/Layout.vue"),
    redirect: "/Dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("../views/Dashboard.vue"),
        meta: {
          title: "សង្ខែបទិន្ន័យ",
          icon: "Odometer",
          short: "ទិន្ន័យ",
          showInNav: true,
        },
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import("../views/Role.vue"),
        meta:{
          title: "សិទ្ធ",
          icon: "User",
          short: "សិទ្ធ",
          showInNav: true,
          permission: "add.role.has.permission",
        }
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("../views/Profile.vue"),
        meta: {
          title: "ប្រវត្តរូប",
          icon: "Setting",
          short: "ប្រវត្ត",
          showInNav: true,
        },
      },
      {
        path: "backup",
        name: "Backup",
        component: () => import("../views/Backup.vue"),
        meta: {
          title: "Backup",
          icon: "Download",
          short: "Backup",
          permission: "view.backup",
          showInNav: true,
        },
      },
      {
        path: "company",
        name: "Company",
        component: () => import("../views/company.vue"),
      },
      {
        path: "customer",
        name: "Customer",
        component: () => import("../views/customer.vue"),
      },
      {
        path: "product",
        name: "Product",
        component: () => import("../views/product.vue"),
      },
      {
        path: "Invoice",
        name: "Invoice",
        component: () => import("../views/Invoice.vue"),
      },
      {
        path: "/:pathMatch(.*)*",
        name: "NotFound",
        component: () => import("../views/NotFound.vue"),
      },
    ],
  },
];

const router = createRouter({
  linkActiveClass: "font-bold",
  linkExactActiveClass: "font-bold",
  history: createWebHistory(),
  routes,
});

router.beforeEach((to) => {
  const userDataStore = useUserDataStore();
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: "Login" };
  }
  if (to.name === "Login" && auth.isLoggedIn) {
    return { name: "Dashboard" };
  }
  if (to.meta.permission) {
    const hasPermission = userDataStore.permissions?.some(
      (p) => p.name === to.meta.permission,
    );
    if (!hasPermission) {
      return { name: "Dashboard" };
    }
  }
  if (to.meta.title) {
    document.title = to.meta.title;
  }
});

export default router;
