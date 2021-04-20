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
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/deployment [get]
func ListNamespaceAllDeployment(c *gin.Context) {
	render := apis.Gin{C: c}
	var query ListDeploymentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}

	if _, ok := c.GetQuery("namespace"); !ok {
		res, err := configs.RestClient.AppsV1().Deployments(metaV1.NamespaceAll).
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
	_, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	if _, ok := c.GetQuery("deployment"); !ok {
		res, err := configs.RestClient.AppsV1().Deployments(query.Namespace).
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
	res, err := configs.RestClient.AppsV1().Deployments(query.Namespace).
		Get(context.TODO(), query.Deployment, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetJson(res)
}
