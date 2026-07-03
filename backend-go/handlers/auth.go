package handlers

import (
	"ad-bi-backend/config"
	"ad-bi-backend/utils"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
	List       interface{} `json:"list"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误：用户名和密码均为必填项",
		})
		return
	}

	db := config.GetDB()
	var userID int64
	var storedPassword string
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		req.Username,
	).Scan(&userID, &storedPassword)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
		return
	}
	if err != nil {
		log.Printf("查询用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询用户信息失败，请稍后重试",
		})
		return
	}

	if !utils.VerifyPassword(storedPassword, req.Password) {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
		return
	}

	token, err := utils.GenerateToken(userID, req.Username)
	if err != nil {
		log.Printf("生成Token失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "生成认证令牌失败，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "登录成功",
		Data: gin.H{
			"user_id":  userID,
			"username": req.Username,
			"token":    token,
		},
	})
}

func Register(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误：用户名和密码均为必填项",
		})
		return
	}

	db := config.GetDB()

	var existingID int64
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, APIResponse{
			Code:    409,
			Message: "用户名已存在，请选择其他用户名",
		})
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("查询用户名是否存在失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询用户名失败，请稍后重试",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "密码加密失败，请稍后重试",
		})
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		req.Username, hashedPassword,
	)
	if err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "创建用户失败，请稍后重试",
		})
		return
	}

	userID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "注册成功",
		Data: gin.H{
			"user_id": userID,
		},
	})
}