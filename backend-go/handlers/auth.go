package handlers

import (
	"ad-bi-backend/config"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 查询用户
	db := config.GetDB()
	var userID int
	var storedPassword string
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		req.Username,
	).Scan(&userID, &storedPassword)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	// 简单密码验证（实际应该用 bcrypt 等加密）
	if storedPassword != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成 Token（实际应该用 JWT）
	token := "sample_token_" + time.Now().Format("20060102150405")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"user_id":  userID,
			"username": req.Username,
			"token":    token,
		},
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	db := config.GetDB()
	
	// 插入用户
	result, err := db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		req.Username, req.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注册失败",
		})
		return
	}

	userID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": gin.H{
			"user_id": userID,
		},
	})
}
