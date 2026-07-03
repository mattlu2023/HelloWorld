<template>
  <div class="campaigns-container">
    <div class="page-header">
      <div class="header-info">
        <h1 class="page-title">广告活动管理</h1>
        <p class="page-subtitle">管理和监控所有广告投放活动</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd" class="btn-primary">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          <span>创建活动</span>
        </el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索活动名称..."
        class="search-input"
        prefix-icon="Search"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="statusFilter" placeholder="状态筛选" class="filter-select">
        <el-option label="全部" value="" />
        <el-option label="进行中" value="active" />
        <el-option label="已结束" value="ended" />
        <el-option label="暂停中" value="paused" />
      </el-select>
    </div>

    <div class="stats-row">
      <div class="mini-stat">
        <span class="mini-stat-value">{{ campaigns.length }}</span>
        <span class="mini-stat-label">活动总数</span>
      </div>
      <div class="mini-stat">
        <span class="mini-stat-value">{{ activeCount }}</span>
        <span class="mini-stat-label">进行中</span>
      </div>
      <div class="mini-stat">
        <span class="mini-stat-value">{{ totalBudget }}</span>
        <span class="mini-stat-label">总预算(元)</span>
      </div>
    </div>

    <div class="table-card">
      <el-table :data="filteredCampaigns" style="width: 100%" class="custom-table">
        <el-table-column prop="name" label="活动名称" min-width="200">
          <template #default="{ row }">
            <div class="campaign-name-cell">
              <div class="campaign-dot" :class="row.status"></div>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small" class="status-tag">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="budget" label="预算" width="120">
          <template #default="{ row }">
            <span class="budget-value">¥{{ row.budget.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="130" />
        <el-table-column prop="end_date" label="结束日期" width="130" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" class="btn-action btn-edit" @click="handleEdit(row)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path>
              </svg>
              编辑
            </el-button>
            <el-button size="small" class="btn-action btn-delete" @click="handleDelete(row)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="empty-state" v-if="filteredCampaigns.length === 0">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
          <line x1="9" y1="9" x2="15" y2="9"></line>
          <line x1="12" y1="12" x2="12" y2="15"></line>
        </svg>
        <p>暂无广告活动</p>
        <el-button type="primary" @click="handleAdd">创建第一个活动</el-button>
      </div>
    </div>

    <el-dialog
      :title="isEdit ? '编辑活动' : '创建活动'"
      v-model="dialogVisible"
      width="500px"
      class="custom-dialog"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="活动描述" prop="description">
          <el-input v-model="form.description" type="textarea" placeholder="请输入活动描述" :rows="3" />
        </el-form-item>
        <el-form-item label="预算" prop="budget">
          <el-input-number v-model="form.budget" :min="0" :precision="0" :step="1000" placeholder="请输入预算" style="width: 100%" />
        </el-form-item>
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker v-model="form.start_date" type="date" placeholder="选择开始日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束日期" prop="end_date">
          <el-date-picker v-model="form.end_date" type="date" placeholder="选择结束日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="选择状态">
            <el-option label="进行中" value="active" />
            <el-option label="暂停中" value="paused" />
            <el-option label="已结束" value="ended" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          {{ isEdit ? '保存修改' : '创建活动' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCampaigns, createCampaign, updateCampaign, deleteCampaign } from '@/api'

const searchQuery = ref('')
const statusFilter = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref(null)

const campaigns = ref([])

const form = reactive({
  id: null,
  name: '',
  description: '',
  budget: 0,
  start_date: '',
  end_date: '',
  status: 'active'
})

const rules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  budget: [{ required: true, message: '请输入预算', trigger: 'blur' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
}

const filteredCampaigns = computed(() => {
  return campaigns.value.filter(item => {
    const matchSearch = !searchQuery.value || item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchStatus = !statusFilter.value || item.status === statusFilter.value
    return matchSearch && matchStatus
  })
})

const activeCount = computed(() => {
  return campaigns.value.filter(item => item.status === 'active').length
})

const totalBudget = computed(() => {
  return campaigns.value.reduce((sum, item) => sum + item.budget, 0).toLocaleString()
})

const getStatusType = (status) => {
  const types = {
    active: 'success',
    ended: 'info',
    paused: 'warning'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    active: '进行中',
    ended: '已结束',
    paused: '暂停中'
  }
  return texts[status] || status
}

const loadCampaigns = async () => {
  try {
    const res = await getCampaigns()
    campaigns.value = res.data
  } catch (error) {
    ElMessage.error('加载活动列表失败')
  }
}

const handleSearch = () => {}

const handleAdd = () => {
  isEdit.value = false
  form.id = null
  form.name = ''
  form.description = ''
  form.budget = 0
  form.start_date = ''
  form.end_date = ''
  form.status = 'active'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.description = row.description || ''
  form.budget = row.budget
  form.start_date = row.start_date
  form.end_date = row.end_date
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除活动「${row.name}」吗？`, '确认删除', {
      type: 'warning'
    })
    await deleteCampaign(row.id)
    ElMessage.success('删除成功')
    loadCampaigns()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitLoading.value = true
    
    if (isEdit.value) {
      await updateCampaign(form.id, {
        name: form.name,
        description: form.description,
        budget: form.budget,
        start_date: form.start_date,
        end_date: form.end_date,
        status: form.status
      })
      ElMessage.success('更新成功')
    } else {
      await createCampaign({
        name: form.name,
        description: form.description,
        budget: form.budget,
        start_date: form.start_date,
        end_date: form.end_date,
        status: form.status
      })
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadCampaigns()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  } finally {
    submitLoading.value = false
  }
}

loadCampaigns()
</script>

<style scoped>
.campaigns-container {
  padding: var(--spacing-8);
  min-height: calc(100vh - var(--header-height));
  background: var(--color-bg);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-6);
}

.header-info {
  display: flex;
  flex-direction: column;
}

.page-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  margin: 0 0 var(--spacing-2);
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-normal);
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-lg);
}

.filter-bar {
  display: flex;
  gap: var(--spacing-4);
  margin-bottom: var(--spacing-6);
}

.search-input {
  width: 300px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: var(--radius-lg);
  border-color: var(--color-border);
}

.filter-select {
  width: 160px;
}

.filter-select :deep(.el-select__wrapper) {
  border-radius: var(--radius-lg);
  border-color: var(--color-border);
}

.stats-row {
  display: flex;
  gap: var(--spacing-6);
  margin-bottom: var(--spacing-6);
}

.mini-stat {
  display: flex;
  flex-direction: column;
  padding: var(--spacing-4) var(--spacing-6);
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  min-width: 140px;
}

.mini-stat-value {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  letter-spacing: -0.02em;
}

.mini-stat-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-top: var(--spacing-1);
}

.table-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.custom-table :deep(.el-table) {
  border: none;
}

.custom-table :deep(.el-table__header-wrapper) {
  background: var(--color-bg);
}

.custom-table :deep(.el-table__header th) {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  font-weight: var(--font-weight-medium);
  font-size: var(--font-size-sm);
  padding: var(--spacing-4) var(--spacing-6);
  border-bottom: 1px solid var(--color-border-light);
}

.custom-table :deep(.el-table__body tr:hover) {
  background: rgba(6, 182, 212, 0.03);
}

.custom-table :deep(.el-table__body td) {
  padding: var(--spacing-4) var(--spacing-6);
  border-bottom: 1px solid var(--color-border-light);
  color: var(--color-text-primary);
}

.campaign-name-cell {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}

.campaign-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.campaign-dot.active {
  background: var(--color-success);
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.campaign-dot.ended {
  background: var(--color-text-muted);
}

.campaign-dot.paused {
  background: var(--color-warning);
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.4);
}

.status-tag {
  border-radius: 20px;
  padding: 2px 12px;
}

.budget-value {
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
}

.btn-action {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  border-radius: var(--radius-md);
  padding: var(--spacing-1) var(--spacing-3);
  transition: all var(--transition-fast);
}

.btn-edit {
  color: var(--color-accent);
  border-color: rgba(6, 182, 212, 0.3);
}

.btn-edit:hover {
  background: rgba(6, 182, 212, 0.1);
  border-color: var(--color-accent);
}

.btn-delete {
  color: var(--color-danger);
  border-color: rgba(239, 68, 68, 0.3);
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: var(--color-danger);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-16);
  color: var(--color-text-muted);
}

.empty-state svg {
  margin-bottom: var(--spacing-4);
  opacity: 0.5;
}

.empty-state p {
  margin: 0 0 var(--spacing-6);
}

.custom-dialog :deep(.el-dialog__header) {
  border-bottom: 1px solid var(--color-border-light);
}

.custom-dialog :deep(.el-dialog__title) {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
}

.custom-dialog :deep(.el-dialog__body) {
  padding: var(--spacing-6);
}

.custom-dialog :deep(.el-dialog__footer) {
  border-top: 1px solid var(--color-border-light);
  padding: var(--spacing-4) var(--spacing-6);
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
  }
  
  .search-input {
    width: 100%;
  }
  
  .stats-row {
    flex-wrap: wrap;
  }
  
  .mini-stat {
    flex: 1;
    min-width: 120px;
  }
}
</style>
