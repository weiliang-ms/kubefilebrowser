package file_browser

import (
	"github.com/gin-gonic/gin"
	"kubefilebrowser/apis"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
	"strings"
)

// Remove
// @Summary Remove
// @description 容器文件浏览器 - 删除
// @Tags FileBrowser
// @Param namespace query FileBrowserQuery true "namespace"
// @Param pods query FileBrowserQuery true "Pod名称"
// @Param container query FileBrowserQuery true "容器名称"
// @Param path query FileBrowserQuery true "路径"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/file_browser/remove [post]
func Remove(c *gin.Context) {
	render := apis.Gin{C: c}
	var query = &FileBrowserQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	query.Command = append([]string{"/kf_tools", "rm"}, strings.Split(query.Path, ",")...)
	bs, err := query.FileBrowser()
	if err != nil {
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	if len(string(bs)) != 0  {
		render.SetJson(string(bs))
		return
	}
	render.SetJson("success")
}
