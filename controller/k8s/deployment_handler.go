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

type ListDeploymentQuery struct {
	Namespace     string `json:"namespace" form:"namespace"`
	Deployment    string `json:"deployment" form:"deployment"`
	FieldSelector string `json:"field_selector" form:"field_selector"`
	LabelSelector string `json:"label_selector" form:"label_selector"`
}

// @Summary ListNamespaceAllDeployment
// @description 命名空间下Deployment资源列表
// @Tags Kubernetes
// @Param namespace query ListDeploymentQuery false "namespace" default(test)
// @Param deployment query ListDeploymentQuery false "deployment"
// @Param field_selector query ListDeploymentQuery false "field_selector"
// @Param label_selector query ListDeploymentQuery false "label_selector"
// @Success 200 {object} controller.JSONResult
// @Failure 500 {object} controller.JSONResult
// @Router /api/k8s/deployment [get]
func ListNamespaceAllDeployment(c *gin.Context) {
	render := controller.Gin{C: c}
	var listDeploymentQuery ListDeploymentQuery
	if err := c.ShouldBindQuery(&listDeploymentQuery); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}

	if _, ok := c.GetQuery("namespace"); !ok {
		res, err := configs.RestClient.AppsV1().Deployments(metaV1.NamespaceAll).
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: listDeploymentQuery.LabelSelector,
				FieldSelector: listDeploymentQuery.FieldSelector,
			})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		render.SetJson(res)
		return
	}
	_, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), listDeploymentQuery.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	if _, ok := c.GetQuery("deployment"); !ok {
		res, err := configs.RestClient.AppsV1().Deployments(listDeploymentQuery.Namespace).
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: listDeploymentQuery.LabelSelector,
				FieldSelector: listDeploymentQuery.FieldSelector,
			})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		render.SetJson(res)
		return
	}
	res, err := configs.RestClient.AppsV1().Deployments(listDeploymentQuery.Namespace).
		Get(context.TODO(), listDeploymentQuery.Deployment, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetJson(res)
}
