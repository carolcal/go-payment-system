package routes

import (

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func SetUpWebRoutes(router *gin.Engine) {

	router.LoadHTMLGlob("web/*.html")
	router.Static("css", "./web/css")
	router.Static("js", "./web/js")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

}