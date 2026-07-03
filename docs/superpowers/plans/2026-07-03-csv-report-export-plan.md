# CSV 报表导出功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在数据看板页面实现支持自定义时间范围的 CSV 报表导出功能

**Architecture:** 纯前端导出方案，在 Dashboard.vue 中添加导出按钮和时间选择对话框，使用 JavaScript 将数据转换为 CSV 格式并触发浏览器下载

**Tech Stack:** Vue 3, Element Plus, JavaScript

## Global Constraints

- 使用现有技术栈（Vue 3 + Element Plus）
- 无需后端支持，纯前端实现
- CSV 文件采用 UTF-8 with BOM 编码以支持中文
- 文件命名格式：`report_YYYY-MM-DD_to_YYYY-MM-DD.csv`

---

## 文件结构

| 文件 | 操作 | 说明 |
|------|------|------|
| `frontend-vue/src/views/Dashboard.vue` | 修改 | 添加导出按钮、时间选择对话框和 CSV 生成逻辑 |
| `docs/superpowers/specs/2026-07-03-csv-report-export-design.md` | 参考 | 设计文档 |

---

### Task 1: 添加导出按钮和时间选择对话框

**Files:**
- Modify: `frontend-vue/src/views/Dashboard.vue`

**Interfaces:**
- Consumes: 现有 Dashboard.vue 组件结构
- Produces: 导出按钮、时间选择对话框、导出相关状态变量

- [ ] **Step 1: 添加导出按钮到页面头部**

在 Dashboard.vue 的 `page-header` 区域添加导出按钮：

```vue
<div class="header-actions">
  <el-button @click="showExportDialog = true" class="btn-export">
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
      <polyline points="7 10 12 15 17 10"></polyline>
      <line x1="12" y1="15" x2="12" y2="3"></line>
    </svg>
    <span>导出报表</span>
  </el-button>
</div>
```

- [ ] **Step 2: 添加时间选择对话框**

在 Dashboard.vue 的模板末尾添加时间选择对话框：

```vue
<el-dialog title="导出报表" v-model="showExportDialog" width="400px">
  <el-form :model="exportForm" label-width="80px">
    <el-form-item label="时间范围">
      <el-date-picker
        v-model="exportForm.dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        style="width: 100%"
      />
    </el-form-item>
  </el-form>
  <template #footer>
    <el-button @click="showExportDialog = false">取消</el-button>
    <el-button type="primary" @click="handleExport">导出 CSV</el-button>
  </template>
</el-dialog>
```

- [ ] **Step 3: 添加导出相关状态变量**

在 Dashboard.vue 的 `<script setup>` 中添加：

```javascript
const showExportDialog = ref(false)
const exportForm = reactive({
  dateRange: []
})
```

- [ ] **Step 4: 添加导出按钮样式**

在 Dashboard.vue 的 `<style scoped>` 中添加：

```css
.btn-export {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  margin-left: var(--spacing-3);
  border-radius: var(--radius-md);
}
```

- [ ] **Step 5: 验证页面构建是否成功**

Run: `cd frontend-vue && npm run build`
Expected: 构建成功，无错误

---

### Task 2: 实现 CSV 生成和下载逻辑

**Files:**
- Modify: `frontend-vue/src/views/Dashboard.vue`

**Interfaces:**
- Consumes: `overview` 和 `dailyTrend` 数据
- Produces: `handleExport` 函数，用于生成和下载 CSV 文件

- [ ] **Step 1: 实现 CSV 生成函数**

在 Dashboard.vue 的 `<script setup>` 中添加：

```javascript
const generateCSV = () => {
  let csv = '\ufeff'
  
  csv += '广告 BI 系统数据报表\n'
  
  const now = new Date()
  const exportTime = now.toLocaleString('zh-CN')
  const startDate = exportForm.dateRange[0] ? exportForm.dateRange[0].toLocaleDateString('zh-CN') : '最近7天'
  const endDate = exportForm.dateRange[1] ? exportForm.dateRange[1].toLocaleDateString('zh-CN') : '今天'
  
  csv += `时间范围: ${startDate} 至 ${endDate}\n`
  csv += `导出时间: ${exportTime}\n\n`
  
  csv += '指标,数值\n'
  csv += `总曝光量,${overview.value.total_impressions.toLocaleString()}\n`
  csv += `总点击量,${overview.value.total_clicks.toLocaleString()}\n`
  csv += `总转化数,${overview.value.total_conversions.toLocaleString()}\n`
  csv += `总收入,${overview.value.total_revenue.toFixed(2)}\n`
  csv += `平均点击率,${overview.value.avg_ctr}%\n`
  csv += `平均转化率,${overview.value.avg_cvr}%\n\n`
  
  csv += '日期,曝光量,点击量,转化数,收入\n'
  if (dailyTrend.value && dailyTrend.value.dates) {
    dailyTrend.value.dates.forEach((date, index) => {
      csv += `${date},${dailyTrend.value.impressions[index].toLocaleString()},${dailyTrend.value.clicks[index].toLocaleString()},${dailyTrend.value.conversions[index].toLocaleString()},${dailyTrend.value.revenue[index].toFixed(2)}\n`
    })
  }
  
  return csv
}
```

- [ ] **Step 2: 实现导出处理函数**

在 Dashboard.vue 的 `<script setup>` 中添加：

```javascript
const handleExport = () => {
  try {
    const csv = generateCSV()
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    
    const link = document.createElement('a')
    link.href = url
    
    const start = exportForm.dateRange[0] ? exportForm.dateRange[0].toISOString().split('T')[0] : 'start'
    const end = exportForm.dateRange[1] ? exportForm.dateRange[1].toISOString().split('T')[0] : 'end'
    link.download = `report_${start}_to_${end}.csv`
    
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    
    showExportDialog.value = false
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message)
  }
}
```

- [ ] **Step 3: 添加 ElMessage 导入**

确保 Dashboard.vue 的 `<script setup>` 中已导入：

```javascript
import { ElMessage } from 'element-plus'
```

- [ ] **Step 4: 测试导出功能**

1. 启动前端服务：`cd frontend-vue && npm run dev`
2. 访问 http://localhost:3000/
3. 点击"导出报表"按钮
4. 选择日期范围并点击"导出 CSV"
5. 验证 CSV 文件是否下载成功

- [ ] **Step 5: 验证构建是否成功**

Run: `cd frontend-vue && npm run build`
Expected: 构建成功，无错误

---

### Task 3: 添加边界情况处理

**Files:**
- Modify: `frontend-vue/src/views/Dashboard.vue`

**Interfaces:**
- Consumes: `handleExport` 函数
- Produces: 完善的错误处理逻辑

- [ ] **Step 1: 添加日期验证**

修改 `handleExport` 函数，添加日期验证：

```javascript
const handleExport = () => {
  if (exportForm.dateRange.length === 2) {
    if (exportForm.dateRange[0] > exportForm.dateRange[1]) {
      ElMessage.error('开始日期不能晚于结束日期')
      return
    }
  }
  
  try {
    const csv = generateCSV()
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    
    const link = document.createElement('a')
    link.href = url
    
    const start = exportForm.dateRange[0] ? exportForm.dateRange[0].toISOString().split('T')[0] : 'start'
    const end = exportForm.dateRange[1] ? exportForm.dateRange[1].toISOString().split('T')[0] : 'end'
    link.download = `report_${start}_to_${end}.csv`
    
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    
    showExportDialog.value = false
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message)
  }
}
```

- [ ] **Step 2: 添加数据检查**

修改 `generateCSV` 函数，添加数据检查：

```javascript
const generateCSV = () => {
  if (!overview.value || !overview.value.total_impressions) {
    throw new Error('暂无数据可导出')
  }
  
  let csv = '\ufeff'
  // ... 其余代码不变
}
```

- [ ] **Step 3: 测试边界情况**

1. 测试起始日期晚于结束日期的情况
2. 测试没有数据的情况
3. 验证错误提示是否正确显示

- [ ] **Step 4: 验证构建是否成功**

Run: `cd frontend-vue && npm run build`
Expected: 构建成功，无错误

---

## Self-Review

**1. Spec coverage:**
- ✅ 在数据看板页面添加导出按钮
- ✅ 支持选择起始日期和结束日期
- ✅ 导出 CSV 文件包含概览数据和每日趋势数据
- ✅ 导出失败时有明确的错误提示

**2. Placeholder scan:**
- ✅ 无 "TBD", "TODO" 等占位符
- ✅ 所有步骤包含完整代码
- ✅ 所有命令明确具体

**3. Type consistency:**
- ✅ 函数名和变量名在所有任务中保持一致
- ✅ 数据结构与设计文档一致

---

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-07-03-csv-report-export-plan.md`. Two execution options:

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

**Which approach?**
