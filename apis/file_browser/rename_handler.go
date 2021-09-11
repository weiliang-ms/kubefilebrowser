package file_browser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kubefilebrowser/apis"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
)


// Rename
// @Summary Rename
// @description 容器文件浏览器 - 重命名
// @Tags FileBrowser
// @Param namespace query FileBrowserQuery true "namespace"
// @Param pods query FileBrowserQuery true "Pod名称"
// @Param container query FileBrowserQuery true "容器名称"
// @Param path query FileBrowserQuery true "新路径"
// @Param old_path query FileBrowserQuery true "旧路径"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/file_browser/rename [post]
func Rename(c *gin.Context) {
	render := apis.Gin{C: c}
	var query = &FileBrowserQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	if query.OldPath == "" {
		render.SetError(utils.CODE_ERR_PARAM, fmt.Errorf("file path does not exist"))
		return
	}
	query.Command = []string{"/kf_tools", "mv", query.OldPath, query.Path}
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
