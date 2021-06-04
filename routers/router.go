package routers

import (
	_ "embed"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kubefilebrowser/apis/file_browser"
	"kubefilebrowser/apis/k8s"
	"kubefilebrowser/configs"
	_ "kubefilebrowser/docs"
	"kubefilebrowser/routers/middleware"
	"kubefilebrowser/utils/logs"
	"net/http/pprof"
)

func Router() *gin.Engine {
	router := gin.New()
	// 设置文件上传大小限制为8G
	router.MaxMultipartMemory = 32 << 20
	// middleware
	router.Use(
		logs.Logger(),
		cors.Default(),
		gin.Recovery(),
		gzip.Gzip(gzip.DefaultCompression),
		middleware.NoCache(),
		middleware.DenyMiddleware(),
		middleware.RequestIDMiddleware(),
		middleware.PromMiddleware(nil),
	)
	// prometheus metrics
	router.GET("/metrics", middleware.PromHandler(promhttp.Handler()))
	// swagger doc
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// api
	apiGroup := router.Group("/api")
	{
		k8sGroup := apiGroup.Group("/k8s")
		{
			//k8sGroup.GET("/ws", k8s.WatchPods)
			k8sGroup.GET("/namespace", k8s.ListNamespace)
			k8sGroup.GET("/deployment", k8s.ListNamespaceAllDeployment)
			k8sGroup.GET("/pods", k8s.ListNamespaceAllPods)
			k8sGroup.GET("/nodes", k8s.ListAllNodes)
			k8sGroup.GET("/status", k8s.PodStatus)
			k8sGroup.POST("/upload", k8s.Copy2Container)
			k8sGroup.POST("/multi_upload", k8s.MultiCopy2Container)
			k8sGroup.GET("/download", k8s.Copy2Local)
			k8sGroup.GET("/terminal", k8s.Terminal)
			k8sGroup.GET("/exec", k8s.Exec)
		}
		FileBrowserGroup := apiGroup.Group("/file_browser")
		{
			FileBrowserGroup.GET("/list", file_browser.ListFile)
			FileBrowserGroup.GET("/open", file_browser.OpenFile)
			FileBrowserGroup.POST("/create_file", file_browser.CreateFile)
			FileBrowserGroup.POST("/create_dir", file_browser.CreateDir)
			FileBrowserGroup.POST("/rename", file_browser.Rename)
			FileBrowserGroup.POST("/remove", file_browser.Remove)
		}
	}
	if configs.Config.RunMode != gin.DebugMode {
		return router
	}
	// debug
	debugGroup := router.Group("/debug/pprof")
	{
		debugGroup.GET("/", gin.WrapF(pprof.Index))
		debugGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debugGroup.GET("/profile", gin.WrapF(pprof.Profile))
		debugGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		debugGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		debugGroup.GET("/trace", gin.WrapF(pprof.Trace))
		debugGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		debugGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		debugGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		debugGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		debugGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		debugGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
	return router
}
