<template>
  <div class="login-page">
    <div class="login-card">
      <AppTabs
        v-model="activeTab"
        :tabs="[
          { name: 'phone', label: 'ចូលតាមរយៈលេខទូរសព្ទ' },
          { name: 'qr', label: 'ចូលតាមរយៈQR' },
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
              v-model="form.phone"
              label="លេខទូរសព្ទ"
              prop="phone"
              placeholder="បញ្ចូលលេខទូរសព្ទ"
              prefix-icon="Phone"
              clearable="true"
              autofocus="true"
              type="number"
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
        <template #qr>
          <div class="qr-scanner-wrap">
            <div v-if="!scanning && !loading" class="qr-start">
              <p class="qr-hint">ចុចប៊ូតុងខាងក្រោម ដើម្បីស្កែន QR Code</p>
              <AppButton type="primary" size="large" @click="startScanner">
                បើកកាមេរ៉ាស្កែន
              </AppButton>
            </div>
            <div v-show="scanning" class="qr-viewport">
              <div id="qr-reader"></div>
              <p class="qr-hint">ដាក់ QR Code នៅក្នុងប្រអប់ស្កែន</p>
              <AppButton type="primary" size="large" @click="stopScanner">
                បោះបង់
              </AppButton>
            </div>

            <div v-if="loading" class="qr-loading">
              <el-icon class="is-loading" size="32"><Loading /></el-icon>
              <p>កំពុងចូលប្រព័ន្ធ...</p>
            </div>
          </div>
        </template>
      </AppTabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onUnmounted, nextTick } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useAuthStore } from "../stores/auth";
import { login, loginByQr } from "../api/services";
import { Html5Qrcode } from "html5-qrcode";
import AppButton from "../../components/AppButton.vue";
import AppInput from "../../components/AppInput.vue";
import AppTabs from "../../components/AppTabs.vue";
import { useNotification } from "../../composables/useNotification.js";
const notify = useNotification();
const router = useRouter();
const auth = useAuthStore();
const formRef = ref();
const loading = ref(false);
const activeTab = ref("phone");
const scanning = ref(false);
let html5QrCode = null;

async function startScanner() {
  scanning.value = true;
  await nextTick(); // wait for #qr-reader to mount

  html5QrCode = new Html5Qrcode("qr-reader");
  try {
    await html5QrCode.start(
      { facingMode: "environment" }, // rear camera
      { fps: 10, qrbox: { width: 180, height: 180 } },
      (decodedText) => {
        stopScanner();
        handleQrLogin(decodedText); // auto-login on scan
      },
      () => {}, // ignore per-frame errors
    );
  } catch (e) {
    ElMessage.error("មិនអាចបើកកាមេរ៉ាបាន: " + e);
    scanning.value = false;
  }
}

function stopScanner() {
  scanning.value = false;
  html5QrCode?.stop().catch(() => {});
  html5QrCode = null;
}

const form = reactive({ phone: "", password: "" });
const rules = {
  phone: [{ required: true, message: "សូមបញ្ចូលលេខទូរសព្ទ", trigger: "blur" }],
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
  } finally {
    loading.value = false;
  }
}

async function handleQrLogin(token) {
  if (!token) return ElMessage.warning("QR token មិនត្រឹមត្រូវ");
  loading.value = true;
  try {
    const res = await loginByQr({ qr_token: token });
    auth.setAuth(res.data.data);
    router.push("/dashboard");
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "QR login failed");
  } finally {
    loading.value = false;
  }
}

onUnmounted(() => stopScanner());
watch(activeTab, (tab) => {
  if (tab !== "qr") stopScanner();
});
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
