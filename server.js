const express = require('express');
const sqlite3 = require('sqlite3').verbose();
const cors = require('cors');
const bodyParser = require('body-parser');
const path = require('path');

const app = express();
const PORT = 3000;

app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(express.static('public'));

const db = new sqlite3.Database('./database.sqlite', (err) => {
  if (err) {
    console.error('数据库连接失败:', err.message);
  } else {
    console.log('已连接到 SQLite 数据库');
    db.run(`CREATE TABLE IF NOT EXISTS submissions (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL,
      gender TEXT NOT NULL,
      email TEXT NOT NULL,
      message TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`, (err) => {
      if (err) {
        console.error('创建表失败:', err.message);
      } else {
        console.log('数据库表已准备就绪');
        db.run('ALTER TABLE submissions ADD COLUMN gender TEXT', (err) => {
          if (err && !err.message.includes('duplicate')) {
            console.error('添加 gender 列失败:', err.message);
          }
        });
      }
    });
  }
});

app.post('/api/submit', (req, res) => {
  const { name, gender, email, message } = req.body;
  
  if (!name || !gender || !email || !message) {
    return res.status(400).json({ error: '请填写所有必填字段' });
  }
  
  const sql = 'INSERT INTO submissions (name, gender, email, message) VALUES (?, ?, ?, ?)';
  
  db.run(sql, [name, gender, email, message], function(err) {
    if (err) {
      console.error('插入数据失败:', err.message);
      return res.status(500).json({ error: '提交失败，请稍后重试' });
    }
    
    res.json({
      success: true,
      message: '提交成功',
      id: this.lastID
    });
  });
});

app.get('/api/submissions', (req, res) => {
  const sql = 'SELECT * FROM submissions ORDER BY created_at DESC';
  
  db.all(sql, [], (err, rows) => {
    if (err) {
      console.error('查询数据失败:', err.message);
      return res.status(500).json({ error: '查询失败，请稍后重试' });
    }
    
    res.json({
      success: true,
      data: rows
    });
  });
});

app.get('/api/submissions/:id', (req, res) => {
  const { id } = req.params;
  const sql = 'SELECT * FROM submissions WHERE id = ?';
  
  db.get(sql, [id], (err, row) => {
    if (err) {
      console.error('查询数据失败:', err.message);
      return res.status(500).json({ error: '查询失败，请稍后重试' });
    }
    
    if (!row) {
      return res.status(404).json({ error: '未找到该提交记录' });
    }
    
    res.json({
      success: true,
      data: row
    });
  });
});

app.listen(PORT, () => {
  console.log(`服务器运行在 http://localhost:${PORT}`);
});
