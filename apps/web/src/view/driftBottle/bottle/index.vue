<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索内容/昵称"
          style="width: 200px; margin-right: 10px"
          clearable
          @keyup.enter="getTableData"
        />
        <el-select
          v-model="searchStatus"
          placeholder="状态筛选"
          style="width: 120px; margin-right: 10px"
          clearable
        >
          <el-option label="漂流中" :value="1" />
          <el-option label="已被捞起" :value="2" />
          <el-option label="已回复" :value="3" />
        </el-select>
        <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
      </div>
      <el-table :data="tableData" style="width: 100%" row-key="ID">
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="昵称" prop="nickname" width="120" />
        <el-table-column align="left" label="内容" prop="content" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.status)">
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="设备ID" prop="deviceId" width="160" show-overflow-tooltip />
        <el-table-column align="left" label="投瓶时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="160">
          <template #default="scope">
            <el-button type="primary" link icon="chat-dot-round" @click="viewReplies(scope.row)">查看回复</el-button>
            <el-button type="danger" link icon="delete" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 查看回复弹窗 -->
    <el-dialog v-model="replyDialogVisible" title="回复列表" width="700px">
      <div v-if="currentBottle" class="bottle-content-box">
        <div class="bottle-content-label">瓶子内容：</div>
        <div class="bottle-content-text">{{ currentBottle.content }}</div>
      </div>
      <el-table :data="replyList" style="width: 100%; margin-top: 16px" row-key="ID">
        <el-table-column align="left" label="回复昵称" prop="nickname" width="120" />
        <el-table-column align="left" label="回复内容" prop="content" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="回复时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="100">
          <template #default="scope">
            <el-button type="danger" link icon="delete" @click="handleDeleteReply(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="replyList.length === 0" class="no-reply-text">暂无回复</div>
      <template #footer>
        <el-button @click="replyDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { adminGetBottleList, adminGetReplyList, adminDeleteBottle, adminDeleteReply } from '@/api/driftBottle'

defineOptions({ name: 'DriftBottleList' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchKeyword = ref('')
const searchStatus = ref(undefined)

const statusLabel = (status) => {
  const map = { 1: '漂流中', 2: '已被捞起', 3: '已回复' }
  return map[status] || '未知'
}

const statusTagType = (status) => {
  const map = { 1: 'primary', 2: 'warning', 3: 'success' }
  return map[status] || 'info'
}

const getTableData = async () => {
  const res = await adminGetBottleList({
    page: page.value,
    pageSize: pageSize.value,
    keyword: searchKeyword.value,
    status: searchStatus.value || 0
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除这个漂流瓶吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await adminDeleteBottle({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 回复相关
const replyDialogVisible = ref(false)
const replyList = ref([])
const currentBottle = ref(null)

const viewReplies = async (row) => {
  currentBottle.value = row
  replyDialogVisible.value = true
  const res = await adminGetReplyList({ bottleId: row.ID, page: 1, pageSize: 100 })
  if (res.code === 0) {
    replyList.value = res.data.list
  }
}

const handleDeleteReply = (row) => {
  ElMessageBox.confirm('确定要删除这条回复吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await adminDeleteReply({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      replyList.value = replyList.value.filter(r => r.ID !== row.ID)
    }
  })
}

onMounted(() => {
  getTableData()
})
</script>

<style scoped>
.bottle-content-box {
  background: #f5f7fa;
  border-radius: 6px;
  padding: 12px 16px;
}
.bottle-content-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
}
.bottle-content-text {
  font-size: 14px;
  color: #303133;
  line-height: 1.6;
}
.no-reply-text {
  text-align: center;
  color: #909399;
  padding: 20px 0;
}
</style>
