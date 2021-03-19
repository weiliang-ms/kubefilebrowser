package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kubecp/configs"
	"kubecp/controller"
	"kubecp/controller/k8s"
	"kubecp/controller/static"
	_ "kubecp/docs"
	"kubecp/logs"
	"kubecp/utils"
	"kubecp/utils/denyip"
	"net/http"
)

var (
	checker *denyip.Checker
	err     error
)

func init() {
	if !utils.InSliceString("*", configs.Config.IPWhiteList) && len(configs.Config.IPWhiteList) != 0 {
		checker, err = denyip.NewChecker(configs.Config.IPWhiteList)
	}
	if err != nil {
		logs.Fatal(err)
	}
}

func Router() *gin.Engine {
	router := gin.New()
	// 设置文件上传大小限制为8G
	router.MaxMultipartMemory = 32 << 20
	router.Use(logs.Logger(), gin.Recovery(), gzip.Gzip(gzip.DefaultCompression), cors.Default())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// 静态资源
	router.StaticFile("/", "static/index.html")
	router.Static("/static", "static")
	router.LoadHTMLFiles("static/index.html")
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	
	staticGroup := router.Group("/test")
	{
		staticGroup.GET("/test", static.StatusHtml)
		staticGroup.GET("/upload", static.UploadHtml)
		staticGroup.GET("/multi_upload", static.MultiUploadHtml)
		staticGroup.GET("/download", static.DownloadHtml)
		staticGroup.GET("/terminal", static.TerminalHtml)
	}

	apiGroup := router.Group("/api", handlersMiddleware())
	{
		k8sGroup := apiGroup.Group("/k8s")
		{
			//k8sGroup.GET("/ws", k8s.WatchPods)
			k8sGroup.GET("/namespace", k8s.ListNamespace)
			k8sGroup.GET("/deployment", k8s.ListNamespaceAllDeployment)
			k8sGroup.GET("/pods", k8s.ListNamespaceAllPods)
			k8sGroup.GET("/status", k8s.PodStatus)
			k8sGroup.POST("/upload", k8s.Copy2Container)
			k8sGroup.POST("/multi_upload", k8s.MultiCopy2Container)
			k8sGroup.GET("/download", k8s.Copy2Local)
			k8sGroup.GET("/terminal", k8s.Terminal)
			k8sGroup.GET("/exec", k8s.Exec)
			k8sGroup.GET("/file_browser", k8s.FileBrowser)
		}
	}
	return router
}

func handlersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		render := controller.Gin{C: c}
		reqIPAddr := denyip.GetRemoteIP(c.Request)
		if !utils.InSliceString("*", configs.Config.IPWhiteList) && len(configs.Config.IPWhiteList) != 0 {
			reeIPadLenOffset := len(reqIPAddr) - 1
			for i := reeIPadLenOffset; i >= 0; i-- {
				err = checker.IsAuthorized(reqIPAddr[i])
				if err != nil {
					logs.Error(err)
					render.SetError(utils.CODE_ERR_NO_PRIV, err)
					return
				}
			}
		}
		c.Next()
	}
}
