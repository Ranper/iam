package apiserver

import (
	"os"

	"github.com/gin-gonic/gin"
)

func initRouter(g *gin.Engine) {
	installMiddlewares(g)
	installController(g)
}

func installMiddlewares(g *gin.Engine) {

}

func installController(g *gin.Engine) *gin.Engine {
	g.GET("/hostname", func(c *gin.Context) {
		hostName, _ := os.Hostname()
		c.JSON(200, gin.H{
			"hostname": hostName,
		})
	})

	return g
}
