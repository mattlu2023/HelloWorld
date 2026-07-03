<template>
  <div class="dashboard-layout">
    <el-aside width="260px" class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <line x1="9" y1="9" x2="15" y2="9"></line>
            <line x1="12" y1="12" x2="12" y2="15"></line>
          </svg>
          <span class="logo-text">广告 BI</span>
        </div>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        background-color="transparent"
        text-color="#cbd5e1"
        active-text-color="#ffffff"
        router
        class="sidebar-menu"
      >
        <el-menu-item index="/">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
          </svg>
          <span>数据看板</span>
        </el-menu-item>
        <el-menu-item index="/campaigns">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="8" y1="6" x2="21" y2="6"></line>
            <line x1="8" y1="12" x2="21" y2="12"></line>
            <line x1="8" y1="18" x2="21" y2="18"></line>
            <line x1="3" y1="6" x2="3.01" y2="6"></line>
            <line x1="3" y1="12" x2="3.01" y2="12"></line>
            <line x1="3" y1="18" x2="3.01" y2="18"></line>
          </svg>
          <span>广告活动</span>
        </el-menu-item>
        <el-menu-item index="/ad-units">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="7" height="7"></rect>
            <rect x="14" y="3" width="7" height="7"></rect>
            <rect x="14" y="14" width="7" height="7"></rect>
            <rect x="3" y="14" width="7" height="7"></rect>
          </svg>
          <span>广告单元</span>
        </el-menu-item>
        <el-menu-item index="/analysis">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>
          </svg>
          <span>用户分析</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-container">
      <el-header class="header">
        <div class="header-left">
          <h1 class="page-title">{{ currentPageTitle }}</h1>
          <p class="page-subtitle">{{ currentPageSubtitle }}</p>
        </div>
        <div class="header-right">
          <el-button type="primary" size="small" @click="showExportDialog = true" class="export-btn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            <span>导出报表</span>
          </el-button>
          <div class="user-info">
            <div class="user-avatar">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </div>
            <span class="username">{{ username }}</span>
          </div>
          <el-button type="default" size="small" @click="handleLogout" class="logout-btn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            <span>退出</span>
          </el-button>
        </div>
      </el-header>

      <el-main class="main-content">
        <div class="stats-grid">
          <div class="stat-card stat-card-impressions">
            <div class="stat-icon-wrapper">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path>
                <circle cx="12" cy="12" r="3"></circle>
              </svg>
            </div>
            <div class="stat-content">
              <p class="stat-label">总曝光量</p>
              <p class="stat-value">{{ overview.total_impressions.toLocaleString() }}</p>
              <div class="stat-trend positive">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 10 10 20 4 14"></polyline>
                </svg>
                <span>+12.5%</span>
              </div>
            </div>
          </div>

          <div class="stat-card stat-card-clicks">
            <div class="stat-icon-wrapper">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M15 18l-6-6 6-6"></path>
              </svg>
            </div>
            <div class="stat-content">
              <p class="stat-label">总点击量</p>
              <p class="stat-value">{{ overview.total_clicks.toLocaleString() }}</p>
              <div class="stat-trend positive">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 10 10 20 4 14"></polyline>
                </svg>
                <span>+8.3%</span>
              </div>
            </div>
          </div>

          <div class="stat-card stat-card-conversions">
            <div class="stat-icon-wrapper">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
            </div>
            <div class="stat-content">
              <p class="stat-label">总转化量</p>
              <p class="stat-value">{{ overview.total_conversions.toLocaleString() }}</p>
              <div class="stat-trend positive">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 10 10 20 4 14"></polyline>
                </svg>
                <span>+15.2%</span>
              </div>
            </div>
          </div>

          <div class="stat-card stat-card-revenue">
            <div class="stat-icon-wrapper">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="12" y1="1" x2="12" y2="23"></line>
                <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
              </svg>
            </div>
            <div class="stat-content">
              <p class="stat-label">总收入</p>
              <p class="stat-value">¥{{ overview.total_revenue.toLocaleString() }}</p>
              <div class="stat-trend positive">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 10 10 20 4 14"></polyline>
                </svg>
                <span>+22.1%</span>
              </div>
            </div>
          </div>
        </div>

        <div class="charts-section">
          <div class="chart-card large">
            <div class="card-header">
              <h3 class="card-title">每日趋势</h3>
              <div class="card-actions">
                <el-button size="small" type="default" :class="{ active: timeRange === 'week' }" @click="timeRange = 'week'">本周</el-button>
                <el-button size="small" type="default" :class="{ active: timeRange === 'month' }" @click="timeRange = 'month'">本月</el-button>
                <el-button size="small" type="default" :class="{ active: timeRange === 'quarter' }" @click="timeRange = 'quarter'">本季度</el-button>
              </div>
            </div>
            <div ref="trendChartRef" class="chart-container"></div>
          </div>

          <div class="chart-card">
            <div class="card-header">
              <h3 class="card-title">转化漏斗</h3>
            </div>
            <div ref="funnelChartRef" class="chart-container"></div>
          </div>
        </div>
      </el-main>
    </el-container>

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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getOverview, getDailyTrend, getConversionFunnel } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '用户')
const activeMenu = ref('/')
const timeRange = ref('week')

const showExportDialog = ref(false)
const exportForm = reactive({
  dateRange: []
})

const overview = reactive({
  total_impressions: 0,
  total_clicks: 0,
  total_conversions: 0,
  total_cost: 0,
  total_revenue: 0,
  average_ctr: 0,
  average_conversion_rate: 0
})

const currentPageTitle = computed(() => '数据看板')
const currentPageSubtitle = computed(() => '实时监控广告投放效果')

const trendChartRef = ref(null)
const funnelChartRef = ref(null)
let trendChart = null
let funnelChart = null

const loadOverview = async () => {
  try {
    const res = await getOverview()
    Object.assign(overview, res.data)
  } catch (error) {
    console.error('加载概览数据失败:', error)
  }
}

const loadTrendChart = async () => {
  try {
    const res = await getDailyTrend()
    const data = res.data.reverse()
    
    const dates = data.map(item => item.date)
    const impressions = data.map(item => item.impressions)
    const clicks = data.map(item => item.clicks)
    
    if (!trendChartRef.value) return
    trendChart = echarts.init(trendChartRef.value)
    
    trendChart.setOption({
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(15, 23, 42, 0.95)',
        borderColor: 'transparent',
        borderRadius: 8,
        textStyle: { color: '#ffffff' }
      },
      legend: {
        data: ['曝光量', '点击量'],
        textStyle: { color: '#64748b' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dates,
        axisLine: { lineStyle: { color: '#e2e8f0' } },
        axisLabel: { color: '#94a3b8' }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f1f5f9' } },
        axisLabel: { color: '#94a3b8' }
      },
      series: [
        {
          name: '曝光量',
          type: 'line',
          data: impressions,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(6, 182, 212, 0.2)' },
              { offset: 1, color: 'rgba(6, 182, 212, 0)' }
            ])
          },
          lineStyle: { color: '#06b6d4', width: 3 },
          itemStyle: { color: '#06b6d4' }
        },
        {
          name: '点击量',
          type: 'line',
          data: clicks,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(59, 130, 246, 0.2)' },
              { offset: 1, color: 'rgba(59, 130, 246, 0)' }
            ])
          },
          lineStyle: { color: '#3b82f6', width: 3 },
          itemStyle: { color: '#3b82f6' }
        }
      ]
    })
  } catch (error) {
    console.error('加载趋势图表失败:', error)
  }
}

const loadFunnelChart = async () => {
  try {
    const res = await getConversionFunnel()
    const data = res.data
    
    if (!funnelChartRef.value) return
    funnelChart = echarts.init(funnelChartRef.value)
    
    funnelChart.setOption({
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(15, 23, 42, 0.95)',
        borderColor: 'transparent',
        borderRadius: 8,
        textStyle: { color: '#ffffff' },
        formatter: '{b}: {c}'
      },
      series: [
        {
          name: '转化漏斗',
          type: 'funnel',
          left: '10%',
          top: '10%',
          bottom: '10%',
          width: '80%',
          min: 0,
          max: 10000,
          minSize: '0%',
          maxSize: '100%',
          sort: 'descending',
          gap: 4,
          label: {
            show: true,
            position: 'inside',
            color: '#ffffff',
            fontSize: 12,
            fontWeight: 600
          },
          itemStyle: {
            borderRadius: 6,
            borderColor: '#ffffff',
            borderWidth: 2
          },
          emphasis: {
            label: { fontSize: 14 }
          },
          data: data.map((item, index) => ({
            name: item.step_name,
            value: item.user_count,
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
                { offset: 0, color: ['#06b6d4', '#22d3ee', '#3b82f6', '#8b5cf6'][index] || '#06b6d4' },
                { offset: 1, color: ['#22d3ee', '#3b82f6', '#8b5cf6', '#a855f7'][index] || '#22d3ee' }
              ])
            }
          }))
        }
      ]
    })
  } catch (error) {
    console.error('加载漏斗图表失败:', error)
  }
}

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  ElMessage.success('已退出登录')
  router.push('/login')
}

const dailyTrend = reactive({
  dates: [],
  impressions: [],
  clicks: [],
  conversions: [],
  revenue: []
})

const loadDailyTrendData = async () => {
  try {
    const res = await getDailyTrend()
    const data = res.data.reverse()
    dailyTrend.dates = data.map(item => item.date)
    dailyTrend.impressions = data.map(item => item.impressions)
    dailyTrend.clicks = data.map(item => item.clicks)
    dailyTrend.conversions = data.map(item => item.conversions || 0)
    dailyTrend.revenue = data.map(item => item.revenue || 0)
  } catch (error) {
    console.error('加载每日趋势数据失败:', error)
  }
}

const generateCSV = () => {
  if (!overview.total_impressions) {
    throw new Error('暂无数据可导出')
  }

  let csv = '\ufeff'

  csv += '广告 BI 系统数据报表\n'

  const now = new Date()
  const exportTime = now.toLocaleString('zh-CN')
  const startDate = exportForm.dateRange[0] ? exportForm.dateRange[0].toLocaleDateString('zh-CN') : '最近7天'
  const endDate = exportForm.dateRange[1] ? exportForm.dateRange[1].toLocaleDateString('zh-CN') : '今天'

  csv += `时间范围: ${startDate} 至 ${endDate}\n`
  csv += `导出时间: ${exportTime}\n\n`

  csv += '指标,数值\n'
  csv += `总曝光量,${overview.total_impressions.toLocaleString()}\n`
  csv += `总点击量,${overview.total_clicks.toLocaleString()}\n`
  csv += `总转化数,${overview.total_conversions.toLocaleString()}\n`
  csv += `总收入,${overview.total_revenue.toFixed(2)}\n`
  csv += `平均点击率,${overview.average_ctr}%\n`
  csv += `平均转化率,${overview.average_conversion_rate}%\n\n`

  csv += '日期,曝光量,点击量,转化数,收入\n'
  if (dailyTrend.dates && dailyTrend.dates.length > 0) {
    dailyTrend.dates.forEach((date, index) => {
      csv += `${date},${dailyTrend.impressions[index].toLocaleString()},${dailyTrend.clicks[index].toLocaleString()},${dailyTrend.conversions[index].toLocaleString()},${dailyTrend.revenue[index].toFixed(2)}\n`
    })
  }

  return csv
}

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

onMounted(() => {
  loadOverview()
  loadTrendChart()
  loadFunnelChart()
  loadDailyTrendData()
  
  window.addEventListener('resize', () => {
    trendChart?.resize()
    funnelChart?.resize()
  })
})
</script>

<style scoped>
.dashboard-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-primary-light) 100%);
  color: white;
  box-shadow: var(--shadow-xl);
}

.sidebar-header {
  padding: var(--spacing-6) var(--spacing-4);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  color: var(--color-accent);
}

.logo-text {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  letter-spacing: -0.02em;
}

.sidebar-menu {
  padding: var(--spacing-4) 0;
}

.main-container {
  flex: 1;
  min-height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border-light);
  padding: 0 var(--spacing-8);
  height: var(--header-height);
}

.header-left {
  display: flex;
  flex-direction: column;
}

.page-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  margin: 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-2) var(--spacing-4);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
}

.user-avatar {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border-radius: 50%;
  color: white;
}

.username {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}

.logout-btn:hover {
  color: var(--color-danger);
  background-color: rgba(239, 68, 68, 0.1);
}

.main-content {
  background: var(--color-bg);
  padding: var(--spacing-8);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--spacing-6);
  margin-bottom: var(--spacing-8);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-6);
  padding: var(--spacing-6);
  border-radius: var(--radius-xl);
  background: var(--color-bg-card);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-normal);
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
}

.stat-card-impressions::before {
  background: linear-gradient(90deg, #06b6d4, #22d3ee);
}

.stat-card-clicks::before {
  background: linear-gradient(90deg, #3b82f6, #60a5fa);
}

.stat-card-conversions::before {
  background: linear-gradient(90deg, #10b981, #34d399);
}

.stat-card-revenue::before {
  background: linear-gradient(90deg, #f59e0b, #fbbf24);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.stat-icon-wrapper {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-lg);
  flex-shrink: 0;
}

.stat-card-impressions .stat-icon-wrapper {
  background: rgba(6, 182, 212, 0.1);
  color: #06b6d4;
}

.stat-card-clicks .stat-icon-wrapper {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.stat-card-conversions .stat-icon-wrapper {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.stat-card-revenue .stat-icon-wrapper {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--spacing-1);
}

.stat-value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  margin: 0 0 var(--spacing-2);
  letter-spacing: -0.02em;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.stat-trend.positive {
  color: var(--color-success);
}

.stat-trend.negative {
  color: var(--color-danger);
}

.charts-section {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: var(--spacing-6);
}

.chart-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.chart-card.large {
  grid-column: span 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-6);
  border-bottom: 1px solid var(--color-border-light);
}

.card-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
  margin: 0;
}

.card-actions {
  display: flex;
  gap: var(--spacing-2);
}

.card-actions .el-button {
  border-radius: var(--radius-md);
  font-size: var(--font-size-xs);
  padding: var(--spacing-1) var(--spacing-3);
}

.card-actions .el-button.active {
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border: none;
  color: white;
}

.chart-container {
  padding: var(--spacing-6);
  height: 400px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-section {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.export-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  border-radius: var(--radius-md);
}
</style>