require('dotenv').config();
const express = require('express');
const cors = require('cors');
const axios = require('axios');
const rateLimit = require('express-rate-limit');

const app = express();
const PORT = process.env.PORT || 3001;
const GO_BACKEND_URL = process.env.GO_BACKEND_URL || 'http://localhost:8080';

// 中间件配置
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// 请求限流
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 分钟
  max: 100 // 最多 100 个请求
});
app.use('/api/', limiter);

// 代理 Go 后端的 API 请求
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

// 健康检查接口
app.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    message: '中间层服务运行正常',
    timestamp: new Date().toISOString()
  });
});

// 错误处理中间件
app.use((err, req, res, next) => {
  console.error('服务器错误:', err);
  res.status(500).json({
    code: 500,
    message: '服务器内部错误',
    error: process.env.NODE_ENV === 'development' ? err.message : undefined
  });
});

// 404 处理
app.use((req, res) => {
  res.status(404).json({
    code: 404,
    message: '接口不存在'
  });
});

// 启动服务器
app.listen(PORT, () => {
  console.log(`中间层服务启动在 http://localhost:${PORT}`);
  console.log(`Go 后端地址：${GO_BACKEND_URL}`);
});
