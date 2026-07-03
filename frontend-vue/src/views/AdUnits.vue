<template>
  <div class="adunits-container">
    <div class="page-header">
      <div class="header-info">
        <h1 class="page-title">广告单元管理</h1>
        <p class="page-subtitle">管理和配置广告展示单元</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd" class="btn-primary">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          <span>创建单元</span>
        </el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索单元名称..."
        class="search-input"
        prefix-icon="Search"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="typeFilter" placeholder="类型筛选" class="filter-select">
        <el-option label="全部" value="" />
        <el-option label="横幅广告" value="banner" />
        <el-option label="信息流" value="feed" />
        <el-option label="弹窗广告" value="popup" />
      </el-select>
    </div>

    <div class="adunits-grid">
      <div class="adunit-card" v-for="unit in filteredAdUnits" :key="unit.id">
        <div class="card-header">
          <div class="card-icon" :class="unit.type">
            <svg v-if="unit.type === 'banner'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2"></rect>
            </svg>
            <svg v-else-if="unit.type === 'feed'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M8 16h1"></path>
              <path d="M8 12h1"></path>
              <path d="M8 8h1"></path>
              <path d="M12 16h.01"></path>
              <path d="M12 12h.01"></path>
              <path d="M12 8h.01"></path>
              <path d="M16 16h1"></path>
              <path d="M16 12h1"></path>
              <path d="M16 8h1"></path>
            </svg>
            <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="11" width="18" height="11" rx="2"></rect>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
          </div>
          <el-tag :type="getTypeTag(unit.type)" size="small" class="type-tag">
            {{ getTypeText(unit.type) }}
          </el-tag>
        </div>
        
        <h3 class="card-title">{{ unit.name }}</h3>
        <p class="card-desc">{{ unit.description }}</p>
        
        <div class="card-stats">
          <div class="stat-item">
            <span class="stat-value">{{ unit.impressions }}</span>
            <span class="stat-label">曝光</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ unit.clicks }}</span>
            <span class="stat-label">点击</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ unit.ctr }}%</span>
            <span class="stat-label">CTR</span>
          </div>
        </div>
        
        <div class="card-footer">
          <span class="card-date">创建于 {{ unit.created_at }}</span>
          <div class="card-actions">
            <el-button size="small" class="btn-icon" title="编辑">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path>
              </svg>
            </el-button>
            <el-button size="small" class="btn-icon btn-danger" title="删除">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="empty-state" v-if="filteredAdUnits.length === 0">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="7" height="7"></rect>
        <rect x="14" y="3" width="7" height="7"></rect>
        <rect x="14" y="14" width="7" height="7"></rect>
        <rect x="3" y="14" width="7" height="7"></rect>
      </svg>
      <p>暂无广告单元</p>
      <el-button type="primary" @click="handleAdd">创建第一个单元</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')
const typeFilter = ref('')

const adUnits = ref([
  { id: 1, name: '首页横幅广告', description: '网站首页顶部横幅位，展示主推产品', type: 'banner', impressions: 125000, clicks: 8320, ctr: 6.65, created_at: '2024-01-15' },
  { id: 2, name: '信息流广告位', description: '内容信息流中的原生广告', type: 'feed', impressions: 89000, clicks: 6150, ctr: 6.91, created_at: '2024-02-20' },
  { id: 3, name: '弹窗促销广告', description: '节假日促销活动弹窗', type: 'popup', impressions: 45000, clicks: 3280, ctr: 7.29, created_at: '2024-03-10' },
  { id: 4, name: '侧边栏横幅', description: '页面侧边栏固定广告位', type: 'banner', impressions: 67000, clicks: 4520, ctr: 6.75, created_at: '2024-04-05' },
  { id: 5, name: 'APP信息流', description: '移动端APP信息流广告', type: 'feed', impressions: 156000, clicks: 10890, ctr: 6.98, created_at: '2024-05-12' },
  { id: 6, name: '登录弹窗广告', description: '用户登录后弹窗推荐', type: 'popup', impressions: 32000, clicks: 2150, ctr: 6.72, created_at: '2024-06-20' }
])

const filteredAdUnits = computed(() => {
  return adUnits.value.filter(item => {
    const matchSearch = !searchQuery.value || item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchType = !typeFilter.value || item.type === typeFilter.value
    return matchSearch && matchType
  })
})

const getTypeTag = (type) => {
  const tags = { banner: 'primary', feed: 'success', popup: 'warning' }
  return tags[type] || 'info'
}

const getTypeText = (type) => {
  const texts = { banner: '横幅广告', feed: '信息流', popup: '弹窗广告' }
  return texts[type] || type
}

const handleSearch = () => {}
const handleAdd = () => {}
</script>

<style scoped>
.adunits-container {
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

.adunits-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-6);
}

.adunit-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  padding: var(--spacing-6);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-normal);
}

.adunit-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-4);
}

.card-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-lg);
}

.card-icon.banner {
  background: rgba(6, 182, 212, 0.1);
  color: #06b6d4;
}

.card-icon.feed {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.card-icon.popup {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.type-tag {
  border-radius: 20px;
  padding: 2px 10px;
}

.card-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
  margin: 0 0 var(--spacing-2);
}

.card-desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--spacing-4);
  line-height: 1.6;
}

.card-stats {
  display: flex;
  gap: var(--spacing-6);
  padding: var(--spacing-4);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-4);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
}

.stat-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-top: var(--spacing-1);
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--spacing-4);
  border-top: 1px solid var(--color-border-light);
}

.card-date {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.card-actions {
  display: flex;
  gap: var(--spacing-2);
}

.btn-icon {
  width: 32px;
  height: 32px;
  padding: 0;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}

.btn-icon:hover {
  color: var(--color-accent);
  background: rgba(6, 182, 212, 0.1);
}

.btn-icon.btn-danger:hover {
  color: var(--color-danger);
  background: rgba(239, 68, 68, 0.1);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-16);
  color: var(--color-text-muted);
  grid-column: 1 / -1;
}

.empty-state svg {
  margin-bottom: var(--spacing-4);
  opacity: 0.5;
}

.empty-state p {
  margin: 0 0 var(--spacing-6);
}

@media (max-width: 1200px) {
  .adunits-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .adunits-grid {
    grid-template-columns: 1fr;
  }
  
  .filter-bar {
    flex-direction: column;
  }
  
  .search-input {
    width: 100%;
  }
}
</style>