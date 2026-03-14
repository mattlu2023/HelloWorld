import request from './request'

// 登录
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 注册
export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

// 获取数据概览
export function getOverview() {
  return request({
    url: '/stats/overview',
    method: 'get'
  })
}

// 获取每日趋势
export function getDailyTrend() {
  return request({
    url: '/stats/daily-trend',
    method: 'get'
  })
}

// 获取转化漏斗
export function getConversionFunnel() {
  return request({
    url: '/stats/funnel',
    method: 'get'
  })
}

// 广告活动管理
export function getCampaigns() {
  return request({
    url: '/campaigns',
    method: 'get'
  })
}

export function getCampaignById(id) {
  return request({
    url: `/campaigns/${id}`,
    method: 'get'
  })
}

export function createCampaign(data) {
  return request({
    url: '/campaigns',
    method: 'post',
    data
  })
}

export function updateCampaign(id, data) {
  return request({
    url: `/campaigns/${id}`,
    method: 'put',
    data
  })
}

export function deleteCampaign(id) {
  return request({
    url: `/campaigns/${id}`,
    method: 'delete'
  })
}

// 广告单元管理
export function getAdUnits() {
  return request({
    url: '/ad-units',
    method: 'get'
  })
}

export function getAdUnitById(id) {
  return request({
    url: `/ad-units/${id}`,
    method: 'get'
  })
}

export function createAdUnit(data) {
  return request({
    url: '/ad-units',
    method: 'post',
    data
  })
}

export function updateAdUnit(id, data) {
  return request({
    url: `/ad-units/${id}`,
    method: 'put',
    data
  })
}

export function deleteAdUnit(id) {
  return request({
    url: `/ad-units/${id}`,
    method: 'delete'
  })
}
