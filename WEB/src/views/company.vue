<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import { getcompany, addcompany, updatecompany } from "../api/services.js";
import { useAuthStore } from "../stores/auth";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
import AppSelect from "../../components/AppSelect.vue";

const notify = useNotification();
const userDataStore = useUserDataStore();
const companies = ref([]);
const loading = ref(false);
const useloading = useLoading();
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const dialogVisible = ref(false);
const isEdit = ref(false);
const editId = ref(null);
const formRef = ref();
const auth = useAuthStore();

const CurrencyOption = [
  { label: "USD", value: "USD" },
  { label: "KHR", value: "KHR" },
];

const StatusOption = [
  { label: "ACTIVE", value: "ACTIVE" },
  { label: "SUSPENDED", value: "SUSPENDED" },
];

const form = reactive({
  name: "",
  base_currency: "USD",
  status: "ACTIVE",
});

const branchcolumns = [
  { prop: "name", slot: "name", label: "ឈ្មោះ", minwidth: 120 },
  { prop: "address", label: "ទីតាំងសាខា", minwidth: 120 },
  { prop: "status", label: "ស្ថានភាព", width: 120 },
];

const canAddCompany = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.company"),
);

const canEditCompany = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "update.company"),
);

const rules = {
  name: [{ required: true, message: "Company name is required" }],
  base_currency: [{ required: true, message: "Currency is required" }],
  status: [{ required: true, message: "Status is required" }],
};

async function fetchCompanies() {
  loading.value = true;
  try {
    const res = await getcompany({
      page: page.value,
      page_size: pageSize.value,
    });
    companies.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  isEdit.value = false;
  form.name = "";
  form.base_currency = "USD";
  form.status = "ACTIVE";
  dialogVisible.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, {
    name: row.name || "",
    base_currency: row.base_currency || "USD",
    status: row.status || "ACTIVE",
  });
  dialogVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    const payload = {
      name: form.name,
      base_currency: form.base_currency,
      status: form.status,
    };
    if (isEdit.value) {
      await updatecompany(editId.value, payload);
      notify.success("កែប្រែបានជោគជ័យ");
    } else {
      await addcompany(payload);
      notify.success("បង្កើតក្រុមហ៊ុនបានជោគជ័យ");
    }
    dialogVisible.value = false;
    fetchCompanies();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

onMounted(fetchCompanies);
</script>

<template>
  <div>
    <div class="page-header">
      <AppButton
        v-if="canAddCompany"
        type="primary"
        @click="openCreate"
        :block="false"
        size="default"
      >
        បន្ថែមក្រុមហ៊ុន
      </AppButton>
    </div>

    <el-card class="table-card">
      <AppTable
        expandable
        :data="companies"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchCompanies"
        :columns="[
          { prop: 'name', label: 'ឈ្មោះ', minWidth: 150 },
          { prop: 'base_currency', label: 'រូបិយប័ណ្ណ', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
        ]"
      >
        <template #status="{ row }">
          <el-tag
            :type="row.status === 'ACTIVE' ? 'success' : 'danger'"
            size="small"
          >
            {{ row.status }}
          </el-tag>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
              v-if="canEditCompany"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="openEdit(row)"
            />
          </el-tooltip>
        </template>

        <template #expand="{ row }">
          <el-divider content-position="left">
            <el-text>
              សាខាសរុប {{ row.branches.length }}
              <el-tooltip content="ថែមសាខា" placement="top"
                ><AppButton
                  circle
                  size="small"
                  type="primary"
                  icon="Plus"
                ></AppButton
              ></el-tooltip>
            </el-text>
          </el-divider>
          <AppTable
            expandable
            :data="row.branches"
            :columns="branchcolumns"
            :show-pagination="false"
          >
            <template #name="{ row: branchRow }">
              <el-text>
                {{ branchRow.name }} |
                <el-text size="small" type="primary">{{
                  branchRow.code
                }}</el-text>
              </el-text>
            </template>
            <template #actions>
              <el-tooltip content="កែប្រែ" placement="top">
                <AppButton size="small" icon="Edit" type="warning" circle />
              </el-tooltip>
            </template>
          </AppTable>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="dialogVisible"
      :title="isEdit ? 'កែប្រែក្រុមហ៊ុន' : 'បន្ថែមក្រុមហ៊ុន'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="ឈ្មោះក្រុមហ៊ុន" prop="name">
          <el-input
            v-model.trim="form.name"
            size="large"
            placeholder="បញ្ចូលឈ្មោះក្រុមហ៊ុន"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="form.base_currency"
              :options="CurrencyOption"
              size="default"
              placeholder="រូបិយប័ណ្ណ"
              clearable
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="form.status"
              :options="StatusOption"
              placeholder="ស្ថានភាព"
              clearable
            />
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <AppButton
          @click="dialogVisible = false"
          size="large"
          :block="false"
          type="warning"
        >
          បោះបង់
        </AppButton>
        <AppButton
          @click="handleSave"
          type="primary"
          :loading="saving"
          size="large"
          :block="false"
        >
          {{ isEdit ? "កែប្រែ" : "បង្កើត" }}
        </AppButton>
      </template>
    </AppDialog>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 6px;
}
</style>
