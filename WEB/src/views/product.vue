<script setup>
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getproduct,
  addproduct,
  updateproduct,
  getcompanynopagitaion,
} from "../api/services.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppForm from "../../components/form/AppForm.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";

const notify = useNotification();
const userDataStore = useUserDataStore();

const products = ref([]);
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

// filters
const filters = reactive({
  name: "",
  company_id: null,
});

const companyoptions = ref([]);
async function fetchCompanyOptions() {
  try {
    const res = await getcompanynopagitaion();
    companyoptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }));
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || "Failed to load company");
  }
}

const StatusOption = [
  { label: "ACTIVE", value: "ACTIVE" },
  { label: "INACTIVE", value: "INACTIVE" },
];

// --- Create form: backend accepts an array of products (ProductRequestCreate.product) ---
// so the create dialog lets you queue up several names before saving them all at once.
function blankRow() {
  return { key: Date.now() + Math.random(), name: "", default_price: "" };
}
const createRows = ref([blankRow()]);

function addRow() {
  createRows.value.push(blankRow());
}
function removeRow(key) {
  if (createRows.value.length === 1) return;
  createRows.value = createRows.value.filter((r) => r.key !== key);
}

// --- Edit form: single product, name + status ---
const form = reactive({
  name: "",
  default_price: 0,
  status: "ACTIVE",
});

const rules = {
  name: [{ required: true, message: "Product name is required" }],
};

const canAddProduct = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Product"),
);
const canEditProduct = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "update.Product"),
);

async function fetchProducts() {
  loading.value = true;
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    };
    if (filters.name) params.name = filters.name;
    if (filters.company_id) params.company_id = filters.company_id;
    const res = await getproduct(params);
    products.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  isEdit.value = false;
  editId.value = null;
  createRows.value = [blankRow()];
  dialogVisible.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, {
    name: row.name || "",
    default_price: row.default_price || 0,
    status: row.status || "ACTIVE",
  });
  dialogVisible.value = true;
}

async function handleSave() {
  saving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    if (isEdit.value) {
      const payload = {
        name: form.name,
        default_price: form.default_price,
        status: form.status,
      };
      await updateproduct(editId.value, payload);
      notify.success("កែប្រែផលិតផលបានជោគជ័យ");
    } else {
      // matches ProductRequestCreate: { product: [{ name, default_price }, ...] }
      const rows = createRows.value
        .map((r) => ({
          name: r.name.trim(),
          default_price:
            r.default_price === "" || r.default_price === null || isNaN(Number(r.default_price))
              ? 0
              : Number(r.default_price),
        }))
        .filter((r) => r.name !== "");

      if (rows.length === 0) {
        notify.error("សូមបញ្ចូលឈ្មោះផលិតផលយ៉ាងតិចមួយ");
        saving.value = false;
        useloading.hide();
        return;
      }

      const payload = {
        product: rows,
      };
      await addproduct(payload);
      notify.success("បង្កើតផលិតផលបានជោគជ័យ");
    }
    dialogVisible.value = false;
    fetchProducts();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    saving.value = false;
    useloading.hide();
  }
}

watch(
  () => filters.company_id,
  async () => {
    page.value = 1;
    fetchProducts();
  },
);

watch(
  () => filters.name,
  async () => {
    page.value = 1;
    fetchProducts();
  },
);

onMounted(() => {
  fetchProducts();
  fetchCompanyOptions();
});
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 6 },
        { slot: 'company', span: 6 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #name>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះ"
          clearable
          
        />
      </template>
      <template #company>
        <AppSelect
          v-model="filters.company_id"
          :options="companyoptions"
          label="ក្រុមហ៊ុន"
          placeholder="ក្រុមហ៊ុន"
          clearable
          size="large"
        />
      </template>
      <template #create>
        <AppButton
          v-if="canAddProduct"
          type="primary"
          @click="openCreate"
          :block="false"
          size="large"
        >
          បន្ថែមផលិតផល
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="products"
        :loading="loading"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchProducts"
        :columns="[
          { prop: 'name', label: 'ឈ្មោះឥវ៉ាន់', minWidth: 180 },
          { slot: 'default_price', label: 'តម្លៃលក់រាយ', minWidth: 180 },
          { prop: 'company_name', label: 'ក្រុមហ៊ុន', minWidth: 150 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
        ]"
      >

      <template #default_price="{row}">
        <el-text>{{ row.default_price }} <el-text size="small">{{ row.company_currency }}</el-text></el-text>
      </template>
        <template #status="{ row }">
          <el-text :type="row.status === 'ACTIVE' ? 'success' : 'info'" size="small">
            {{ row.status }}
          </el-text>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
              v-if="canEditProduct"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="openEdit(row)"
            />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="dialogVisible"
      :title="isEdit ? 'កែប្រែឥវ៉ាន់' : 'បន្ថែមឥវ៉ាន់'"
      width="500px"
      :showDefaultFooter="false"
    >
      <!-- EDIT: single product -->
      <AppForm
        v-if="isEdit"
        ref="formRef"
        :model="form"
        :rules="rules"
        :show-actions="true"
        @submit="handleSave"
        submitText="រក្សាទុក"
      >
        <AppInput
          v-model="form.name"
          label="ឈ្មោះឥវ៉ាន់"
          prop="name"
          placeholder="បញ្ចូលឈ្មោះឥវ៉ាន់"
        />
         <AppInput
          v-model.number="form.default_price"
          type="number"
          label="តម្លៃលក់រាយ"
          placeholder="បញ្ចូលតម្លៃលក់រាយ"
        />

        <AppSelect
          v-model="form.status"
          :options="StatusOption"
          label="ស្ថានភាព"
          placeholder="ស្ថានភាព"
          clearable
        />
      </AppForm>

      <!-- CREATE: one or more products at once -->
      <div v-else class="create-rows">
        <div v-for="row in createRows" :key="row.key">
          <el-row :gutter="20">
            <el-col :span="20">
              <el-row :gutter="20">
                <el-col :span="12">
<AppInput
            v-model="row.name"
            
            placeholder="បញ្ចូលឈ្មោះឥវ៉ាន់"
          />
                </el-col>
                <el-col :span="12">
<AppInput
            v-model="row.default_price"
            
            placeholder="បញ្ចូលតម្លៃលក់រាយ"
          />
                </el-col>
              </el-row>

            </el-col>
            <el-col :span="4">
          <AppButton
            v-if="createRows.length > 1"
            size="small"
            icon="Delete"
            type="danger"
            circle
            @click="removeRow(row.key)"
          />
            </el-col>
          </el-row>
        </div>

        <AppButton size="large" type="default" icon="Plus" @click="addRow">
          បន្ថែមមួយទៀត
        </AppButton>

        <div class="create-actions">
          <AppButton type="primary" :loading="saving" @click="handleSave">
            រក្សាទុក
          </AppButton>
        </div>
      </div>
    </AppDialog>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 6px;
}

.create-rows {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.create-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}

.create-row :deep(.app-input) {
  flex: 1;
}

.create-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
</style>