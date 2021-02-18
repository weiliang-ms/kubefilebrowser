package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const statusHtml = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Status</title>
    <script src="https://code.jquery.com/jquery-3.4.1.min.js" crossorigin="anonymous"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/css/select2.min.css" rel="stylesheet" />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/js/select2.min.js"></script>
    <script src="http://lib.h-ui.net/layer/3.1.1/layer.js"></script>
    <style>
        html, body {
            margin:0px;
            width: 100%;
            height: 100%;
            background-color:#eee;
        }
        .top-center-bottom>div{
		  position:absolute;
		}
		.top-center-bottom .top{
		  top:0;
		  height:100px;
		  width:100%;
		}
		.top-center-bottom .center{
		  bottom:100px;
		  top:100px;
		  width:100%;
          min-height:400px;
		}
        #container{
            height:100%;
            margin:20px;
            padding:15px;
            border-radius: 15px;
        }
        .left {
            float: left;
            width: 270px;
            height: 300px;
        }
        .right {
            margin-left: 280px;
            height: 300px;
        }
        .short_select{
    		background:#fafdfe;
    		height:28px;
    		width:252px;
    		line-height:28px;
    		border:1px solid #9bc0dd;
    		-moz-border-radius:2px;
    		-webkit-border-radius:2px;
    		border-radius:4px;
		}
        /*表格样式*/
        table {
            width: 100%;
            background: #ccc;
            margin: 10px auto;
            border-collapse: collapse;
            /*border-collapse:collapse合并内外边距
            (去除表格单元格默认的2个像素内外边距*/
        }
        th,td {
            height: 25px;
            line-height: 25px;
            text-align: left;
            border: 1px solid #ccc;
        }
        th {
            background: #eee;
            font-weight: normal;
        }
        tr {
            background: #fff;
        }
        tr:hover {
            background: #cc0;
        }
        td a {
            color: #06f;
            text-decoration: none;
        }
        td a:hover {
            color: #06f;
            text-decoration: underline;
        }
		.fake-btn {
			position: relative;
			display: inline-block;
			background: #D0EEFF;
			border: 1px solid #99D3F5;
			border-radius: 4px;
			padding: 4px 12px;
			overflow: hidden;
			color: #1E88C7;
			text-decoration: none;
			text-indent: 0;
			line-height: 20px;
		}
		.fake-btn:active {
			box-shadow: 0 1px 5px 1px rgba(0, 255, 255, 0.3) inset;
		}
        .main{
	        text-align: center; /*让div内部文字居中*/
	        border-radius: 20px;
	        margin: auto;
	        position: absolute;
	        top: 0;
	        left: 0;
	        right: 0;
	        bottom: 0;
      }
    </style>
    <script>
        $(function() {
			$("#namespace").empty();
            $("#deployment").empty();
            $("#tbody").empty();
			//获取数据
			$.ajax({
				type:"GET",
				contentType:"application/json",
				dataType:"json",
				url:"/api/k8s/namespace",
				dataType: "json",
				success:function(res){
					if (res.code != 200) {
						alert(res.info.message);
					}
					var data = res.data.items;
					for(var key in data){
						$("#namespace").append("<option value='" + data[key].metadata.name + "'>" + data[key].metadata.name + "</option>");
					};
                    $("#namespace").select2('val','1')
				}
			});
		});
    </script>
</head>
<body>
    <div class="top-center-bottom">
        <div class="top">
          <div class="main">
            <br/>
            <input class="fake-btn" type="button" value="首页" onclick="javascrtpt:window.location.href='/'">&nbsp;
            <input class="fake-btn" type="button" value="上传到单个pod内多容器" onclick="javascrtpt:window.location.href='/upload'">&nbsp;
            <input class="fake-btn" type="button" value="上传到多个pod内所有容器" onclick="javascrtpt:window.location.href='/multi_upload'">&nbsp;
            <input class="fake-btn" type="button" value="从容器内下载" onclick="javascrtpt:window.location.href='/download'">&nbsp;
          </div>
        </div>
        <div class="center">
          <div class="left">
            <p style="margin:0px 0px 0px 15px">
              <label for="namespace">命名空间:&nbsp;&nbsp;&nbsp;&nbsp;</label>
              <select id="namespace" class="short_select" name="namespace"></select>
              <script type="text/javascript">
                $("#namespace").select2({placeholder: '请选择命名空间'});
              </script>
              <br>
              <label for="deployment">无状态集:&nbsp;&nbsp;&nbsp;&nbsp;</label>
              <select id="deployment" class="short_select" name="deployment" multiple="multiple" style="height:600px" size="50"></select>
              <script type="text/javascript">
                $("#deployment").select2({
                  placeholder: '请选择deployment，最多选择20个',
                  allowClear: true,
                  maximumSelectionLength:20,
                });
              </script>
              <br>
              <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>
              <br>
              <span class="fake-btn" id="fake-btn">
			    查看
		      </span>
            </p>
          </div>
          <div class="right">
            <table id="status">
                <thead>
                    <tr>
                        <th>pod名称</th>
                        <th>容器名称</th>
                        <th>镜像地址</th>
                        <th>镜像版本</th>
                        <th>镜像拉取策略</th>
                        <th>运行状态</th>
                        <th>cpu使用</th>
                        <th>内存使用</th>
                    </tr>
                </thead>
                <tbody id="tbody"></tbody>
            </table>
          </div>
        </div>
    </div>
    <script>
        $.ajaxSetup({
          layerIndex:-1, //保存当前请求对应的提示框index,用于后面关闭使用
          //在请求显示提示框
          beforeSend: function(jqXHR, settings) {
            this.layerIndex = layer.load(1);
          },
          //请求完毕后（不管成功还是失败），关闭提示框
          complete: function () {
            layer.close(this.layerIndex);
          },
          //请求失败时，弹出错误信息
          error: function (jqXHR, status, e) {
            layer.alert('数据请求失败，请后再试!');
          }
        });
        function getDeployments() {
			$("#deployment").empty();
            $("#tbody").empty();
			var namespace=$("#namespace option:selected");
			if (namespace.text() == "") {
				 layer.alert('命名空间为空');
				return
			}
			//获取数据
			$.ajax({
				type:"GET",
				contentType:"application/json",
				dataType:"json",
				url:"/api/k8s/deployment?namespace="+namespace.text(),
				dataType: "json",
				success:function(res){
					if (res.code != 200) {
						alert(res.info.message);
					}
					var data = res.data.items;
					for(var key in data){
						$("#deployment").append("<option value='" + data[key].metadata.name + "'>" + data[key].metadata.name + "</option>");
					};
                    $("#deployment").select2('val','1')
				},
                error: function() {
                   layer.alert("错误");
                }
			});
		};
        function statusTable() {
			$("#tbody").empty();
			var namespace=$("#namespace option:selected");
			if (namespace.text() == "") {
				 layer.alert('命名空间不能为空');
				return
			}
            var deployment=$("#deployment option:selected");
			if (deployment.text() == "") {
				 layer.alert('deployment不能为空');
				return
			}
            var deploymentArray = new Array();
			$("#deployment option:selected").each(function(index, obj) {
				deploymentArray.push($(obj).text());
			});
			var deploymentStr=deploymentArray.join(",");
			//获取数据
			$.ajax({
				type:"GET",
				contentType:"application/json",
				dataType:"json",
				url:"/api/k8s/status?namespace="+namespace.text()+"&deployment="+deploymentStr,
				dataType: "json",
				success:function(res) {
					if (res.code != 200) {
						 layer.alert(res.info.message);
                        return;
					};
                    if (res.data == null) {
                         layer.alert("亲, 没找到你想要的数据, 请换一个试试");
                        return;
                    };
                    var rData = res.data;
                    for(var i in rData){
                        var pod_name =rData[i].pod_name;
                        var cData = rData[i].containers;
                        for(var j in cData) {
                            html = $("<tr></tr>");
                            var str = ""
                            if (j==0) {
                                str = "<td>" + pod_name + "</td>";
                            } else {
                                str = "<td></td>";
                            };
                            html.append(str);
                            html.append("<td>" + cData[j].name + "</td>");
                            html.append("<td>" + cData[j].image + "</td>");
                            html.append("<td>" + cData[j].version + "</td>");
                            html.append("<td>" + cData[j].image_pull_policy + "</td>");
                            html.append("<td>" + cData[j].state + "</td>");
                            html.append("<td>" + cData[j].cpu + "</td>");
                            html.append("<td>" + cData[j].ram + "</td>");
                            html.appendTo("#tbody");
                        };
                    };
				},
                error: function() {
                     layer.alert("错误");
                }
			});
		};
        // 选择命名空间后加载Deployment列表
        $("#namespace").on("select2:select", function() {
            //console.log("");
			getDeployments();
		});
        $("#fake-btn").unbind("click").click(function() {
            //console.log("");
			statusTable();
		});
    </script>
</body>
</html>`

func StatusHtml(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, statusHtml)
}
