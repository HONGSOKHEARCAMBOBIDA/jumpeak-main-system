import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { refreshToken } from "../api/services";
import { logoutUser } from "../api/services";
export const useAuthStore = defineStore("auth", () => {
  const accessToken = ref(localStorage.getItem("access_token") || "");
  const refreshTokenVal = ref(localStorage.getItem("refresh_token") || "");
  const permission = ref(
    JSON.parse(localStorage.getItem("permissions") || "null"),
  );
  const user = ref(JSON.parse(localStorage.getItem("user") || "null"));

  const isLoggedIn = computed(() => !!accessToken.value);

  function setAuth(data) {
    accessToken.value = data.access_token;
    refreshTokenVal.value = data.refresh_token;
    user.value = { id: data.id, name: data.name };
    permission.value = data.permissions;
    localStorage.setItem("access_token", data.access_token);
    // localStorage.setItem('user', JSON.stringify({ id: data.id, name: data.name }))
    //  localStorage.setItem('permissions', JSON.stringify(data.permissions))
  }

  async function refresh() {
    //  const res = await refreshToken({ refresh_token: refreshTokenVal.value })
    const res = await refreshToken({});
    const resData = res.data.data;

    accessToken.value = resData.access_token;
    refreshTokenVal.value = resData.refresh_token;
    permission.value = resData.permissions;

    localStorage.setItem("access_token", resData.access_token);
    //   localStorage.setItem('permissions', JSON.stringify(resData.permissions))
  }

  function logout() {
    // const res = logoutUser();
    accessToken.value = "";
    refreshTokenVal.value = "";
    user.value = null;
    permission.value = null;
    localStorage.removeItem("access_token");
    
    // localStorage.removeItem('user')
    //   localStorage.removeItem('permissions')
  }

  return {
    accessToken,
    refreshTokenVal,
    user,
    isLoggedIn,
    setAuth,
    refresh,
    logout,
    permission,
  };
});
