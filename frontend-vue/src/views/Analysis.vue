<template>
  <div class="analysis-container">
    <div class="page-header">
      <div class="header-info">
        <h1 class="page-title">用户行为分析</h1>
        <p class="page-subtitle">深入分析用户交互行为和转化路径</p>
      </div>
      <div class="header-actions">
        <el-button type="default" @click="handleExport" class="btn-export">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          <span>导出报告</span>
        </el-button>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card">
        <div class="card-header">
          <h3 class="card-title">用户来源分布</h3>
        </div>
        <div ref="sourceChartRef" class="chart-container"></div>
      </div>

      <div class="chart-card">
        <div class="card-header">
          <h3 class="card-title">设备类型占比</h3>
        </div>
        <div ref="deviceChartRef" class="chart-container"></div>
      </div>

      <div class="chart-card large">
        <div class="card-header">
          <h3 class="card-title">用户活跃趋势</h3>
          <div class="card-actions">
            <el-button size="small" type="default" :class="{ active: activePeriod === 'day' }" @click="activePeriod = 'day'">日</el-button>
            <el-button size="small" type="default" :class="{ active: activePeriod === 'week' }" @click="activePeriod = 'week'">周</el-button>
            <el-button size="small" type="default" :class="{ active: activePeriod === 'month' }" @click="activePeriod = 'month'">月</el-button>
          </div>
        </div>
        <div ref="activeChartRef" class="chart-container"></div>
      </div>

      <div class="chart-card">
        <div class="card-header">
          <h3 class="card-title">用户行为排行</h3>
        </div>
        <div class="ranking-list">
          <div class="ranking-item" v-for="(item, index) in behaviorRanking" :key="item.name">
            <div class="rank-badge" :class="'rank-' + (index + 1)">
              {{ index + 1 }}
            </div>
            <div class="rank-info">
              <span class="rank-name">{{ item.name }}</span>
              <div class="rank-bar">
                <div class="rank-bar-fill" :style="{ width: item.percent + '%' }"></div>
              </div>
            </div>
            <span class="rank-value">{{ item.count }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'

const activePeriod = ref('week')

const behaviorRanking = ref([
  { name: '页面浏览', count: 156800, percent: 100 },
  { name: '点击广告', count: 124500, percent: 79.4 },
  { name: '商品详情', count: 89200, percent: 56.9 },
  { name: '加入购物车', count: 45600, percent: 29.1 },
  { name: '完成购买', count: 23800, percent: 15.2 }
])

const sourceChartRef = ref(null)
const deviceChartRef = ref(null)
const activeChartRef = ref(null)
let sourceChart = null
let deviceChart = null
let activeChart = null

const initSourceChart = () => {
  if (!sourceChartRef.value) return
  sourceChart = echarts.init(sourceChartRef.value)
  
  sourceChart.setOption({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(15, 23, 42, 0.95)',
      borderColor: 'transparent',
      borderRadius: 8,
      textStyle: { color: '#ffffff' },
      formatter: '{b}: {c} ({d}%)'
    },
    series: [
      {
        name: '用户来源',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 6,
          borderColor: '#ffffff',
          borderWidth: 2
        },
        label: {
          show: true,
          fontSize: 12,
          color: '#64748b'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold'
          }
        },
        data: [
          { value: 45000, name: '搜索引擎', itemStyle: { color: '#06b6d4' } },
          { value: 32000, name: '社交媒体', itemStyle: { color: '#3b82f6' } },
          { value: 28000, name: '直接访问', itemStyle: { color: '#10b981' } },
          { value: 18000, name: '邮件营销', itemStyle: { color: '#f59e0b' } },
          { value: 12000, name: '其他', itemStyle: { color: '#94a3b8' } }
        ]
      }
    ]
  })
}

const initDeviceChart = () => {
  if (!deviceChartRef.value) return
  deviceChart = echarts.init(deviceChartRef.value)
  
  deviceChart.setOption({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(15, 23, 42, 0.95)',
      borderColor: 'transparent',
      borderRadius: 8,
      textStyle: { color: '#ffffff' },
      formatter: '{b}: {c} ({d}%)'
    },
    series: [
      {
        name: '设备类型',
        type: 'pie',
        radius: '65%',
        center: ['50%', '50%'],
        itemStyle: {
          borderRadius: 6,
          borderColor: '#ffffff',
          borderWidth: 2
        },
        label: {
          show: true,
          fontSize: 12,
          color: '#64748b'
        },
        data: [
          { value: 68000, name: '移动端', itemStyle: { color: '#06b6d4' } },
          { value: 52000, name: 'PC端', itemStyle: { color: '#3b82f6' } },
          { value: 15000, name: '平板', itemStyle: { color: '#8b5cf6' } }
        ]
      }
    ]
  })
}

const initActiveChart = () => {
  if (!activeChartRef.value) return
  activeChart = echarts.init(activeChartRef.value)
  
  const dates = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  const activeUsers = [12500, 13800, 14200, 13500, 15600, 18900, 17200]
  const newUsers = [2300, 2600, 2800, 2500, 3100, 4200, 3800]
  
  activeChart.setOption({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.95)',
      borderColor: 'transparent',
      borderRadius: 8,
      textStyle: { color: '#ffffff' }
    },
    legend: {
      data: ['活跃用户', '新用户'],
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
        name: '活跃用户',
        type: 'bar',
        data: activeUsers,
        barWidth: '40%',
        itemStyle: {
          borderRadius: [4, 4, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#06b6d4' },
            { offset: 1, color: '#22d3ee' }
          ])
        }
      },
      {
        name: '新用户',
        type: 'bar',
        data: newUsers,
        barWidth: '40%',
        itemStyle: {
          borderRadius: [4, 4, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#8b5cf6' },
            { offset: 1, color: '#a855f7' }
          ])
        }
      }
    ]
  })
}

const handleExport = () => {}

onMounted(() => {
  initSourceChart()
  initDeviceChart()
  initActiveChart()
  
  window.addEventListener('resize', () => {
    sourceChart?.resize()
    deviceChart?.resize()
    activeChart?.resize()
  })
})
</script>

<style scoped>
.analysis-container {
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

.btn-export {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-normal);
}

.btn-export:hover {
  background: var(--color-bg);
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-6);
}

.chart-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.chart-card.large {
  grid-column: span 2;
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
  height: 320px;
}

.ranking-list {
  padding: var(--spacing-6);
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4) 0;
  border-bottom: 1px solid var(--color-border-light);
}

.ranking-item:last-child {
  border-bottom: none;
}

.rank-badge {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-bold);
  color: white;
  flex-shrink: 0;
}

.rank-1 {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
}

.rank-2 {
  background: linear-gradient(135deg, #94a3b8 0%, #cbd5e1 100%);
}

.rank-3 {
  background: linear-gradient(135deg, #d97706 0%, #f59e0b 100%);
}

.rank-4, .rank-5 {
  background: var(--color-border);
  color: var(--color-text-secondary);
}

.rank-info {
  flex: 1;
}

.rank-name {
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  display: block;
  margin-bottom: var(--spacing-2);
}

.rank-bar {
  height: 6px;
  background: var(--color-bg);
  border-radius: 3px;
  overflow: hidden;
}

.rank-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border-radius: 3px;
  transition: width 0.6s ease-out;
}

.rank-value {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
  min-width: 80px;
  text-align: right;
}

@media (max-width: 768px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .chart-card.large {
    grid-column: span 1;
  }
  
  .chart-container {
    height: 280px;
  }
}
</style>