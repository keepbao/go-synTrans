package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/keepbao/go-synTrans/server/controller"
	"github.com/keepbao/go-synTrans/server/ws"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run(start chan int, end chan interface{}) {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.POST("/api/v1/files", controller.FileController)
	router.POST("/api/v1/texts", controller.TextsController)
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
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
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	port := 27149
	start <- port
	runErr := router.Run(fmt.Sprintf(":%d", port))
	if runErr != nil {
		end <- runErr
		log.Fatal(runErr)
	}
}
