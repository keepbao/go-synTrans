package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keepbao/go-synTrans/config"
	"github.com/keepbao/go-synTrans/server/controller"
	"github.com/keepbao/go-synTrans/server/ws"
)

//go:embed frontend/dist/*
//打包成exe程序
var FS embed.FS

func Run() {
	hub := ws.NewHub()
	go hub.Run()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "<h1>Hello World</h1>")
	})
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.POST("/api/v1/files", controller.FileController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.POST("/api/v1/texts", controller.TextsController)
	router.StaticFS("/static", http.FS(staticFiles))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":" + config.GetPort())
}
