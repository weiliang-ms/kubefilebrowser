# kubefileborwser

kubernetes container file browser

## 启动可选环境变量

| 名称 | 类型 | 默认值 |
| ---- | ---- | ---- |
| RUN_MODE | string | default:"debug" |
| HTTP_PORT | string | default:"9999" |
| HTTP_ADDR | string | default:"0.0.0.0" |
| LOG_LEVEL | string | default:"debug" |
| IP_WHITE_LIST | []string | default:"*" |
| DOWNLOAD_TMP | string | default:"tmp" |
| KUBECONFIG | string | default:"~/.kube/config" |


## debug模式启动日志

```text
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (5 handlers)
[GIN-debug] GET    /                         --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (5 handlers)
[GIN-debug] HEAD   /                         --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (5 handlers)
[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (5 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (5 handlers)
[GIN-debug] GET    /api/k8s/namespace        --> kubefilebrowser/apis/k8s.ListNamespace (6 handlers)
[GIN-debug] GET    /api/k8s/deployment       --> kubefilebrowser/apis/k8s.ListNamespaceAllDeployment (6 handlers)
[GIN-debug] GET    /api/k8s/pods             --> kubefilebrowser/apis/k8s.ListNamespaceAllPods (6 handlers)
[GIN-debug] GET    /api/k8s/status           --> kubefilebrowser/apis/k8s.PodStatus (6 handlers)
[GIN-debug] POST   /api/k8s/upload           --> kubefilebrowser/apis/k8s.Copy2Container (6 handlers)
[GIN-debug] POST   /api/k8s/multi_upload     --> kubefilebrowser/apis/k8s.MultiCopy2Container (6 handlers)
[GIN-debug] GET    /api/k8s/download         --> kubefilebrowser/apis/k8s.Copy2Local (6 handlers)
[GIN-debug] GET    /api/k8s/terminal         --> kubefilebrowser/apis/k8s.Terminal (6 handlers)
[GIN-debug] GET    /api/k8s/exec             --> kubefilebrowser/apis/k8s.Exec (6 handlers)
[GIN-debug] GET    /api/k8s/file_browser     --> kubefilebrowser/apis/k8s.FileBrowser (6 handlers)
```

## Index.html
![kubefilebrowser_index_html](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/index_html.jpg)

## Swagger

![kubefilebrowser swagger image](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/swagger.jpg)
