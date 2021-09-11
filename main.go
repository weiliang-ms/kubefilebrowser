package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"kubefilebrowser/configs"
	_ "kubefilebrowser/configs"
	"kubefilebrowser/routers"
	_ "kubefilebrowser/routers"
	"kubefilebrowser/utils/logs"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	BuildAt string
	GitHash string
)

func init() {
	configs.Init(fmt.Sprintf("Hash: %s\nBuildDate: %s", GitHash, BuildAt))
}

// @title KubeFileBrowser Swagger
// @version 1.0
// @description Kubernetes FileBrowser
// @BasePath /
// @query.collection.format multi
func main() {
	gin.SetMode(configs.Config.RunMode)
	logs.Info("Start up...")
	r := routers.Router()
	r.Use(feMw("/"))
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", configs.Config.HTTPAddr, configs.Config.HTTPPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 300,
		ReadTimeout:  time.Second * 300,
		IdleTimeout:  time.Second * 300,
		Handler:      r, // Pass our instance of gin in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logs.Error("shutting down")
	os.Exit(0)
}

//go:embed static/*
var staticFS embed.FS

const fsBase = "static"

//feMw 使用go.16新的特性embed 到包前端编译后的代码. 替代nginx.   one binary rules them all
func feMw(urlPrefix string) gin.HandlerFunc {
	const indexHtml = "index.html"
	return func(c *gin.Context) {
		urlPath := strings.TrimSpace(c.Request.URL.Path)
		if urlPath == urlPrefix {
			urlPath = path.Join(urlPrefix, indexHtml)
		}
		urlPath = path.Join(fsBase, urlPath)
		f, err := staticFS.Open(urlPath)
		if err != nil {
			//NoRoute
			bs, err := staticFS.ReadFile(path.Join(fsBase, "/", indexHtml))
			if err != nil {
				logs.Error(err, "embed fs")
				return
			}
			c.Status(200)
			c.Writer.Write(bs)
			c.Abort()
			return
		}
		fi, err := f.Stat()
		if strings.HasSuffix(urlPath, ".html") {
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "text/html; charset=utf-8")
		}

		if strings.HasSuffix(urlPath, ".js") {
			c.Header("Content-Type", "text/javascript; charset=utf-8")
		}
		if strings.HasSuffix(urlPath, ".css") {
			c.Header("Content-Type", "text/css; charset=utf-8")
		}

		if err != nil || !fi.IsDir() {
			bs, err := staticFS.ReadFile(urlPath)
			if err != nil {
				logs.Error(err, "embed fs")
				return
			}
			c.Status(200)
			c.Writer.Write(bs)
			c.Abort()
		}
	}
}
