<template>
  <div class="dashboard-layout">
    <!-- 侧边栏 -->
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h3>广告 BI 系统</h3>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据看板</span>
        </el-menu-item>
        <el-menu-item index="/campaigns">
          <el-icon><List /></el-icon>
          <span>广告活动</span>
        </el-menu-item>
        <el-menu-item index="/ad-units">
          <el-icon><Grid /></el-icon>
          <span>广告单元</span>
        </el-menu-item>
        <el-menu-item index="/analysis">
          <el-icon><TrendCharts /></el-icon>
          <span>用户分析</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <h2>数据看板</h2>
        </div>
        <div class="header-right">
          <span class="username">{{ username }}</span>
          <el-button type="danger" size="small" @click="handleLogout">
            退出登录
          </el-button>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="main-content">
        <!-- 数据卡片 -->
        <el-row :gutter="20" class="stats-cards">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-item">
                <div class="stat-label">总曝光量</div>
                <div class="stat-value">{{ overview.total_impressions.toLocaleString() }}</div>
                <div class="stat-icon impressions">👁</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-item">
                <div class="stat-label">总点击量</div>
                <div class="stat-value">{{ overview.total_clicks.toLocaleString() }}</div>
                <div class="stat-icon clicks">👆</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-item">
                <div class="stat-label">总转化量</div>
                <div class="stat-value">{{ overview.total_conversions.toLocaleString() }}</div>
                <div class="stat-icon conversions">✅</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-item">
                <div class="stat-label">总收入</div>
                <div class="stat-value">¥{{ overview.total_revenue.toLocaleString() }}</div>
                <div class="stat-icon revenue">💰</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 图表区 -->
        <el-row :gutter="20" class="charts-row">
          <el-col :span="16">
            <el-card>
              <template #header>
                <div class="card-header">
                  <span>每日趋势</span>
                </div>
              </template>
              <div ref="trendChartRef" style="height: 400px"></div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card>
              <template #header>
                <div class="card-header">
                  <span>转化漏斗</span>
                </div>
              </template>
              <div ref="funnelChartRef" style="height: 400px"></div>
            </el-card>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getOverview, getDailyTrend, getConversionFunnel } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '用户')
const activeMenu = ref('/')

const overview = reactive({
  total_impressions: 0,
  total_clicks: 0,
  total_conversions: 0,
  total_cost: 0,
  total_revenue: 0,
  average_ctr: 0,
  average_conversion_rate: 0
})

const trendChartRef = ref(null)
const funnelChartRef = ref(null)
let trendChart = null
let funnelChart = null

// 加载概览数据
const loadOverview = async () => {
  try {
    const res = await getOverview()
    Object.assign(overview, res.data)
  } catch (error) {
    console.error('加载概览数据失败:', error)
  }
}

// 加载趋势图表
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
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['曝光量', '点击量']
      },
      xAxis: {
        type: 'category',
        data: dates
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '曝光量',
          type: 'line',
          data: impressions,
          smooth: true,
          itemStyle: { color: '#667eea' }
        },
        {
          name: '点击量',
          type: 'line',
          data: clicks,
          smooth: true,
          itemStyle: { color: '#764ba2' }
        }
      ]
    })
  } catch (error) {
    console.error('加载趋势图表失败:', error)
  }
}

// 加载漏斗图表
const loadFunnelChart = async () => {
  try {
    const res = await getConversionFunnel()
    const data = res.data
    
    if (!funnelChartRef.value) return
    funnelChart = echarts.init(funnelChartRef.value)
    
    funnelChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c}'
      },
      series: [
        {
          name: '转化漏斗',
          type: 'funnel',
          left: '10%',
          top: 60,
          bottom: 60,
          width: '80%',
          min: 0,
          max: 10000,
          minSize: '0%',
          maxSize: '100%',
          sort: 'descending',
          gap: 2,
          label: {
            show: true,
            position: 'inside'
          },
          data: data.map(item => ({
            name: item.step_name,
            value: item.user_count
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

onMounted(() => {
  loadOverview()
  loadTrendChart()
  loadFunnelChart()
  
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
  background: #304156;
  color: #fff;
}

.logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.logo h3 {
  color: #fff;
  margin: 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: #606266;
}

.main-content {
  background: #f5f7fa;
  padding: 20px;
}

.stats-cards {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  position: relative;
  padding: 10px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-icon {
  position: absolute;
  right: 20px;
  top: 20px;
  font-size: 40px;
  opacity: 0.3;
}

.charts-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
