package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubefilebrowser/apis"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
)

type ListNodesQuery struct {
	Name          string `json:"name" form:"name"`
	FieldSelector string `json:"field_selector" form:"field_selector"`
	LabelSelector string `json:"label_selector" form:"label_selector"`
}

// @Summary ListAllNodes
// @description 节点资源列表
// @Tags Kubernetes
// @Param name query ListNodesQuery false "name"
// @Param field_selector query ListNodesQuery false "field_selector"
// @Param label_selector query ListNodesQuery false "label_selector"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/nodes [get]
func ListAllNodes(c *gin.Context) {
	render := apis.Gin{C: c}
	var query ListNodesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}

	if _, ok := c.GetQuery("name"); !ok {
		res, err := configs.RestClient.CoreV1().Nodes().
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: query.LabelSelector,
				FieldSelector: query.FieldSelector,
			})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		render.SetJson(res)
		return
	}
	res, err := configs.RestClient.CoreV1().Nodes().
		Get(context.TODO(), query.Name, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetJson(res)
}
