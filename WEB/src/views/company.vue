<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getcompany,
  addcompany,
  updatecompany,
  addbranch,
  updatebranch,
  createUser,
  getrole,
  updateUser,
} from "../api/services.js";
import { useAuthStore } from "../stores/auth";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppForm from "../../components/form/AppForm.vue";
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

const showDetail = ref(false);
const selectedInvoice = ref(null);

function openDetail(row) {
  selectedInvoice.value = row;
  showDetail.value = true;
}

const showSubDetail = ref(false);
const selectedSubInvoice = ref(null);
const selectedSubInvoiceCompany = ref(null);
function openSubDetail(row, companyRow) {
  selectedSubInvoice.value = row;
  selectedSubInvoiceCompany.value = companyRow; // NEW
  showSubDetail.value = true;
}

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

// branch dialog state

const BranchStatusOption = [
  { label: "ACTIVE", value: "ACTIVE" },
  { label: "INACTIVE", value: "INACTIVE" },
];

const branchDialogVisible = ref(false);
const isEditBranch = ref(false);
const editBranchId = ref(null);
const branchSaving = ref(false);
const branchFormRef = ref();

const branchform = reactive({
  company_id: null,
  name: "",
  address: "",
  phone: "",
  status: "ACTIVE",
});

const branchRules = {
  name: [{ required: true, message: "Branch name is required" }],
  address: [{ required: false }],
  status: [{ required: true, message: "Status is required" }],
};

// open dialog to ADD a branch under a specific company
function openCreateBranch(companyRow) {
  isEditBranch.value = false;
  editBranchId.value = null;
  branchform.company_id = companyRow.id;
  branchform.name = "";
  branchform.address = "";
  branchform.phone = "";
  branchform.status = "ACTIVE";
  branchDialogVisible.value = true;
}

// open dialog to EDIT an existing branch
function openEditBranch(companyRow, branchRow) {
  isEditBranch.value = true;
  editBranchId.value = branchRow.id;
  branchform.company_id = companyRow.id;
  branchform.name = branchRow.name || "";
  branchform.address = branchRow.address || "";
  branchform.phone = branchRow.phone || "";
  branchform.status = branchRow.status || "ACTIVE";
  branchDialogVisible.value = true;
}

async function handleSaveBranch() {
  await branchFormRef.value.validate();
  branchSaving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    if (isEditBranch.value) {
      const payload = {
        company_id: branchform.company_id,
        name: branchform.name,
        address: branchform.address,
        phone: branchform.phone,
        status: branchform.status,
      };
      await updatebranch(editBranchId.value, payload);
      notify.success("កែប្រែសាខាបានជោគជ័យ");
    } else {
      // create request has no `status` field per BranchRequestCreate
      const payload = {
        company_id: branchform.company_id,
        name: branchform.name,
        address: branchform.address,
        phone: branchform.phone,
      };
      await addbranch(payload);
      notify.success("បង្កើតសាខាបានជោគជ័យ");
    }
    branchDialogVisible.value = false;
    fetchCompanies(); // refresh so nested branches update
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    branchSaving.value = false;
    useloading.hide();
  }
}
// branch

const branchcolumns = [
  { prop: "name", slot: "name", label: "ឈ្មោះ", minwidth: 120 },
  { prop: "address", label: "ទីតាំងសាខា", minwidth: 120 },
  { prop: "phone", label: "លេខទូរសព្ទសាខា", minwidth: 120 },
  { slot: "status", label: "ស្ថានភាព", width: 120 },
];

const usercolumns = [
  { prop: "name", label: "ឈ្មោះ", minwidth: 120 },
  { prop: "email", label: "អុីម៉ែល", minwidth: 120 },
  { slot: "role_name", label: "តួនាទី", minwidth: 120 },
  { slot: "manage_branch", label: "សិទ្ធមើលសាខា", minwidth: 120 },
  { slot: "status", label: "ស្ថានភាព", minwidth: 120 },
];

function getBranchNames(companyRow, branchIds) {
  if (!Array.isArray(branchIds) || !companyRow?.branches) return [];
  return branchIds
    .map((id) => companyRow.branches.find((b) => b.id === id)?.name)
    .filter(Boolean);
}

const canAddCompany = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.company"),
);

const canEditCompany = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "update.company"),
);

const canAddBranch = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Branch"),
);

const canEditBranch = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "update.Branch"),
);

const canAddUser = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.user"),
);
const canEditUser = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.user"),
);

const ManageBranchOption = [
  { label: "ONE", value: "ONE" },
  { label: "MULTIPLE", value: "MULTIPLE" },
  { label: "ALL", value: "ALL" },
];

const UserStatusOption = [
  { label: "ACTIVE", value: "ACTIVE" },
  { label: "DISABLED", value: "DISABLED" },
];

const roles = ref([]);
const roleOptions = computed(() =>
  roles.value.map((r) => ({ label: r.display_name || r.name, value: r.id })),
);

async function fetchRoles() {
  try {
    const res = await getrole();
    roles.value = res.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  }
}

// user dialog state

const userDialogVisible = ref(false);
const isEditUser = ref(false);
const editUserId = ref(null);
const userSaving = ref(false);
const userFormRef = ref();
const currentCompanyBranches = ref([]);

const branchOptions = computed(() =>
  currentCompanyBranches.value.map((b) => ({ label: b.name, value: b.id })),
);

const userform = reactive({
  branch_id: null,
  name: "",
  role_id: null,
  manage_branch: "ONE",
  branch_ids: [],
  status: "",
});

const userRules = {
  name: [{ required: true, message: "Name is required" }],
  role_id: [{ required: true, message: "Role is required" }],
  manage_branch: [{ required: true, message: "Manage branch is required" }],
};

function openCreateUser(companyRow, branchRow) {
  isEditUser.value = false;
  editUserId.value = null;
  currentCompanyBranches.value = companyRow.branches || [];
  userform.branch_id = branchRow.id;
  userform.name = "";
  userform.role_id = null;
  userform.manage_branch = "ONE";
  userform.branch_ids = [];
  userDialogVisible.value = true;
}

function openEditUser(companyRow, branchRow, userRow) {
  isEditUser.value = true;
  editUserId.value = userRow.id;
  currentCompanyBranches.value = companyRow.branches || [];
  userform.branch_id = branchRow.id;
  userform.name = userRow.name || "";
  userform.role_id = userRow.role_id || null;
  userform.manage_branch = userRow.manage_branch || "ONE";
  userform.status = userRow.status;
  userform.branch_ids = Array.isArray(userRow.branch_ids)
    ? userRow.branch_ids.map((b) => b.branch_id)
    : [];
  userDialogVisible.value = true;
}

async function handleSaveUser() {
  await userFormRef.value.validate();
  userSaving.value = true;
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    const payload = {
      branch_id: userform.branch_id,
      name: userform.name,
      role_id: userform.role_id,
      manage_branch: userform.manage_branch,
      status: userform.status,
    };
    // only send branch_ids when it's actually needed — matches
    // input.ManageBranch == MULTIPLE check on the Go side
    if (userform.manage_branch === "MULTIPLE") {
      payload.branch_ids = userform.branch_ids;
    }

    if (isEditUser.value) {
      await updateUser(editUserId.value, payload);
      notify.success("កែប្រែបុគ្គលិកបានជោគជ័យ");
    } else {
      await createUser(payload);
      notify.success("បង្កើតបុគ្គលិកបានជោគជ័យ");
    }
    userDialogVisible.value = false;
    fetchCompanies();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    userSaving.value = false;
    useloading.hide();
  }
}

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
    console.log(companies.value);
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

onMounted(() => {
  fetchCompanies();
  fetchRoles();
});
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
          <el-text
            :type="row.status === 'ACTIVE' ? 'success' : 'danger'"
            size="small"
          >
            {{ row.status }}
          </el-text>
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
          <el-tooltip content="ថែមសាខា" placement="top">
            <AppButton
              v-if="canAddBranch"
              circle
              size="small"
              type="primary"
              icon="Plus"
              @click="openCreateBranch(row)"
            />
          </el-tooltip>
          <el-tooltip content="មើលលំអិត" placement="top">
            <AppButton
              size="small"
              icon="View"
              type="success"
              circle
              @click="openDetail(row)"
            />
          </el-tooltip>
        </template>

        <template #expand="{ row: companyRow }">
          <el-divider content-position="left">
            <el-text>
              សាខាសរុប {{ companyRow.branches.length }}
              <el-tooltip content="ថែមសាខា" placement="top">
                <AppButton
                  v-if="canAddBranch"
                  circle
                  size="small"
                  type="primary"
                  icon="Plus"
                  @click="openCreateBranch(companyRow)"
                />
              </el-tooltip>
            </el-text>
          </el-divider>

          <AppTable
            expandable
            :data="companyRow.branches"
            :columns="branchcolumns"
            :show-pagination="false"
          >
            <template #status="{ row: branchRow }">
              <el-text
                :type="branchRow.status === 'ACTIVE' ? 'success' : 'danger'"
                size="small"
              >
                {{ branchRow.status }}
              </el-text>
            </template>

            <template #name="{ row: branchRow }">
              <el-text>
                {{ branchRow.name }} |
                <el-text size="small" type="primary">{{
                  branchRow.code
                }}</el-text>
              </el-text>
            </template>

            <template #actions="{ row: branchRow }">
              <el-tooltip content="កែប្រែ" placement="top">
                <AppButton
                  v-if="canEditBranch"
                  size="small"
                  icon="Edit"
                  type="warning"
                  circle
                  @click="openEditBranch(companyRow, branchRow)"
                />
              </el-tooltip>
              <el-tooltip content="ថែមបុគ្គលិក" placement="top">
                <AppButton
                  v-if="canAddUser"
                  circle
                  size="small"
                  type="primary"
                  icon="Plus"
                  @click="openCreateUser(companyRow, branchRow)"
                />
              </el-tooltip>
            </template>

            <template #expand="{ row: branchRow }">
              <el-divider content-position="left">
                <el-text>
                  បុគ្គលិកសរុប {{ branchRow.users.length }}
                  <el-tooltip content="ថែមបុគ្គលិក" placement="top">
                    <AppButton
                      v-if="canAddUser"
                      circle
                      size="small"
                      type="primary"
                      icon="Plus"
                      @click="openCreateUser(companyRow, branchRow)"
                    />
                  </el-tooltip>
                </el-text>
              </el-divider>

              <AppTable
                expandable
                :data="branchRow.users"
                :columns="usercolumns"
                :show-pagination="false"
              >
                <template #manage_branch="{ row: userRow }">
                  <el-text>
                    {{ userRow.manage_branch }}
                    <template v-if="userRow.manage_branch === 'MULTIPLE'">
                      |
                      <el-text size="small" type="primary">
                        {{
                          getBranchNames(
                            companyRow,
                            Array.isArray(userRow.branch_ids)
                              ? userRow.branch_ids.map((b) => b.branch_id)
                              : [],
                          ).join(", ")
                        }}
                      </el-text>
                    </template>
                  </el-text>
                </template>

                <template #role_name="{ row: userRow }">
                  <el-text>
                    {{ userRow.role_display_name }} |
                    <el-text size="small" type="primary">{{
                      userRow.role_name
                    }}</el-text>
                  </el-text>
                </template>

                <template #status="{ row: userRow }">
                  <el-text
                    :type="userRow.status === 'ACTIVE' ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ userRow.status }}
                  </el-text>
                </template>

                <template #actions="{ row: userRow }">
                  <el-tooltip content="កែប្រែ" placement="top">
                    <AppButton
                      v-if="canEditUser"
                      size="small"
                      icon="Edit"
                      type="warning"
                      circle
                      @click="openEditUser(companyRow, branchRow, userRow)"
                    />
                  </el-tooltip>
                </template>
              </AppTable>
            </template>
          </AppTable>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="showDetail"
      title="សាខា"
      width="60%"
      :showDefaultFooter="false"
    >
      <AppTable
        expandable
        :data="selectedInvoice?.branches || []"
        :columns="branchcolumns"
        :show-pagination="false"
      >
        <template #status="{ row: branchRow }">
          <el-text
            :type="branchRow.status === 'ACTIVE' ? 'success' : 'danger'"
            size="small"
          >
            {{ branchRow.status }}
          </el-text>
        </template>

        <template #name="{ row: branchRow }">
          <el-text>
            {{ branchRow.name }} |
            <el-text size="small" type="primary">{{ branchRow.code }}</el-text>
          </el-text>
        </template>

        <template #actions="{ row: branchRow }">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
              v-if="canEditBranch"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="openEditBranch(selectedInvoice, branchRow)"
            />
          </el-tooltip>
          <el-tooltip content="ថែមបុគ្គលិក" placement="top">
            <AppButton
              v-if="canAddUser"
              circle
              size="small"
              type="primary"
              icon="Plus"
              @click="openCreateUser(selectedInvoice, branchRow)"
            />
          </el-tooltip>
          <el-tooltip content="មើលលំអិត" placement="top">
            <AppButton
              size="small"
              icon="View"
              type="success"
              circle
              @click="openSubDetail(branchRow, selectedInvoice)"
            />
          </el-tooltip>
        </template>
      </AppTable>
    </AppDialog>

    <AppDialog
      v-model="showSubDetail"
      title="បុគ្គលិក"
      width="60%"
      :showDefaultFooter="false"
    >
      <AppTable
        expandable
        :data="selectedSubInvoice?.users || []"
        :columns="usercolumns"
        :show-pagination="false"
      >
        <template #manage_branch="{ row: userRow }">
          <el-text>
            {{ userRow.manage_branch }}
            <template v-if="userRow.manage_branch === 'MULTIPLE'">
              |
              <el-text size="small" type="primary">
                {{
                  getBranchNames(
                    companyRow,
                    Array.isArray(userRow.branch_ids)
                      ? userRow.branch_ids.map((b) => b.branch_id)
                      : [],
                  ).join(", ")
                }}
              </el-text>
            </template>
          </el-text>
        </template>

        <template #role_name="{ row: userRow }">
          <el-text>
            {{ userRow.role_display_name }} |
            <el-text size="small" type="primary">{{
              userRow.role_name
            }}</el-text>
          </el-text>
        </template>

        <template #status="{ row: userRow }">
          <el-text
            :type="userRow.status === 'ACTIVE' ? 'success' : 'danger'"
            size="small"
          >
            {{ userRow.status }}
          </el-text>
        </template>

        <template #actions="{ row: userRow }">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
              v-if="canEditUser"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="
                openEditUser(
                  selectedSubInvoiceCompany,
                  selectedSubInvoice,
                  userRow,
                )
              "
            />
          </el-tooltip>
        </template>
      </AppTable>
    </AppDialog>

    <AppDialog
      v-model="dialogVisible"
      :title="isEdit ? 'កែប្រែក្រុមហ៊ុន' : 'បន្ថែមក្រុមហ៊ុន'"
      width="500px"
      :showDefaultFooter="false"
    >
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :show-actions="true"
        @submit="handleSave"
        submitText="រក្សាទុក"
      >
        <AppInput
          v-model="form.name"
          label="ឈ្មោះក្រុមហ៊ុន"
          prop="name"
          placeholder="បញ្ចូលឈ្មោះក្រុមហ៊ុន"
        >
        </AppInput>

        <AppSelect
          v-model="form.base_currency"
          :options="CurrencyOption"
          size="default"
          placeholder="រូបិយប័ណ្ណ"
          clearable
        />

        <AppSelect
          v-model="form.status"
          :options="StatusOption"
          placeholder="ស្ថានភាព"
          clearable
        />
      </AppForm>
    </AppDialog>

    <AppDialog
      v-model="branchDialogVisible"
      :title="isEditBranch ? 'កែប្រែសាខា' : 'បន្ថែមសាខា'"
      width="500px"
      :showDefaultFooter="false"
    >
      <AppForm
        ref="branchFormRef"
        :model="branchform"
        :rules="branchRules"
        :show-actions="true"
        @submit="handleSaveBranch"
        submitText="រក្សាទុក"
      >
        <AppInput
          label="ឈ្មោះសាខា"
          prop="name"
          v-model="branchform.name"
          placeholder="បញ្ចូលឈ្មោះសាខា"
        >
        </AppInput>
        <AppInput
          label="ទីតាំងសាខា"
          prop="address"
          v-model="branchform.address"
          placeholder="បញ្ចូលទីតាំងសាខា"
        >
        </AppInput>
        <AppInput
          label="លេខទូរសព្ទសាខា"
          v-model="branchform.phone"
          placeholder="បញ្ចូលលេខទូរសព្ទសាខា"
        >
        </AppInput>
        <AppSelect
          v-if="isEditBranch"
          v-model="branchform.status"
          :options="BranchStatusOption"
          placeholder="ស្ថានភាព"
          label="ស្ថានភាព"
          prop="status"
          clearable
        />
      </AppForm>
    </AppDialog>

    <AppDialog
      v-model="userDialogVisible"
      :title="isEditUser ? 'កែប្រែបុគ្គលិក' : 'បន្ថែមបុគ្គលិក'"
      width="500px"
      :showDefaultFooter="false"
    >
      <AppForm
        ref="userFormRef"
        :model="userform"
        :rules="userRules"
        :show-actions="true"
        @submit="handleSaveUser"
        submitText="រក្សាទុក"
      >
        <AppInput
          label="ឈ្មោះ"
          prop="name"
          v-model="userform.name"
          placeholder="បញ្ចូលឈ្មោះបុគ្គលិក"
        />

        <AppSelect
          v-model="userform.role_id"
          :options="roleOptions"
          label="តួនាទី"
          prop="role_id"
          placeholder="ជ្រើសរើសតួនាទី"
          clearable
        />

        <AppSelect
          v-model="userform.manage_branch"
          :options="ManageBranchOption"
          label="សិទ្ធិមើលសាខា"
          prop="manage_branch"
        />

        <AppSelect
          v-if="userform.manage_branch === 'MULTIPLE'"
          v-model="userform.branch_ids"
          :options="branchOptions"
          label="សាខា"
          prop="branch_ids"
          placeholder="ជ្រើសរើសសាខា"
          multiple
          clearable
        />
        <AppSelect
          v-model="userform.status"
          :options="UserStatusOption"
          label="ស្ថានភាព"
        />
      </AppForm>
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
