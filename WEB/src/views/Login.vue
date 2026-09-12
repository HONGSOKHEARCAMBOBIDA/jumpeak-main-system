<template>
  <div class="login-page">
    <div class="login-card">
      <AppTabs
        v-model="activeTab"
        :tabs="[
          { name: 'phone', label: 'ចូលប្រព័ន្ធ' },
        ]"
        tab-position="top"
        stretch="true"
      >
        <template #phone>
          <el-form
            :model="form"
            :rules="rules"
            ref="formRef"
            @submit.prevent="handleLogin"
            label-position="top"
          >
            <AppInput
              v-model="form.email"
              label="អុីម៉ែល"
              prop="phone"
              placeholder="បញ្ចូលអុីម៉ែល"
              prefix-icon="Message"
              clearable="true"
              autofocus="true"
            >
            </AppInput>
            <AppInput
              v-model="form.password"
              label="ពាក្យសម្ងាត់"
              prop="password"
              type="password"
              placeholder="បញ្ចូលពាក្យសម្ងាត់"
              prefix-icon="Lock"
              @enter="handleLogin"
            >
            </AppInput>
            <AppButton
              native-type="submit"
              :loading="loading"
              type="primary"
              block="false"
            >
              ចូលប្រព័ន្ធ
            </AppButton>
          </el-form>
        </template>
      </AppTabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { login, } from "../api/services";
import AppButton from "../../components/AppButton.vue";
import AppInput from "../../components/AppInput.vue";
import AppTabs from "../../components/AppTabs.vue";
import { useNotification } from "../../composables/useNotification.js";
const notify = useNotification();
const router = useRouter();
const auth = useAuthStore();
const formRef = ref();
const loading = ref(false);


const form = reactive({ 
  email: "", 
  password: ""
 });
const rules = {
  email: [{ required: true, message: "សូមបញ្ចូលអុីម៉ែល", trigger: "blur" }],
  password: [
    { required: true, message: "សូមបញ្ជូលពាក្យសម្ងាត់", trigger: "blur" },
  ],
};

async function handleLogin() {
  await formRef.value.validate();
  loading.value = true;
  try {
    const res = await login(form);
    auth.setAuth(res.data.data);
    router.push("/Dashboard");
    notify.success("ចូលប្រព័ន្ធបានជោគជ័យ")
  } catch (e) {
    notify.error(e.response?.data?.message || "Login failed");
    console.log(e.response?.data?.message)
  } finally {
    loading.value = false;
  }
}

</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  border: 1px solid #8aaff5;
  border-radius: 2px;
  padding: 48px 40px;
  width: 420px;
}

.login-btn {
  width: 100%;
  margin-top: 8px;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
}

.qr-scanner-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 0 16px;
  min-height: 200px;
}

.qr-start,
.qr-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 24px 0;
}

.qr-viewport {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

#qr-reader {
  width: 100% !important;
  max-width: 280px;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid #409eff;
}

#qr-reader__scan_region img,
#qr-reader__dashboard {
  display: none !important;
}

.qr-hint {
  color: #888;
  font-size: 13px;
  margin: 10px 0 0;
  text-align: center;
}

@media (max-width: 768px) {
  .login-page {
    padding: 25px;
  }

  .login-card {
    width: 100%;
    max-width: 420px;
    padding: 24px 20px;
    border-radius: 25px;
  }
}
</style>
