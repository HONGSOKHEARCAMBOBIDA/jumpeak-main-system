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
        path: "company",
        name: "Company",
        component: () => import("../views/Company.vue"),
        meta: {
          title: "ក្រុមហ៑ុន",
          icon: "OfficeBuilding",
          short: "ក្រុមហ៑ុន",
          showInNav: true,
        },
      },
      {
        path: "users",
        name: "Users",
        component: () => import("../views/Users.vue"),
        meta: {
          title: "បុគ្គលិក",
          icon: "User",
          short: "បុគ្គលិក",
          showInNav: true,
        },
      },
      {
        path: "leavetype",
        name: "LeaveType",
        component: () => import("../views/LeaveType.vue"),
        meta: {
          title: "ប្រភេទច្បាប់",
          icon: "List",
          short: "ប្រភេទច្បាប់",
          showInNav: false,
        }
      },
      {
        path: "leaverequest",
        name: "LeaveRequest",
        component: () => import("../views/LeaveRequest.vue"),
        meta: {
          title: "សុំច្បាប់",
          icon: "Promotion",
          short: "សុំច្បាប់",
          showInNav: false,
        }
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
        path: "attendance",
        name: "Attendance",
        component: () => import("../views/Attendance.vue"),
        meta: {
          title: "វត្តមាន",
          icon: "Clock",
          short: "វត្តមាន",
          showInNav: true,
        },
      },
      {
        path: "payroll",
        name: "Payroll",
        component: () => import("../views/Payroll.vue"),
        meta: {
          title: "បេីកប្រាក់ខែ",
          icon: "Money",
          short: "បេីកប្រាក់ខែ",
          permission: "add.payroll",
          showInNav: true,
        },
      },
      {
        path: "payrolllist",
        name: "Payrolllist",
        component: () => import("../views/PayrollList.vue"),
        meta: {
          title: "របាយការណ៍បើកប្រាក់ខែ",
          icon: "Money",
          short: "របាយការណ៍បើកប្រាក់ខែ",
          showInNav: true,
        },
      },
      {
        path: "createattendance",
        name: "CreateAttendance",
        component: () => import("../views/CreateAttendance.vue"),
        meta: { title: "បង្កេីតវត្តមាន",showInNav: false }, // no showInNav -> hidden from sidebar
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
