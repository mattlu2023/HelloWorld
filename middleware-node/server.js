require('dotenv').config();
const express = require('express');
const cors = require('cors');
const axios = require('axios');
const rateLimit = require('express-rate-limit');

const app = express();
const PORT = process.env.PORT || 3002;
const GO_BACKEND_URL = process.env.GO_BACKEND_URL || 'http://localhost:8080';

app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 100
});
app.use('/api/', limiter);

const generateToken = () => {
  return 'mock-token-' + Math.random().toString(36).substring(2, 15);
};

let campaignsData = [
  {
    id: 1,
    name: '春季促销活动',
    description: '春季新品推广活动，覆盖全渠道',
    status: 'active',
    start_date: '2024-03-01',
    end_date: '2024-03-31',
    budget: 500000,
    spent: 325678.50,
    impressions: 3567890,
    clicks: 245678,
    conversions: 12345,
    revenue: 789012.34,
    created_at: '2024-02-28 10:00:00',
    updated_at: '2024-03-15 14:30:00'
  },
  {
    id: 2,
    name: '品牌曝光计划',
    description: '提升品牌知名度，增加市场份额',
    status: 'active',
    start_date: '2024-02-15',
    end_date: '2024-04-15',
    budget: 800000,
    spent: 567890.00,
    impressions: 6789012,
    clicks: 345678,
    conversions: 8901,
    revenue: 567890.12,
    created_at: '2024-02-10 09:00:00',
    updated_at: '2024-03-14 16:45:00'
  },
  {
    id: 3,
    name: '会员专享活动',
    description: '针对会员用户的专属优惠活动',
    status: 'paused',
    start_date: '2024-03-05',
    end_date: '2024-03-25',
    budget: 200000,
    spent: 156789.00,
    impressions: 1234567,
    clicks: 89012,
    conversions: 4567,
    revenue: 234567.89,
    created_at: '2024-03-01 11:00:00',
    updated_at: '2024-03-10 10:00:00'
  },
  {
    id: 4,
    name: '新品上市推广',
    description: '夏季新品发布推广',
    status: 'ended',
    start_date: '2024-01-01',
    end_date: '2024-02-28',
    budget: 600000,
    spent: 600000.00,
    impressions: 8901234,
    clicks: 567890,
    conversions: 23456,
    revenue: 1234567.89,
    created_at: '2023-12-28 08:00:00',
    updated_at: '2024-02-28 23:59:00'
  },
  {
    id: 5,
    name: '节日营销活动',
    description: '春节期间的促销活动',
    status: 'ended',
    start_date: '2024-01-20',
    end_date: '2024-02-10',
    budget: 400000,
    spent: 389012.34,
    impressions: 4567890,
    clicks: 289012,
    conversions: 15678,
    revenue: 890123.45,
    created_at: '2024-01-15 09:00:00',
    updated_at: '2024-02-10 23:59:00'
  },
  {
    id: 6,
    name: '电商大促活动',
    description: '618大促预热活动',
    status: 'active',
    start_date: '2024-03-10',
    end_date: '2024-06-20',
    budget: 1000000,
    spent: 234567.89,
    impressions: 2345678,
    clicks: 156789,
    conversions: 6789,
    revenue: 456789.01,
    created_at: '2024-03-08 10:00:00',
    updated_at: '2024-03-15 15:00:00'
  }
];

let adUnitsData = [
  {
    id: 1,
    name: '首页横幅广告',
    type: 'banner',
    position: '首页顶部',
    size: '728x90',
    status: 'active',
    impressions: 4567890,
    clicks: 345678,
    conversions: 8901,
    revenue: 234567.89,
    ecpm: 51.35,
    ctr: 7.57,
    created_at: '2024-01-01 08:00:00',
    updated_at: '2024-03-15 10:00:00'
  },
  {
    id: 2,
    name: '详情页信息流',
    type: 'feed',
    position: '商品详情页',
    size: '300x250',
    status: 'active',
    impressions: 3256789,
    clicks: 234567,
    conversions: 12345,
    revenue: 345678.90,
    ecpm: 106.14,
    ctr: 7.20,
    created_at: '2024-01-05 09:00:00',
    updated_at: '2024-03-14 14:30:00'
  },
  {
    id: 3,
    name: '侧边栏推荐',
    type: 'sidebar',
    position: '网站侧边栏',
    size: '160x600',
    status: 'active',
    impressions: 1890123,
    clicks: 89012,
    conversions: 4567,
    revenue: 98765.43,
    ecpm: 52.25,
    ctr: 4.71,
    created_at: '2024-01-10 10:00:00',
    updated_at: '2024-03-12 11:00:00'
  },
  {
    id: 4,
    name: '弹窗广告',
    type: 'popup',
    position: '页面弹窗',
    size: '600x400',
    status: 'paused',
    impressions: 890123,
    clicks: 56789,
    conversions: 3456,
    revenue: 67890.12,
    ecpm: 76.27,
    ctr: 6.38,
    created_at: '2024-02-01 08:00:00',
    updated_at: '2024-03-10 09:00:00'
  },
  {
    id: 5,
    name: '视频前贴片',
    type: 'video',
    position: '视频播放前',
    size: '1280x720',
    status: 'active',
    impressions: 2345678,
    clicks: 167890,
    conversions: 7890,
    revenue: 456789.01,
    ecpm: 194.74,
    ctr: 7.16,
    created_at: '2024-02-15 11:00:00',
    updated_at: '2024-03-15 16:00:00'
  },
  {
    id: 6,
    name: '搜索结果广告',
    type: 'search',
    position: '搜索结果页',
    size: '450x120',
    status: 'active',
    impressions: 5678901,
    clicks: 456789,
    conversions: 18901,
    revenue: 567890.12,
    ecpm: 99.99,
    ctr: 8.04,
    created_at: '2024-01-01 08:00:00',
    updated_at: '2024-03-15 17:00:00'
  }
];

app.post('/api/v1/login', (req, res) => {
  const { username, password } = req.body;
  if ((username === 'admin' && password === 'admin123') ||
      (username === 'user' && password === 'user123')) {
    res.json({
      code: 0,
      message: '登录成功',
      data: {
        token: generateToken(),
        user: {
          id: 1,
          username: username,
          email: username + '@example.com',
          role: username === 'admin' ? 'admin' : 'user',
          avatar: ''
        }
      }
    });
  } else {
    res.json({ code: 401, message: '用户名或密码错误' });
  }
});

app.post('/api/v1/register', (req, res) => {
  const { username, password, email } = req.body;
  res.json({
    code: 0,
    message: '注册成功',
    data: {
      token: generateToken(),
      user: {
        id: Date.now(),
        username,
        email,
        role: 'user',
        avatar: ''
      }
    }
  });
});

app.get('/api/v1/stats/overview', (req, res) => {
  const total_impressions = campaignsData.reduce((sum, c) => sum + (c.impressions || 0), 0);
  const total_clicks = campaignsData.reduce((sum, c) => sum + (c.clicks || 0), 0);
  const total_conversions = campaignsData.reduce((sum, c) => sum + (c.conversions || 0), 0);
  const total_revenue = campaignsData.reduce((sum, c) => sum + (c.revenue || 0), 0);
  
  res.json({
    code: 0,
    message: 'success',
    data: {
      total_impressions,
      total_clicks,
      total_conversions,
      total_revenue,
      avg_ctr: total_impressions > 0 ? ((total_clicks / total_impressions) * 100).toFixed(2) : 0,
      avg_cvr: total_clicks > 0 ? ((total_conversions / total_clicks) * 100).toFixed(2) : 0,
      avg_cpm: total_impressions > 0 ? ((total_revenue / total_impressions) * 1000).toFixed(2) : 0,
      avg_cpc: total_clicks > 0 ? (total_revenue / total_clicks).toFixed(2) : 0
    }
  });
});

app.get('/api/v1/stats/daily-trend', (req, res) => {
  const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日'];
  res.json({
    code: 0,
    message: 'success',
    data: {
      dates: days,
      impressions: [1856234, 2134567, 1987654, 2345678, 2198765, 1678901, 1484543],
      clicks: [125678, 145678, 134567, 156789, 148901, 112345, 68487],
      conversions: [6789, 7890, 7234, 8345, 7987, 5678, 2755],
      revenue: [423456.78, 487654.32, 456789.01, 534567.89, 501234.56, 389012.34, 203063.50]
    }
  });
});

app.get('/api/v1/stats/funnel', (req, res) => {
  const total_impressions = campaignsData.reduce((sum, c) => sum + (c.impressions || 0), 0);
  const total_clicks = campaignsData.reduce((sum, c) => sum + (c.clicks || 0), 0);
  const total_conversions = campaignsData.reduce((sum, c) => sum + (c.conversions || 0), 0);
  
  res.json({
    code: 0,
    message: 'success',
    data: [
      { name: '曝光', value: total_impressions, rate: 100 },
      { name: '点击', value: total_clicks, rate: total_impressions > 0 ? ((total_clicks / total_impressions) * 100).toFixed(2) : 0 },
      { name: '访问', value: Math.floor(total_clicks * 0.88), rate: total_impressions > 0 ? ((total_clicks * 0.88 / total_impressions) * 100).toFixed(2) : 0 },
      { name: '注册', value: Math.floor(total_clicks * 0.14), rate: total_impressions > 0 ? ((total_clicks * 0.14 / total_impressions) * 100).toFixed(2) : 0 },
      { name: '转化', value: total_conversions, rate: total_impressions > 0 ? ((total_conversions / total_impressions) * 100).toFixed(2) : 0 }
    ]
  });
});

app.get('/api/v1/campaigns', (req, res) => {
  res.json({ code: 0, message: 'success', data: campaignsData });
});

app.get('/api/v1/campaigns/:id', (req, res) => {
  const campaign = campaignsData.find(c => c.id === parseInt(req.params.id));
  if (campaign) {
    res.json({ code: 0, message: 'success', data: campaign });
  } else {
    res.json({ code: 404, message: '活动不存在' });
  }
});

app.post('/api/v1/campaigns', (req, res) => {
  const newCampaign = {
    id: Date.now(),
    ...req.body,
    spent: 0,
    impressions: 0,
    clicks: 0,
    conversions: 0,
    revenue: 0,
    created_at: new Date().toLocaleString(),
    updated_at: new Date().toLocaleString()
  };
  campaignsData.push(newCampaign);
  res.json({ code: 0, message: '创建成功', data: newCampaign });
});

app.put('/api/v1/campaigns/:id', (req, res) => {
  const index = campaignsData.findIndex(c => c.id === parseInt(req.params.id));
  if (index !== -1) {
    campaignsData[index] = {
      ...campaignsData[index],
      ...req.body,
      updated_at: new Date().toLocaleString()
    };
    res.json({ code: 0, message: '更新成功', data: campaignsData[index] });
  } else {
    res.json({ code: 404, message: '活动不存在' });
  }
});

app.delete('/api/v1/campaigns/:id', (req, res) => {
  const index = campaignsData.findIndex(c => c.id === parseInt(req.params.id));
  if (index !== -1) {
    campaignsData.splice(index, 1);
    res.json({ code: 0, message: '删除成功' });
  } else {
    res.json({ code: 404, message: '活动不存在' });
  }
});

app.get('/api/v1/ad-units', (req, res) => {
  res.json({ code: 0, message: 'success', data: adUnitsData });
});

app.get('/api/v1/ad-units/:id', (req, res) => {
  const adUnit = adUnitsData.find(a => a.id === parseInt(req.params.id));
  if (adUnit) {
    res.json({ code: 0, message: 'success', data: adUnit });
  } else {
    res.json({ code: 404, message: '广告单元不存在' });
  }
});

app.post('/api/v1/ad-units', (req, res) => {
  const newAdUnit = {
    id: Date.now(),
    ...req.body,
    impressions: 0,
    clicks: 0,
    conversions: 0,
    revenue: 0,
    ecpm: 0,
    ctr: 0,
    created_at: new Date().toLocaleString(),
    updated_at: new Date().toLocaleString()
  };
  adUnitsData.push(newAdUnit);
  res.json({ code: 0, message: '创建成功', data: newAdUnit });
});

app.put('/api/v1/ad-units/:id', (req, res) => {
  const index = adUnitsData.findIndex(a => a.id === parseInt(req.params.id));
  if (index !== -1) {
    adUnitsData[index] = {
      ...adUnitsData[index],
      ...req.body,
      updated_at: new Date().toLocaleString()
    };
    res.json({ code: 0, message: '更新成功', data: adUnitsData[index] });
  } else {
    res.json({ code: 404, message: '广告单元不存在' });
  }
});

app.delete('/api/v1/ad-units/:id', (req, res) => {
  const index = adUnitsData.findIndex(a => a.id === parseInt(req.params.id));
  if (index !== -1) {
    adUnitsData.splice(index, 1);
    res.json({ code: 0, message: '删除成功' });
  } else {
    res.json({ code: 404, message: '广告单元不存在' });
  }
});

app.get('/api/v1/analysis', (req, res) => {
  res.json({
    code: 0,
    message: 'success',
    data: {
      source_distribution: [
        { name: '搜索引擎', value: 45 },
        { name: '社交媒体', value: 28 },
        { name: '直接访问', value: 15 },
        { name: '广告投放', value: 8 },
        { name: '其他', value: 4 }
      ],
      device_distribution: [
        { name: '移动端', value: 62 },
        { name: 'PC端', value: 35 },
        { name: '平板', value: 3 }
      ],
      active_trend: {
        dates: ['03-09', '03-10', '03-11', '03-12', '03-13', '03-14', '03-15'],
        pv: [125678, 134567, 145678, 138901, 156789, 142345, 167890],
        uv: [45678, 48901, 52345, 49876, 56789, 51234, 58901]
      },
      behavior_ranking: [
        { name: '商品浏览', count: 456789, percentage: 35.2 },
        { name: '加入购物车', count: 234567, percentage: 18.1 },
        { name: '下单购买', count: 189012, percentage: 14.5 },
        { name: '查看详情', count: 167890, percentage: 13.0 },
        { name: '收藏商品', count: 123456, percentage: 9.5 },
        { name: '分享转发', count: 67890, percentage: 5.3 },
        { name: '评论互动', count: 56789, percentage: 4.4 }
      ]
    }
  });
});

app.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    message: '中间层服务运行正常',
    timestamp: new Date().toISOString()
  });
});

app.use('/api/v1', async (req, res) => {
  try {
    const url = `${GO_BACKEND_URL}/api/v1${req.url}`;
    const config = {
      method: req.method,
      url: url,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': req.headers.authorization
      }
    };

    if (['POST', 'PUT'].includes(req.method)) {
      config.data = req.body;
    }

    const response = await axios(config);
    res.json(response.data);
  } catch (error) {
    console.error('代理请求错误:', error.message);
    
    if (error.response) {
      res.status(error.response.status).json(error.response.data);
    } else {
      res.status(500).json({
        code: 500,
        message: '后端服务不可用',
        error: error.message
      });
    }
  }
});

app.use((err, req, res, next) => {
  console.error('服务器错误:', err);
  res.status(500).json({
    code: 500,
    message: '服务器内部错误',
    error: process.env.NODE_ENV === 'development' ? err.message : undefined
  });
});

app.use((req, res) => {
  res.status(404).json({
    code: 404,
    message: '接口不存在'
  });
});

app.listen(PORT, () => {
  console.log(`中间层服务启动在 http://localhost:${PORT}`);
  console.log(`Go 后端地址：${GO_BACKEND_URL}`);
});
