-- 广告分析 BI 系统数据库初始化脚本
-- MySQL 8.0+

-- 创建数据库
CREATE DATABASE IF NOT EXISTS ad_bi_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ad_bi_system;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    role ENUM('admin', 'user') DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 广告活动表
CREATE TABLE IF NOT EXISTS ad_campaigns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '活动名称',
    description TEXT COMMENT '活动描述',
    status ENUM('active', 'paused', 'ended') DEFAULT 'active' COMMENT '活动状态',
    budget DECIMAL(10, 2) DEFAULT 0.00 COMMENT '预算',
    start_date DATE COMMENT '开始日期',
    end_date DATE COMMENT '结束日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_dates (start_date, end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 广告单元表
CREATE TABLE IF NOT EXISTS ad_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    campaign_id INT NOT NULL COMMENT '所属活动 ID',
    name VARCHAR(200) NOT NULL COMMENT '广告单元名称',
    ad_type ENUM('banner', 'video', 'native', 'interstitial') DEFAULT 'banner' COMMENT '广告类型',
    placement VARCHAR(100) COMMENT '投放位置',
    creative_url VARCHAR(500) COMMENT '创意素材 URL',
    landing_url VARCHAR(500) COMMENT '落地页 URL',
    status ENUM('active', 'paused') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES ad_campaigns(id) ON DELETE CASCADE,
    INDEX idx_campaign (campaign_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 广告数据统计表（按天聚合）
CREATE TABLE IF NOT EXISTS ad_stats_daily (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    stat_date DATE NOT NULL COMMENT '统计日期',
    campaign_id INT NOT NULL COMMENT '活动 ID',
    ad_unit_id INT NOT NULL COMMENT '广告单元 ID',
    impressions BIGINT DEFAULT 0 COMMENT '曝光量',
    clicks BIGINT DEFAULT 0 COMMENT '点击量',
    conversions BIGINT DEFAULT 0 COMMENT '转化量',
    cost DECIMAL(10, 2) DEFAULT 0.00 COMMENT '消耗',
    revenue DECIMAL(10, 2) DEFAULT 0.00 COMMENT '收入',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES ad_campaigns(id) ON DELETE CASCADE,
    FOREIGN KEY (ad_unit_id) REFERENCES ad_units(id) ON DELETE CASCADE,
    UNIQUE KEY uk_date_campaign_unit (stat_date, campaign_id, ad_unit_id),
    INDEX idx_date (stat_date),
    INDEX idx_campaign (campaign_id),
    INDEX idx_ad_unit (ad_unit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户行为表
CREATE TABLE IF NOT EXISTS user_actions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(100) COMMENT '用户 ID',
    session_id VARCHAR(100) COMMENT '会话 ID',
    action_type ENUM('view', 'click', 'conversion', 'bounce') NOT NULL COMMENT '行为类型',
    ad_unit_id INT COMMENT '广告单元 ID',
    campaign_id INT COMMENT '活动 ID',
    page_url VARCHAR(500) COMMENT '页面 URL',
    referrer VARCHAR(500) COMMENT '来源 URL',
    device_type ENUM('desktop', 'mobile', 'tablet') DEFAULT 'desktop' COMMENT '设备类型',
    os VARCHAR(50) COMMENT '操作系统',
    browser VARCHAR(50) COMMENT '浏览器',
    ip_address VARCHAR(45) COMMENT 'IP 地址',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES ad_campaigns(id) ON DELETE SET NULL,
    FOREIGN KEY (ad_unit_id) REFERENCES ad_units(id) ON DELETE SET NULL,
    INDEX idx_user (user_id),
    INDEX idx_session (session_id),
    INDEX idx_action_type (action_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 转化漏斗表
CREATE TABLE IF NOT EXISTS conversion_funnel (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    funnel_date DATE NOT NULL COMMENT '日期',
    campaign_id INT NOT NULL COMMENT '活动 ID',
    step_name VARCHAR(50) NOT NULL COMMENT '漏斗步骤',
    step_order INT NOT NULL COMMENT '步骤顺序',
    user_count INT DEFAULT 0 COMMENT '用户数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES ad_campaigns(id) ON DELETE CASCADE,
    UNIQUE KEY uk_date_campaign_step (funnel_date, campaign_id, step_name),
    INDEX idx_date (funnel_date),
    INDEX idx_campaign (campaign_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入测试数据
-- 测试用户
INSERT INTO users (username, password, email, role) VALUES
('admin', 'admin123', 'admin@example.com', 'admin'),
('user1', 'user123', 'user1@example.com', 'user');

-- 测试广告活动
INSERT INTO ad_campaigns (name, description, status, budget, start_date, end_date) VALUES
('双 11 促销活动', '双 11 期间的大型促销活动', 'active', 50000.00, '2024-11-01', '2024-11-11'),
('新年特惠活动', '新年特惠，全场打折', 'active', 30000.00, '2024-01-01', '2024-01-15'),
('春季新品发布', '春季新品上市推广', 'paused', 20000.00, '2024-03-01', '2024-03-31');

-- 测试广告单元
INSERT INTO ad_units (campaign_id, name, ad_type, placement, status) VALUES
(1, '双 11 首页 Banner', 'banner', 'homepage_top', 'active'),
(1, '双 11 视频广告', 'video', 'video_pre_roll', 'active'),
(2, '新年横幅广告', 'banner', 'sidebar', 'active'),
(3, '春季开屏广告', 'interstitial', 'app_splash', 'paused');

-- 测试统计数据（最近 7 天）
INSERT INTO ad_stats_daily (stat_date, campaign_id, ad_unit_id, impressions, clicks, conversions, cost, revenue) VALUES
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 1, 1, 10000, 500, 50, 1000.00, 5000.00),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 1, 1, 12000, 600, 60, 1200.00, 6000.00),
(DATE_SUB(CURDATE(), INTERVAL 4 DAY), 1, 1, 11000, 550, 55, 1100.00, 5500.00),
(DATE_SUB(CURDATE(), INTERVAL 3 DAY), 1, 1, 13000, 650, 65, 1300.00, 6500.00),
(DATE_SUB(CURDATE(), INTERVAL 2 DAY), 1, 1, 15000, 750, 75, 1500.00, 7500.00),
(DATE_SUB(CURDATE(), INTERVAL 1 DAY), 1, 1, 14000, 700, 70, 1400.00, 7000.00),
(CURDATE(), 1, 1, 8000, 400, 40, 800.00, 4000.00),

(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 1, 2, 5000, 300, 30, 800.00, 3000.00),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 1, 2, 6000, 360, 36, 960.00, 3600.00),
(DATE_SUB(CURDATE(), INTERVAL 4 DAY), 1, 2, 5500, 330, 33, 880.00, 3300.00),
(DATE_SUB(CURDATE(), INTERVAL 3 DAY), 1, 2, 6500, 390, 39, 1040.00, 3900.00),
(DATE_SUB(CURDATE(), INTERVAL 2 DAY), 1, 2, 7000, 420, 42, 1120.00, 4200.00),
(DATE_SUB(CURDATE(), INTERVAL 1 DAY), 1, 2, 6800, 408, 41, 1088.00, 4080.00),
(CURDATE(), 1, 2, 4000, 240, 24, 640.00, 2400.00),

(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 2, 3, 8000, 400, 40, 600.00, 4000.00),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 2, 3, 9000, 450, 45, 675.00, 4500.00),
(DATE_SUB(CURDATE(), INTERVAL 4 DAY), 2, 3, 8500, 425, 42, 637.50, 4250.00),
(DATE_SUB(CURDATE(), INTERVAL 3 DAY), 2, 3, 9500, 475, 47, 712.50, 4750.00),
(DATE_SUB(CURDATE(), INTERVAL 2 DAY), 2, 3, 10000, 500, 50, 750.00, 5000.00),
(DATE_SUB(CURDATE(), INTERVAL 1 DAY), 2, 3, 9800, 490, 49, 735.00, 4900.00),
(CURDATE(), 2, 3, 5000, 250, 25, 375.00, 2500.00);

-- 测试转化漏斗数据
INSERT INTO conversion_funnel (funnel_date, campaign_id, step_name, step_order, user_count) VALUES
(CURDATE(), 1, '曝光', 1, 10000),
(CURDATE(), 1, '点击', 2, 500),
(CURDATE(), 1, '访问落地页', 3, 400),
(CURDATE(), 1, '加入购物车', 4, 100),
(CURDATE(), 1, '完成购买', 5, 50);

-- 创建视图：广告活动总览
CREATE OR REPLACE VIEW v_campaign_overview AS
SELECT 
    c.id,
    c.name,
    c.status,
    c.budget,
    c.start_date,
    c.end_date,
    COALESCE(SUM(s.impressions), 0) as total_impressions,
    COALESCE(SUM(s.clicks), 0) as total_clicks,
    COALESCE(SUM(s.conversions), 0) as total_conversions,
    COALESCE(SUM(s.cost), 0) as total_cost,
    COALESCE(SUM(s.revenue), 0) as total_revenue,
    CASE 
        WHEN SUM(s.impressions) > 0 THEN ROUND(SUM(s.clicks) / SUM(s.impressions) * 100, 2)
        ELSE 0 
    END as ctr,
    CASE 
        WHEN SUM(s.clicks) > 0 THEN ROUND(SUM(s.conversions) / SUM(s.clicks) * 100, 2)
        ELSE 0 
    END as conversion_rate
FROM ad_campaigns c
LEFT JOIN ad_stats_daily s ON c.id = s.campaign_id
GROUP BY c.id, c.name, c.status, c.budget, c.start_date, c.end_date;

-- 创建视图：每日趋势
CREATE OR REPLACE VIEW v_daily_trend AS
SELECT 
    stat_date,
    SUM(impressions) as impressions,
    SUM(clicks) as clicks,
    SUM(conversions) as conversions,
    SUM(cost) as cost,
    SUM(revenue) as revenue,
    ROUND(SUM(clicks) / NULLIF(SUM(impressions), 0) * 100, 2) as ctr
FROM ad_stats_daily
GROUP BY stat_date
ORDER BY stat_date;

SELECT '数据库初始化完成！' as message;
