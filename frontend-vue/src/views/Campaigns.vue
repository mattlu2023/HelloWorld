<template>
  <div class="campaigns-container">
    <h2>广告活动管理</h2>
    <p>广告活动列表和 CRUD 操作界面（待实现）</p>
    
    <el-table :data="campaigns" style="width: 100%; margin-top: 20px">
      <el-table-column prop="name" label="活动名称" />
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="budget" label="预算" />
      <el-table-column prop="start_date" label="开始日期" />
      <el-table-column prop="end_date" label="结束日期" />
      <el-table-column label="操作">
        <template #default>
          <el-button size="small">编辑</el-button>
          <el-button size="small" type="danger">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const campaigns = ref([
  { name: '双 11 促销活动', status: 'active', budget: 50000, start_date: '2024-11-01', end_date: '2024-11-11' },
  { name: '新年特惠活动', status: 'active', budget: 30000, start_date: '2024-01-01', end_date: '2024-01-15' }
])
</script>

<style scoped>
.campaigns-container {
  padding: 20px;
}
</style>
