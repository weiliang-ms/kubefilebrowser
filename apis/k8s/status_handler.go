package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	deploymentutil "k8s.io/kubernetes/pkg/controller/deployment/util"
	"kubefilebrowser/apis"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
	"strings"
	"sync"
)

type StatusQuery struct {
	Namespace     string   `json:"namespace" form:"namespace" binding:"required"`
	Deployment    []string `json:"deployment" form:"deployment" `
	FieldSelector string   `json:"field_selector" form:"field_selector"`
	LabelSelector string   `json:"label_selector" form:"label_selector"`
}

type ResPods struct {
	PodName    string         `json:"pod_name"`
	Containers []ResContainer `json:"containers,omitempty"`
	Error      error          `json:"error,omitempty"`
}

type ResContainer struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Image           string `json:"image,omitempty"`
	ImagePullPolicy string `json:"image_pull_policy,omitempty"`
	State           string `json:"state,omitempty"`
	Restart         int32  `json:"restart,omitempty"`
	Cpu             string `json:"cpu,omitempty"`
	Ram             string `json:"ram,omitempty"`
	Version         string `json:"version,omitempty"`
	Os              string `json:"os,omitempty"`
	Arch            string `json:"arch,omitempty"`
}

// PodStatus
// @Summary PodStatus
// @description 获取pod中container状态
// @Tags Kubernetes
// @Param namespace query StatusQuery true "namespace" default(test)
// @Param deployment query StatusQuery false "deployment"
// @Param field_selector query StatusQuery false "field_selector"
// @Param label_selector query StatusQuery false "label_selector"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/status [get]
func PodStatus(c *gin.Context) {
	render := apis.Gin{C: c}
	var query StatusQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}

	if _, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{}); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	deployments, err := query.deploymentList(c)
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	var resPods []ResPods
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, p := range query.podsList(deployments) {
		wg.Add(1)
		go func(w *sync.WaitGroup, pod coreV1.Pod) {
			defer w.Done()
			mu.Lock()
			defer mu.Unlock()
			resPods = append(resPods, query.podStats(pod))
		}(&wg, p)
	}
	wg.Wait()
	render.SetJson(resPods)
}

func (query *StatusQuery) deploymentList(c *gin.Context) ([]appsV1.Deployment, error) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var deployments []appsV1.Deployment
	_d, ok := c.GetQuery("deployment")
	if ok && _d != "" {
		for _, d := range query.Deployment {
			wg.Add(1)
			go func(wg *sync.WaitGroup, name string) {
				defer wg.Done()
				deployment, err := configs.RestClient.AppsV1().Deployments(query.Namespace).
					Get(context.TODO(), name, metaV1.GetOptions{})
				if err != nil {
					logs.Error(err)
					return
				}
				mu.Lock()
				defer mu.Unlock()
				deployments = append(deployments, *deployment)
			}(&wg, d)
		}
		wg.Wait()
	} else {
		deploymentList, err := configs.RestClient.AppsV1().Deployments(query.Namespace).
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: query.LabelSelector,
				FieldSelector: query.FieldSelector,
			})
		if err != nil {
			return nil, err
		}
		for _, deployment := range deploymentList.Items {
			deployments = append(deployments, deployment)
		}
	}
	return deployments, nil
}

func (query *StatusQuery) podsList(deployments []appsV1.Deployment) []coreV1.Pod {
	var podList []coreV1.Pod
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, d := range deployments {
		wg.Add(1)
		go func(wg *sync.WaitGroup, deployment appsV1.Deployment) {
			defer wg.Done()
			rsList, err := deploymentutil.ListReplicaSets(&deployment, deploymentutil.RsListFromClient(configs.RestClient.AppsV1()))
			if err != nil {
				logs.Error(err)
				return
			}
			podListFunc := func(namespace string, options metaV1.ListOptions) (*coreV1.PodList, error) {
				return configs.RestClient.CoreV1().Pods(namespace).List(context.TODO(), options)
			}
			pods, err := deploymentutil.ListPods(&deployment, rsList, podListFunc)
			if err != nil {
				logs.Error(err)
				return
			}
			mu.Lock()
			defer mu.Unlock()
			podList = append(podList, pods.Items...)
		}(&wg, d)
	}
	wg.Wait()
	return podList
}

func (query *StatusQuery) podStats(pod coreV1.Pod) ResPods {
	var osType = "linux"
	var arch = "amd64"

	// get pod system arch and type
	node, err := configs.RestClient.CoreV1().Nodes().
		Get(context.TODO(), pod.Spec.NodeName, metaV1.GetOptions{})
	if err == nil {
		if node.Labels["beta.kubernetes.io/os"] != "" {
			osType = node.Labels["beta.kubernetes.io/os"]
		} else if node.Labels["kubernetes.io/os"] != "" {
			osType = node.Labels["kubernetes.io/os"]
		}
		if node.Labels["beta.kubernetes.io/arch"] != "" {
			arch = node.Labels["beta.kubernetes.io/arch"]
		} else if node.Labels["kubernetes.io/arch"] != "" {
			arch = node.Labels["kubernetes.io/arch"]
		}
	}

	var containerMetrics = make(map[string]map[string]string)
	podMetrics, err := configs.MetricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).
		Get(context.Background(), pod.Name, metaV1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return ResPods{Error: err}
	}
	for _, container := range podMetrics.Containers {
		cpuQuantity := container.Usage.Cpu().AsDec().String()
		memQuantity, _ := container.Usage.Memory().AsInt64()
		containerMetrics[fmt.Sprintf("%s.%s", pod.Name, container.Name)] = map[string]string{
			"cpu": cpuQuantity,
			"mem": utils.FormatFileSize(memQuantity),
		}
	}

	var container []ResContainer
	if len(pod.Spec.Containers) > 0 {
		for _k, _v := range pod.Status.ContainerStatuses {
			var state string
			if _v.Ready {
				state = "Running"
			} else {
				state = "Error"
			}
			_i := strings.Split(_v.Image, ":")
			_container := ResContainer{
				ID:              _k + 1,
				Name:            _v.Name,
				Image:           _v.Image,
				State:           state,
				Restart:         _v.RestartCount,
				ImagePullPolicy: fmt.Sprint(pod.Spec.Containers[_k].ImagePullPolicy),
				Version:         _i[len(_i)-1],
				Os:              osType,
				Arch:            arch,
			}
			metrics, ok := containerMetrics[fmt.Sprintf("%s.%s", pod.Name, _v.Name)]
			if ok {
				_container.Cpu = metrics["cpu"]
				_container.Ram = metrics["mem"]
			}
			container = append(container, _container)
		}
	}
	return ResPods{
		PodName:    pod.Name,
		Containers: container,
	}
}
