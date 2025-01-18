// utils/logger.go
package utils

import (
	"github.com/gin-gonic/gin"
)

func SetupLogger() {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Customize the default logger if needed
	// For example, integrate with logrus or zap for better logging
}
