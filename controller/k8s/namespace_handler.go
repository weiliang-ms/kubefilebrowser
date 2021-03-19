package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubecp/configs"
	"kubecp/controller"
	"kubecp/logs"
	"kubecp/utils"
)

type ListNSQuery struct {
	Namespace     string `json:"namespace" form:"namespace"`
	FieldSelector string `json:"field_selector" form:"field_selector"`
	LabelSelector string `json:"label_selector" form:"label_selector"`
}

// @Summary ListNamespace
// @description 命名空间列表
// @Tags Kubernetes
// @Param namespace query ListNSQuery false "namespace"
// @Param field_selector query ListNSQuery false "field_selector"
// @Param label_selector query ListNSQuery false "label_selector"
// @Success 200 {object} controller.JSONResult
// @Failure 500 {object} controller.JSONResult
// @Router /api/k8s/namespace [get]
func ListNamespace(c *gin.Context) {
	render := controller.Gin{C: c}
	var listNSQuery ListNSQuery
	if err := c.ShouldBindQuery(&listNSQuery); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	if _, ok := c.GetQuery("namespace"); !ok {
		res, err := configs.RestClient.CoreV1().Namespaces().
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: listNSQuery.LabelSelector,
				FieldSelector: listNSQuery.FieldSelector,
			})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		render.SetJson(res)
		return
	}
	res, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), listNSQuery.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetJson(res)
}
