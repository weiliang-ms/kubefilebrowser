# kubecp

网页版kubectl cp

## debug模式启动日志

```text
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (4 handlers)
[GIN-debug] GET    /tmp/*filepath            --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (4 handlers)
[GIN-debug] HEAD   /tmp/*filepath            --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (4 handlers)
[GIN-debug] GET    /                         --> kubecp/controller/static.StatusHtml (4 handlers)
[GIN-debug] GET    /upload                   --> kubecp/controller/static.UploadHtml (4 handlers)
[GIN-debug] GET    /multi_upload             --> kubecp/controller/static.MultiUploadHtml (4 handlers)
[GIN-debug] GET    /download                 --> kubecp/controller/static.DownloadHtml (4 handlers)
[GIN-debug] GET    /api/k8s/namespace        --> kubecp/controller/k8s.ListNamespace (5 handlers)
[GIN-debug] GET    /api/k8s/deployment       --> kubecp/controller/k8s.ListNamespaceAllDeployment (5 handlers)
[GIN-debug] GET    /api/k8s/pods             --> kubecp/controller/k8s.ListNamespaceAllPods (5 handlers)
[GIN-debug] GET    /api/k8s/status           --> kubecp/controller/k8s.PodStatus (5 handlers)
[GIN-debug] POST   /api/k8s/upload           --> kubecp/controller/k8s.Copy2Container (5 handlers)
[GIN-debug] POST   /api/k8s/multi_upload     --> kubecp/controller/k8s.MultiCopy2Container (5 handlers)
[GIN-debug] GET    /api/k8s/download         --> kubecp/controller/k8s.Copy2Local (5 handlers)
```

## Swagger

![kubecp swagger image](https://raw.githubusercontent.com/xmapst/kubecp/main/swagger.jpg)