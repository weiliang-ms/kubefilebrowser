package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kubecp/config"
	"kubecp/controller"
	"kubecp/controller/k8s"
	"kubecp/controller/static"
	_ "kubecp/docs"
	"kubecp/logs"
	"kubecp/utils"
	"kubecp/utils/denyip"
)

var (
	checker *denyip.Checker
	err     error
)

func init() {
	if !utils.InSliceString("*", config.Config.IPWhiteList) && len(config.Config.IPWhiteList) != 0 {
		checker, err = denyip.NewChecker(config.Config.IPWhiteList)
	}
	if err != nil {
		logs.Fatal(err)
	}
}

func Router() *gin.Engine {
	router := gin.New()
	// 设置文件上传大小限制为8G
	router.MaxMultipartMemory = 32 << 20
	router.Use(logs.Logger(), gin.Recovery(), gzip.Gzip(gzip.DefaultCompression))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFS("/tmp", gin.Dir(utils.RootPath()+"tmp/", true))
	staticGroup := router.Group("/")
	{
		staticGroup.GET("/", static.StatusHtml)
		staticGroup.GET("/upload", static.UploadHtml)
		staticGroup.GET("/multi_upload", static.MultiUploadHtml)
		staticGroup.GET("/download", static.DownloadHtml)
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
		}
	}
	return router
}

func handlersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		render := controller.Gin{C: c}
		reqIPAddr := denyip.GetRemoteIP(c.Request)
		if !utils.InSliceString("*", config.Config.IPWhiteList) && len(config.Config.IPWhiteList) != 0 {
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
