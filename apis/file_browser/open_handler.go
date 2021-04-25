package file_browser

import (
    "github.com/gin-gonic/gin"
    "kubefilebrowser/apis"
    "kubefilebrowser/utils"
    "kubefilebrowser/utils/logs"
)

// @Summary OpenFile
// @description 容器文件浏览器 - 打开文件
// @Tags FileBrowser
// @Param namespace query FileBrowserQuery true "namespace"
// @Param pods query FileBrowserQuery true "Pod名称"
// @Param container query FileBrowserQuery true "容器名称"
// @Param path query FileBrowserQuery true "路径"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/file_browser/open [get]
func OpenFile(c *gin.Context) {
    render := apis.Gin{C: c}
    var query = &FileBrowserQuery{}
    if err := c.ShouldBindQuery(query); err != nil {
        logs.Error(err)
        render.SetError(utils.CODE_ERR_PARAM, err)
        return
    }
    query.Command = []string{"/kf_tools", "cat", query.Path}
    bs, err := query.fileBrowser()
    if err != nil {
        render.SetError(utils.CODE_ERR_PARAM, err)
        return
    }
    render.SetJson(string(bs))
}
