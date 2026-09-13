<template>
  <div>
    <div id="print-area">
      <el-row :gutter="0" align="middle" class="doc-header">
        <el-col :span="4">
          <div class="header-logo">
            <!-- <el-image :src="url" fit="contain" class="university-logo" /> -->
            <h4 class="label">ក្រុមហ៊ុន {{ company }}</h4>
            <h4 class="label">{{ branch }}</h4>
          </div>
        </el-col>

        <el-col :span="16">
          <h3 class="title-kh">វិក័យប័ត្រ</h3>
          <el-image :src="taktieng" class="taktieng-logo" />
          <h3 class="title-doc">{{ title }}</h3>
        </el-col>

        <el-col :span="4">
          
          <h4 class="tell">លេខទូរសព្ទ: {{ phone }}</h4>
           
        </el-col>
      </el-row>

      <table>
        <thead>
          <tr>
            <th v-if="showIndex">ល.រ</th>
            <th v-for="col in columns" :key="col.key">
              {{ col.label }}
            </th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="(row, index) in rows" :key="row.id ?? index">
            <td v-if="showIndex">{{ index + 1 }}</td>
            <td v-for="col in columns" :key="col.key">
              {{ col.format ? col.format(row[col.key], row) : row[col.key] }}
            </td>
          </tr>
        </tbody>
      </table>
<div class="p-5">
    <el-row :gutter="20">
  <el-col :span="12">
     <h4 class="align-left">អតិថិជន {{ customer }}</h4>
  </el-col>
  <el-col :span="12">
    <h4 class="align-right">សរុប {{ total_amount }}{{ currency }}</h4>
     <h4 class="align-right">កក់មុន {{ paid_amount }}{{ currency }}</h4>
      <h4 class="align-right">នៅខ្វះ {{ outstanding_amount }}{{ currency }}</h4>
  </el-col>
</el-row>
    </div>
</div>
  </div>
</template>

<script setup>
const url = "/logo.png";
const taktieng = "/image.png";
defineProps({
  title: {
    type: String,
    default: "",
  },
  company: {
    type: String,
    default: "",
  },
  branch: {
    type: String,
    default: "",
  },  
  currency: {
    type: String,
    default: "",
  },
  total_amount: {
    type: Number,
    default: "",
  },
  paid_amount: {
    type: Number,
    default: "",
  },
  outstanding_amount: {
    type: Number,
    default: "",
  },
  customer: {
    type: String,
    default: "",
  },
  phone: {
    type: String,
    default:"",
  },
  columns: {
    type: Array,
    required: true,
  },
  rows: {
    type: Array,
    default: () => [],
  },
  showIndex: {
    type: Boolean,
    default: true,
  },
});

function print() {
  window.print();
}

defineExpose({ print });
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Moul&display=swap");

.doc-header {
  margin-bottom: 16px;
}

.align-right {
  margin-top: 10px;
  text-align: right !important;
}

.align-left {
  margin-top: 10px;
  text-align: left !important; 
}

.header-logo {
  padding-top: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.tell {
  padding-top: 100px;
  text-align: right !important;
}

.university-logo {
  width: 70px;
  height: 70px;
}

.taktieng-logo {
  width: 110px;
  height: 25px;
  display: block;
  margin: 4px auto;
}

#print-area {
  padding: 10px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  text-align: center;
  border: 1px solid #000;
  padding: 5px;
}

th {
  text-align: center;
}

.title-kh,
.subtitle-kh,
.title-doc {
  font-family: "Moul", serif;
  font-weight: 200;
  text-align: center;
  margin: 2px 0;
}

.title-kh {
  font-size: 25px;
}

.subtitle-kh {
  font-size: 14px;
}

.label {
  font-family: "Moul", serif;
  font-weight: 200;
  text-align: center;
  font-size: 10px;
  margin: 2px 0;
}

.subtitle-en {
  font-family: Arial, sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-align: center;
  margin: 2px 0;
  letter-spacing: 0.5px;
}

.title-doc {
  font-size: 18px;
  margin-top: 8px;
}
</style>

<style>
/* Not scoped — needs to affect body and elements outside this component */
@media print {
  @page {
    size: A4 portrait;
    margin: 10mm;
  }

  body * {
    background-color: transparent;
    visibility: hidden;
  }

  #print-area,
  #print-area * {
    visibility: visible;
  }

  #print-area {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    padding: 0;
  }

  .print-btn {
    display: none;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    border: 1px solid #000;
    padding: 5px;
  }
}
</style>
