// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/k8s/deployment": {
            "get": {
                "description": "命名空间下Deployment资源列表",
                "tags": [
                    "Kubernetes"
                ],
                "summary": "ListNamespaceAllDeployment",
                "parameters": [
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        },
        "/api/k8s/download": {
            "get": {
                "description": "从容器下载到本地",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Copy2Local",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "namespace",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "nginx-test-76996486df-tdjdf",
                        "description": "pod_name",
                        "name": "pod_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "nginx-0",
                        "description": "container_name",
                        "name": "container_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "/root",
                        "description": "dest_path",
                        "name": "dest_path",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        },
        "/api/k8s/multi_upload": {
            "post": {
                "description": "上传到容器",
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "MultiCopy2Container",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "namespace",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "nginx-test-76996486df-tdjdf",
                        "description": "pod_name",
                        "name": "pod_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "/root/",
                        "description": "dest_path",
                        "name": "dest_path",
                        "in": "query"
                    },
                    {
                        "type": "file",
                        "description": "files",
                        "name": "files",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        },
        "/api/k8s/namespace": {
            "get": {
                "description": "命名空间列表",
                "tags": [
                    "Kubernetes"
                ],
                "summary": "ListNamespace",
                "parameters": [
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        },
        "/api/k8s/pods": {
            "get": {
                "description": "命名空间下Pod资源列表",
                "tags": [
                    "Kubernetes"
                ],
                "summary": "ListNamespaceAllSource",
                "parameters": [
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "pod",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "pod",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "pod",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "pod",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        },
        "/api/k8s/status": {
            "get": {
                "description": "获取pod中container状态",
                "tags": [
                    "Kubernetes"
                ],
                "summary": "PodStatus",
                "parameters": [
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "deployment",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "field_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "label_selector",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        },
        "/api/k8s/upload": {
            "post": {
                "description": "上传到容器",
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Copy2Container",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "namespace",
                        "name": "namespace",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "nginx-test-76996486df-tdjdf",
                        "description": "pod_name",
                        "name": "pod_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "nginx-0",
                        "description": "container_name",
                        "name": "container_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "/root/",
                        "description": "dest_path",
                        "name": "dest_path",
                        "in": "query"
                    },
                    {
                        "type": "file",
                        "description": "files",
                        "name": "files",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.JSONResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.Info": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "消息"
                },
                "ok": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "controller.JSONResult": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 0
                },
                "data": {
                    "type": "object"
                },
                "info": {
                    "$ref": "#/definitions/controller.Info"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "KubeCp Swagger",
	Description: "网页版kubectl cp",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}