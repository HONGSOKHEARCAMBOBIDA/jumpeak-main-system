<template>
  <div class="profile-wrap">
    <el-row :gutter="20">
      <!-- Profile Card -->
      <el-col :xs="24" :sm="24" :md="8" :lg="8">
        <el-card class="profile-card">
          <div class="profile-avatar-section">
            <el-avatar :size="90" style="background: linear-gradient(135deg, #667eea, #764ba2); font-size:32px">
              {{ userName?.charAt(0)?.toUpperCase() || 'U' }}
            </el-avatar>
            <h2>{{ userName }}</h2>
            <p class="profile-id">User ID: {{ userdataStore.userdata?.id }}</p>
          </div>
        </el-card>
      </el-col>

      <!-- Change Password -->
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-card class="pw-card">
          <template #header><span class="card-title">ប្ដូរពាក្យសម្ងាត់</span></template>
          <el-form
            :model="pwForm"
            :rules="pwRules"
            ref="pwFormRef"
            :label-width="isMobile ? '0' : '160px'"
            :label-position="isMobile ? 'top' : 'right'"
            class="pw-form"
          >
            <el-form-item label="ពាក្យសម្ងាត់ថ្មី" prop="new_password">
              <el-input v-model="pwForm.new_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="បញ្ជាក់ពាក្យសម្ងាត់ថ្មី" prop="confirm_password">
              <el-input v-model="pwForm.confirm_password" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="changingPw" class="submit-btn" @click="handleChangePassword">
                កែប្រែ
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserDataStore } from '../stores/user_data'
import { changePassword } from '../api/services'


const pwFormRef = ref()
const changingPw = ref(false)
const userdataStore = useUserDataStore()
// Track viewport for label position / width switching
const userName = computed(() => userdataStore.name || 'User')
const isMobile = ref(window.innerWidth < 768)
function handleResize() {
  isMobile.value = window.innerWidth < 768
}
onMounted(() => window.addEventListener('resize', handleResize))
onUnmounted(() => window.removeEventListener('resize', handleResize))

const pwForm = reactive({ new_password: '', confirm_password: '' })
const pwRules = {
  new_password: [
    { required: true, message: 'Password is required' },
    { min: 6, message: 'Password must be at least 6 characters' }
  ],
  confirm_password: [
    { required: true, message: 'Please confirm password' },
    {
      validator: (_, value, callback) => {
        if (value !== pwForm.new_password) callback(new Error('Passwords do not match'))
        else callback()
      }
    }
  ]
}

async function handleChangePassword() {
  await pwFormRef.value.validate()
  changingPw.value = true
  try {
    await changePassword({ new_password: pwForm.new_password })
    ElMessage.success('Password changed successfully')
    pwForm.new_password = ''
    pwForm.confirm_password = ''
    pwFormRef.value.resetFields()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || 'Failed to change password')
  } finally {
    changingPw.value = false
  }
}
</script>

<style scoped>
.profile-wrap {
  padding: 0 4px;
}

.profile-card {
  border-radius: 6px;
  margin-bottom: 20px;
}

.profile-avatar-section {
  text-align: center;
  padding: 20px 0;
}

.profile-avatar-section h2 {
  margin-top: 16px;
  font-size: 20px;
  font-weight: 700;
}

.profile-id {
  color: #888;
  font-size: 13px;
  margin-top: 4px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}

.pw-form {
  max-width: 460px;
  width: 100%;
}

.submit-btn {
  width: 100%;
}

/* ---- Mobile breakpoints ---- */
@media (max-width: 768px) {
  .pw-card {
    margin-top: 0;
  }

  .pw-form {
    max-width: 100%;
  }

  .profile-avatar-section {
    padding: 14px 0;
  }

  .profile-avatar-section h2 {
    font-size: 18px;
  }

  .submit-btn {
    width: 100%;
    height: 42px;
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .profile-wrap {
    padding: 0 2px;
  }

  .card-title {
    font-size: 14px;
  }
}
</style>