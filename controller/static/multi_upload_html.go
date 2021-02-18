package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const multiUploadHtml = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>UploadToContainer</title>
    <script src="https://code.jquery.com/jquery-3.4.1.min.js" crossorigin="anonymous"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/css/select2.min.css" rel="stylesheet" />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/js/select2.min.js"></script>
    <script src="http://lib.h-ui.net/layer/3.1.1/layer.js"></script>
	<style>
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
        a:link {text-decoration:none;}
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
					<select id="pod" class="short_select" multiple="multiple" style="height:300px" size="30"></select>
                    <script type="text/javascript">
                      $("#pod").select2({
                        placeholder: '请选择pod，最多选择20个',
                        allowClear: true,
                        maximumSelectionLength:20,
                      });
                    </script>
				</td>
			</tr>
			<tr>
				<td>目标路径: </td>
				<td><input class="short_select" type="text" id="dest_path" value="/root/"></td>
			</tr>
		</table>
        <p style="text-align:right">
        <span class="fake-file-btn">
			<a href="/">首页</a>
		</span>
		<span class="fake-file-btn" id="fileFolderOne-btn">
			上传文件<input type="file" id="fileFolderOne" name="files" multiple="true">
		</span>
        <span class="fake-file-btn" id="fileFolderMore-btn">
			上传文件夹<input type="file" id="fileFolderMore" name="files" webkitdirectory mozdirectory accept="*/*">
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

        function fileFolderOne() {
			var files = document.getElementById('fileFolderOne').files;
			var namespace=$("#namespace option:selected");
        	var pod = $("#pod option:selected");
			var container = $("#container option:selected");
			var destPath = document.getElementById("dest_path").value;
			var podArray = new Array();
			$("#pod option:selected").each(function(index, obj) {
				podArray.push($(obj).text());
			});
			var podStr=podArray.join(",");

			//新建一个FormData对象
			var formData = new FormData();
			console.log(files);
			//追加文件数据
			for (i = 0; i < files.length; i++) {
				formData.append("files", files[i]);
			}

			$.ajax({
				type: 'POST',
				url: "/api/k8s/multi_upload?namespace="+namespace.text()+"&pod_name="+podStr+"&dest_path="+destPath,
				timeout: 300 * 1000,
				data: formData,
				processData: false,
				contentType: false,
				success: function(res) {
					if (res.code != 200) {
						 layer.alert(res.info.message);
					}else{
                        $("#fileFolderOne").empty();
                        if ((typeof(res.data.success)!="undefined")&&(typeof(res.data.failure)!="undefined")) {
                          layer.alert(res.data.success+"<br>"+res.data.failure);
                          return;
                        }
						if (typeof(res.data.success)!="undefined") {
                          layer.alert(res.data.success);
                          return;
                        }
                        if (typeof(res.data.failure)!="undefined") {
                          layer.alert(res.data.failure);
                          return;
                        }
					}
				},
                error: function() {
                     layer.alert("错误");
                }
			});
		};

		function fileFolderMore() {
			var files = document.getElementById('fileFolderMore').files;
			var namespace=$("#namespace option:selected");
        	var pod = $("#pod option:selected");
			var container = $("#container option:selected");
			var destPath = document.getElementById("dest_path").value;
			var podArray = new Array();
			$("#pod option:selected").each(function(index, obj) {
				podArray.push($(obj).text());
			});
			var podStr=podArray.join(",");

			//新建一个FormData对象
			var formData = new FormData();
			console.log(files);
			//追加文件数据
			for (i = 0; i < files.length; i++) {
				formData.append("files", files[i]);
			}

			$.ajax({
				type: 'POST',
				url: "/api/k8s/multi_upload?namespace="+namespace.text()+"&pod_name="+podStr+"&dest_path="+destPath,
				timeout: 300 * 1000,
				data: formData,
				processData: false,
				contentType: false,
				success: function(res) {
					if (res.code != 200) {
						 layer.alert(res.info.message);
					}else{
                        $("#fileFolderMore").empty();
                        if ((typeof(res.data.success)!="undefined")&&(typeof(res.data.failure)!="undefined")) {
                          layer.alert(res.data.success+"<br>"+res.data.failure);
                          return;
                        }
						if (typeof(res.data.success)!="undefined") {
                          layer.alert(res.data.success);
                          return;
                        }
                        if (typeof(res.data.failure)!="undefined") {
                          layer.alert(res.data.failure);
                          return;
                        }
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
        // 上传多文件
		$("#fileFolderOne-btn").on("change", "input[type='file']", function() {
			fileFolderOne();
		});
        // 上传多文件夹
        $("#fileFolderMore-btn").on("change", "input[type='file']", function() {
			fileFolderMore();
		});
	</script>
	<!-- 伪装按钮实现，点击后选中文件后使用ajax立刻上传 end -->
</body>
</html>`

func MultiUploadHtml(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, multiUploadHtml)
}
