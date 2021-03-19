package controller

import (
	"github.com/gin-gonic/gin"
	"kubecp/configs"
	"kubecp/utils"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type JSONResult struct {
	Code int         `json:"code" description:"返回码" example:"0000"`
	Info *Info       `json:"info" description:"信息"`
	Data interface{} `json:"data" description:"数据"`
}

type Info struct {
	Ok      bool   `json:"ok" description:"状态" example:"true"`
	Message string `json:"message,omitempty" description:"消息" example:"消息"`
}

func NewRes(data interface{}, err error, code int) *JSONResult {
	if code == 200 {
		code = 0
	}
	codeMsg := ""
	if configs.Config.RunMode == "release" && code != 0 {
		codeMsg = utils.GetMsg(code)
	}

	return &JSONResult{
		Data: data,
		Code: code,
		Info: func() *Info {
			result := NewInfo(err)
			if _, ok := data.(string); ok && len(data.(string)) != 0 {
				result.Message = data.(string)
			}
			if codeMsg != "" {
				result.Message = codeMsg + ",错误信息: " + err.Error()
			}
			return result
		}(),
	}
}

func NewInfo(err error) (info *Info) {
	info = new(Info)
	if err == nil {
		info.Ok = true
		return
	}
	info.Message = err.Error()
	return
}

// Response res
func (g *Gin) SetRes(res interface{}, err error, code int) {
	g.C.JSON(http.StatusOK, NewRes(res, err, code))
}

// Set Json
func (g *Gin) SetJson(res interface{}) {
	g.SetRes(res, nil, utils.CODE_SUCCESS)
}

// Check Error
func (g *Gin) SetError(code int, err error) {
	g.SetRes(nil, err, code)
	g.C.Abort()
}
