package handler

import (
	"net/http"
	"pet-project/db"
	"pet-project/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthCheck 健康检查端点
func HealthCheck(c *gin.Context) {
	logger.Logger.Info("Health check endpoint called", zap.String("clientIP", c.ClientIP()))

	healthInfo := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// 检查数据库连接
	dbErr := db.DB.Exec("SELECT 1").Error
	if dbErr != nil {
		healthInfo["database"] = "unreachable"
		healthInfo["status"] = "unhealthy"
		logger.Logger.Error("Database health check failed", zap.Error(dbErr))
	} else {
		healthInfo["database"] = "reachable"
	}

	// 检查Redis连接
	redisPing := db.Rdb.Ping(c)
	redisErr := redisPing.Err()
	if redisErr != nil {
		healthInfo["redis"] = "unreachable"
		healthInfo["status"] = "unhealthy"
		logger.Logger.Error("Redis health check failed", zap.Error(redisErr))
	} else {
		healthInfo["redis"] = "reachable"
	}

	// 记录健康检查结果
	logger.Logger.Info("Health check completed", zap.String("status", healthInfo["status"].(string)))

	statusCode := http.StatusOK
	if healthInfo["status"] == "unhealthy" {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, healthInfo)
}