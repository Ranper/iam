package middleware

import "github.com/gin-gonic/gin"

// Middlewares 存储注册的中间件
var Middlewares = defaultMiddleware()

func defaultMiddleware() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery": gin.Recovery(),
	}
}
