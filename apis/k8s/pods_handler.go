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

type ListPodsQuery struct {
	Namespace     string `json:"namespace" form:"namespace"`
	Pod           string `json:"pod" form:"pod"`
	FieldSelector string `json:"field_selector" form:"field_selector"`
	LabelSelector string `json:"label_selector" form:"label_selector"`
}

// @Summary ListNamespaceAllSource
// @description 命名空间下Pod资源列表
// @Tags Kubernetes
// @Param namespace query ListPodsQuery false "namespace" default(test)
// @Param pod query ListPodsQuery false "pod"
// @Param field_selector query ListPodsQuery false "field_selector"
// @Param label_selector query ListPodsQuery false "label_selector"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/pods [get]
func ListNamespaceAllPods(c *gin.Context) {
	render := apis.Gin{C: c}
	var listPodsQuery ListPodsQuery
	if err := c.ShouldBindQuery(&listPodsQuery); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}

	if _, ok := c.GetQuery("namespace"); !ok {
		res, err := configs.RestClient.CoreV1().Pods(metaV1.NamespaceAll).
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: listPodsQuery.LabelSelector,
				FieldSelector: listPodsQuery.FieldSelector,
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
		Get(context.TODO(), listPodsQuery.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	if _, ok := c.GetQuery("pod"); !ok {
		res, err := configs.RestClient.CoreV1().Pods(listPodsQuery.Namespace).
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: listPodsQuery.LabelSelector,
				FieldSelector: listPodsQuery.FieldSelector,
			})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		render.SetJson(res)
		return
	}
	res, err := configs.RestClient.CoreV1().Pods(listPodsQuery.Namespace).
		Get(context.TODO(), listPodsQuery.Pod, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetJson(res)
}
