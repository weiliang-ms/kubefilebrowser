package middleware

import (
    "github.com/gin-gonic/gin"
    "kubefilebrowser/apis"
    "kubefilebrowser/configs"
    "kubefilebrowser/utils"
    "kubefilebrowser/utils/denyip"
    "kubefilebrowser/utils/logs"
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

func DenyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        render := apis.Gin{C: c}
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