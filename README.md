# kubefileborwser

kubernetes container file browser

## 启动可选环境变量

| 名称 | 类型 | 默认值 |
| ---- | ---- | ---- |
| RUN_MODE | string | debug |
| HTTP_PORT | string | 9999 |
| HTTP_ADDR | string | 0.0.0.0 |
| LOG_LEVEL | string | debug |
| IP_WHITE_LIST | []string | * |
| KUBECONFIG | string | ~/.kube/config |

+ 部署在k8s内创建使用管理员权限的serviceaccount即可

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

## file_browser
![kubefilebrowser](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/file_browser.jpg)

## terminal
![terminal](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/terminal.jpg)

## Swagger

![kubefilebrowser swagger image](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/swagger.jpg)

## Reference documents

+ [golang 1.16 gin static embed](https://mojotv.cn/golang/golang-html5-websocket-remote-desktop)
+ [vue](https://cli.vuejs.org/config/)
+ [kubectl copy & shell 原理讲解](https://www.yfdou.com/archives/kuberneteszhi-kubectlexeczhi-ling-gong-zuo-yuan-li-shi-xian-copyhe-webshellyi-ji-filebrowser.html)
