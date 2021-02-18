package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const downloadHtml = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>DownloadFromContainer</title>
	<script src="https://code.jquery.com/jquery-3.4.1.min.js" crossorigin="anonymous"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/css/select2.min.css" rel="stylesheet" />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/js/select2.min.js"></script>
    <script src="http://lib.h-ui.net/layer/3.1.1/layer.js"></script>
	<style>
        a:link {text-decoration:none;}
        .ge_table{ }
		.ge_table td{ height:44px; line-height:26px;}
		.hide_border{border: 0px;padding-left: 5px;line-height: 26px;width:205px;height: 26px;}
		.ge_input{border:1px solid #ccc;padding-left:5px;line-height:26px;width:205px; height:26px;}
		.short_select{
    		background:#fafdfe;
    		height:30px;
    		width:600px;
    		line-height:28px;
    		border:1px solid #9bc0dd;
    		-moz-border-radius:2px;
    		-webkit-border-radius:2px;
    		border-radius:4px;
		}
		#containers{
            width:680px;
            margin:20px auto;
            padding:15px;
            background-color:#eee;
          border-radius: 15px;
        }
		.fake-file-btn {
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
		.fake-file-btn:active {
			box-shadow: 0 1px 5px 1px rgba(0, 255, 255, 0.3) inset;
		}
		.fake-file-btn input[type=file] {
			position: absolute;
			font-size: 100px;
			right: 0;
			top: 0;
			opacity: 0;
			filter: alpha(opacity=0);
			cursor: pointer
		}
	</style>
</head>
<body>
	<!-- 伪装按钮实现，点击后选中文件后使用ajax立刻上传 start -->
	<div id="containers">
		<table cellpadding="0" cellspacing="0" border="0" class="ge_table">
			<tr>
				<td>命名空间: </td>
				<!-- td><input class="ge_input" type="text" id="namespace" value="default"></td -->
				<td>
					<select id="namespace" class="short_select"></select>
                    <script type="text/javascript">
                      $("#namespace").select2({placeholder: '请选择命名空间'});
                    </script>
				</td>
			</tr>
			<tr>
				<td>POD名称: </td>
				<!-- td><input class="ge_input" type="text" id="pod" value="nginx-test-76996486df-tdjdf"></td-->
				<td>
					<select id="pod" class="short_select"></select>
                    <script type="text/javascript">
                      $("#pod").select2({placeholder: '请选择命名空间'});
                    </script>
				</td>
			</tr>
			<tr>
				<td>容器名字: </td>
				<!-- td><input class="ge_input" type="text" id="container" value="nginx-0"></td -->
				<td>
					<select id="container" class="short_select"></select>
                    <script type="text/javascript">
                      $("#container").select2({
                        placeholder: '请选择容器'
                      });
                    </script>
				</td>
			</tr>
			<tr>
				<td>目标路径: </td>
				<td><input class="short_select" type="text" id="dest_path" value="/root"></td>
			</tr>
		</table>
        <p style="text-align:right">
        <span class="fake-file-btn">
			<a href="/">首页</a>
		</span>
		<span class="fake-file-btn" id="fake-file-btn">
			下载
		</span>
        </p>
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
		function getNamespace() {
			$("#namespace").empty();
			//获取数据
			$.ajax({
				type:"GET",
				contentType:"application/json",
				dataType:"json",
				url:"/api/k8s/namespace",
				dataType: "json",
				success:function(res){
					if (res.code != 200) {
						 layer.alert(res.info.message);
					}
					var data = res.data.items;
					for(var key in data){
						$("#namespace").append("<option value='" + data[key].metadata.name + "'>" + data[key].metadata.name + "</option>");
					};
                    $("#namespace").select2('val','1')
				},
                error: function() {
                     layer.alert("错误");
                }
			});
		};

		function getPods() {
			$("#pod").empty();
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
				url:"/api/k8s/pods?namespace="+namespace.text(),
				dataType: "json",
				success:function(res){
					if (res.code != 200) {
						 layer.alert(res.info.message);
					}
					var data = res.data.items;
					for(var key in data){
						$("#pod").append("<option value='" + data[key].metadata.name + "'>" + data[key].metadata.name + "</option>");
					};
                    $("#pod").select2('val','1')
				},
                error: function() {
                     layer.alert("错误");
                }
			});
		};

		function getContainer() {
			$("#container").empty();
			var namespace=$("#namespace option:selected");
			var pod=$("#pod option:selected");
			if (pod.text() == "") {
				 layer.alert('pod名称为空');
				return
			}
			//获取数据
			$.ajax({
				type:"GET",
				contentType:"application/json",
				dataType:"json",
				url:"/api/k8s/pods?namespace="+namespace.text()+"&pod="+pod.text(),
				dataType: "json",
				success:function(res){
					if (res.code != 200) {
						 layer.alert(res.info.message);
					}
					var data = res.data.spec.containers;
					for(var key in data){
						$("#container").append("<option value='" + data[key].name + "'>" + data[key].name + "</option>");
					};
                    $("#container").select2('val','1')
				},
                error: function() {
                     layer.alert("错误");
                }
			});
		};

		function downloadFiles() {
			var namespace=$("#namespace option:selected");
        	var pod = $("#pod option:selected");
			var container = $("#container option:selected");
			var destPath = document.getElementById("dest_path").value
			if (container.text() == "") {
				 layer.alert('容器名称为空');
				return
			}
			$.ajax({
				type: 'GET',
				url: "/api/k8s/download?namespace="+namespace.text()+"&pod_name="+pod.text()+"&container_name="+container.text()+"&dest_path="+destPath,
				timeout: 300 * 1000,
				processData: false,
				contentType: false,
				success : function(res){
					if (res.code!=200) {
						 layer.alert(res.info.message);
					}else{
						console.log('收到预期的json数据');
						location.href = res.data;
					}
					
				},
                error: function() {
                     layer.alert("错误");
                }
			});
		};

		// 页面打开即加载各资源列表
		getNamespace();
		
		// 选择命名空间后加载pod列表
		$("#namespace").on("select2:select", function() {
			getPods();
		});
		// 选择pod后加载pod列表
		$("#pod").on("select2:select", function() {
			getContainer();
		});

		$("#fake-file-btn").click(function() {
			downloadFiles();
		});
	</script>
	<!-- 伪装按钮实现，点击后选中文件后使用ajax立刻上传 end -->
</body>
</html>`

func DownloadHtml(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, downloadHtml)
}
