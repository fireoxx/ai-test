<template>
  <div class="app-container">
    <!-- 搜索和操作区域 -->
    <div class="filter-container">
      <el-input
        v-model="listQuery.title"
        placeholder="广告标题"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-select
        v-model="listQuery.position"
        placeholder="广告位置"
        clearable
        style="width: 120px;"
        class="filter-item"
        @change="handleFilter"
      >
        <el-option label="底部" value="bottom" />
        <el-option label="顶部" value="top" />
        <el-option label="弹窗" value="popup" />
      </el-select>
      <el-select
        v-model="listQuery.status"
        placeholder="状态"
        clearable
        style="width: 100px;"
        class="filter-item"
        @change="handleFilter"
      >
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="2" />
      </el-select>
      <el-select
        v-model="listQuery.deviceType"
        placeholder="设备类型"
        clearable
        style="width: 120px;"
        class="filter-item"
        @change="handleFilter"
      >
        <el-option label="全部" value="all" />
        <el-option label="手机" value="mobile" />
        <el-option label="电脑" value="pc" />
      </el-select>
      <el-button
        class="filter-item"
        type="primary"
        icon="Search"
        @click="handleFilter"
      >
        搜索
      </el-button>
      <el-button
        class="filter-item"
        type="success"
        icon="Plus"
        @click="handleCreate"
      >
        新增
      </el-button>
    </div>

    <!-- 数据表格 -->
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" prop="ID" align="center" width="80" />
      <el-table-column label="广告标题" prop="title" min-width="150" />
      <el-table-column label="位置" align="center" width="100">
        <template #default="{ row }">
          <el-tag :type="getPositionType(row.position)">
            {{ getPositionText(row.position) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" align="center" width="80">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="2"
            @change="handleStatusChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="强制弹窗" align="center" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.forcePopup"
            @change="handleForcePopupChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="设备类型" align="center" width="100">
        <template #default="{ row }">
          {{ getDeviceTypeText(row.deviceType) }}
        </template>
      </el-table-column>
      <el-table-column label="排序" prop="sort" align="center" width="80" />
      <el-table-column label="点击/展示" align="center" width="120">
        <template #default="{ row }">
          <div>{{ row.clickCount }}/{{ row.viewCount }}</div>
        </template>
      </el-table-column>
      <el-table-column label="有效期" min-width="180">
        <template #default="{ row }">
          <div v-if="row.startTime || row.endTime">
            {{ row.startTime || '不限' }} ~ {{ row.endTime || '不限' }}
          </div>
          <div v-else>长期有效</div>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" width="180">
        <template #default="{ row }">
          <span>{{ formatDate(row.CreatedAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="200" class-name="small-padding fixed-width">
        <template #default="{ row }">
          <el-button type="primary" size="small" icon="Edit" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button type="danger" size="small" icon="Delete" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="gva-pagination">
      <el-pagination
        v-show="total > 0"
        :current-page="listQuery.page"
        :page-size="listQuery.pageSize"
        :page-sizes="[10, 20, 30, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogFormVisible"
      width="600px"
    >
      <el-form
        ref="dataFormRef"
        :model="temp"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="广告标题" prop="title">
          <el-input v-model="temp.title" placeholder="请输入广告标题" />
        </el-form-item>
        <el-form-item label="广告描述" prop="description">
          <el-input
            v-model="temp.description"
            type="textarea"
            :rows="3"
            placeholder="请输入广告描述"
          />
        </el-form-item>
        <el-form-item label="图片URL" prop="imageUrl">
          <el-input v-model="temp.imageUrl" placeholder="请输入图片URL（可选）" />
        </el-form-item>
        <el-form-item label="跳转链接" prop="link">
          <el-input v-model="temp.link" placeholder="请输入跳转链接" />
        </el-form-item>
        <el-form-item label="广告位置" prop="position">
          <el-select v-model="temp.position" placeholder="请选择广告位置">
            <el-option label="底部" value="bottom" />
            <el-option label="顶部" value="top" />
            <el-option label="弹窗" value="popup" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备类型" prop="deviceType">
          <el-select v-model="temp.deviceType" placeholder="请选择设备类型">
            <el-option label="全部" value="all" />
            <el-option label="手机" value="mobile" />
            <el-option label="电脑" value="pc" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="temp.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="temp.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="强制弹窗">
          <el-switch v-model="temp.forcePopup" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">
            开启后，用户进入页面时会强制弹出此广告
          </span>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="temp.startTime"
            type="datetime"
            placeholder="选择开始时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="temp.endTime"
            type="datetime"
            placeholder="选择结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogFormVisible = false">取消</el-button>
          <el-button type="primary" @click="dialogStatus === 'create' ? createData() : updateData()">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import {
  getAdvertisementList,
  createAdvertisement,
  updateAdvertisement,
  deleteAdvertisement
} from '@/api/advertisement/advertisement'

// 数据
const list = ref([])
const total = ref(0)
const listLoading = ref(true)

// 查询参数
const listQuery = reactive({
  page: 1,
  pageSize: 20,
  title: '',
  position: '',
  status: null,
  deviceType: ''
})

// 对话框
const dialogFormVisible = ref(false)
const dialogStatus = ref('')
const dialogTitle = ref('')
const dataFormRef = ref()

// 临时数据
const temp = ref({
  id: undefined,
  title: '',
  description: '',
  imageUrl: '',
  link: '',
  position: 'bottom',
  status: 1,
  sort: 0,
  startTime: '',
  endTime: '',
  deviceType: 'all',
  forcePopup: false
})

// 表单验证规则
const rules = reactive({
  title: [{ required: true, message: '请输入广告标题', trigger: 'blur' }],
  description: [{ required: true, message: '请输入广告描述', trigger: 'blur' }],
  link: [
    { required: true, message: '请输入跳转链接', trigger: 'blur' },
    { type: 'url', message: '请输入正确的URL地址', trigger: 'blur' }
  ],
  position: [{ required: true, message: '请选择广告位置', trigger: 'change' }],
  deviceType: [{ required: true, message: '请选择设备类型', trigger: 'change' }]
})

// 获取列表
const getList = () => {
  listLoading.value = true
  getAdvertisementList(listQuery)
    .then(response => {
      list.value = response.data.list
      total.value = response.data.total
    })
    .finally(() => {
      listLoading.value = false
    })
}

// 搜索
const handleFilter = () => {
  listQuery.page = 1
  getList()
}

// 分页处理
const handleCurrentChange = (val) => {
  listQuery.page = val
  getList()
}

const handleSizeChange = (val) => {
  listQuery.pageSize = val
  getList()
}

// 重置表单
const resetTemp = () => {
  temp.value = {
    id: undefined,
    title: '',
    description: '',
    imageUrl: '',
    link: '',
    position: 'bottom',
    status: 1,
    sort: 0,
    startTime: '',
    endTime: '',
    deviceType: 'all',
    forcePopup: false
  }
}

// 新增
const handleCreate = () => {
  resetTemp()
  dialogStatus.value = 'create'
  dialogTitle.value = '新增广告'
  dialogFormVisible.value = true
  nextTick(() => {
    dataFormRef.value?.clearValidate()
  })
}

// 编辑
const handleUpdate = (row) => {
  temp.value = Object.assign({}, row)
  dialogStatus.value = 'update'
  dialogTitle.value = '编辑广告'
  dialogFormVisible.value = true
  nextTick(() => {
    dataFormRef.value?.clearValidate()
  })
}

// 创建数据
const createData = () => {
  dataFormRef.value.validate((valid) => {
    if (valid) {
      createAdvertisement(temp.value).then(() => {
        dialogFormVisible.value = false
        ElMessage.success('创建成功')
        getList()
      })
    }
  })
}

// 更新数据
const updateData = () => {
  dataFormRef.value.validate((valid) => {
    if (valid) {
      updateAdvertisement(temp.value).then(() => {
        dialogFormVisible.value = false
        ElMessage.success('更新成功')
        getList()
      })
    }
  })
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该广告吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteAdvertisement({ id: row.ID }).then(() => {
      ElMessage.success('删除成功')
      getList()
    })
  })
}

// 状态切换
const handleStatusChange = (row) => {
  updateAdvertisement({
    id: row.ID,
    status: row.status
  }).then(() => {
    ElMessage.success('状态更新成功')
  })
}

// 强制弹窗切换
const handleForcePopupChange = (row) => {
  updateAdvertisement({
    id: row.ID,
    forcePopup: row.forcePopup
  }).then(() => {
    ElMessage.success('强制弹窗设置已更新')
  })
}

// 辅助函数
const getPositionText = (position) => {
  const map = {
    bottom: '底部',
    top: '顶部',
    popup: '弹窗'
  }
  return map[position] || position
}

const getPositionType = (position) => {
  const map = {
    bottom: 'success',
    top: 'warning',
    popup: 'danger'
  }
  return map[position] || ''
}

const getDeviceTypeText = (deviceType) => {
  const map = {
    all: '全部',
    mobile: '手机',
    pc: '电脑'
  }
  return map[deviceType] || deviceType
}

// 生命周期
onMounted(() => {
  getList()
})
</script>

<style scoped>
.filter-container {
  margin-bottom: 20px;
}

.filter-item {
  margin-right: 10px;
}

.dialog-footer {
  text-align: right;
}
</style>