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
        /** 新增按钮 **/
        #addPath{
            padding:3px;
            display:inline-block;
            background-color:#5ac7d0;
            color:#f1f1f1;
            border-radius: 4px;
        }
        /** 删除按钮 **/
        .removeVar{
            margin:auto;
            padding:5px;
            display:inline-block;
            background-color:#B02109;
            color:#f1f1f1;
            border:1px solid #005;
            border-radius: 4px;
            style="line-height: 15px";
        }
 
        #addPath:hover, .removeVar:hover{
            cursor: pointer;
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
		</table>
        <p style="line-height: 20px">
          <span><a href="#" rel="external nofollow" rel="external nofollow" rel="external nofollow" id="AddMoreFileBox" class="btn btn-info fake-file-btn">添加目标路径</a></span>
          <div id="InputsWrapper">
            <div><input type="text" name="dest_path[]" id="field_1" style="width:300px" class="short_select" value=""/>&nbsp;&nbsp;<a href="#" rel="external nofollow" rel="external nofollow" rel="external nofollow" class="removeVar"><input type='button' value='删除'></a></div>
          </div>
        </p>
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
            layer.alert('数据请求失败，请后再试!', {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
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
					if (res.code != 0) {
						 layer.alert(res.info.message, {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
					}
					var data = res.data.items;
					for(var key in data){
						$("#namespace").append("<option value='" + data[key].metadata.name + "'>" + data[key].metadata.name + "</option>");
					};
                    $("#namespace").select2('val','1')
				},
                error: function() {
                     layer.alert("错误", {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
                }
			});
		};

		function getPods() {
			$("#pod").empty();
			var namespace=$("#namespace option:selected");
			if (namespace.text() == "") {
				 layer.alert('命名空间为空', {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
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
					if (res.code != 0) {
						 layer.alert(res.info.message, {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
					}
					var data = res.data.items;
					for(var key in data){
						$("#pod").append("<option value='" + data[key].metadata.name + "'>" + data[key].metadata.name + "</option>");
					};
                    $("#pod").select2('val','1')
				},
                error: function() {
                     layer.alert("错误", {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
                }
			});
		};

		function getContainer() {
			$("#container").empty();
			var namespace=$("#namespace option:selected");
			var pod=$("#pod option:selected");
			if (pod.text() == "") {
				 layer.alert('pod名称为空', {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
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
					if (res.code != 0) {
						 layer.alert(res.info.message, {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
					}
					var data = res.data.spec.containers;
					for(var key in data){
						$("#container").append("<option value='" + data[key].name + "'>" + data[key].name + "</option>");
					};
                    $("#container").select2('val','1')
				},
                error: function() {
                     layer.alert("错误", {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
                }
			});
		};

        $(document).ready(function() {
            var MaxInputs    = 8; //maximum input boxes allowed
            var InputsWrapper  = $("#InputsWrapper"); //Input boxes wrapper ID
            var AddButton    = $("#AddMoreFileBox"); //Add button ID
            var x = InputsWrapper.length; //initlal text box count
            var FieldCount=1; //to keep track of text box added
            $(AddButton).click(function (e) //on add input button click
            {
                if(x <= MaxInputs) //max input box allowed
                {
                    FieldCount++; //text box added increment
                    //add input box
                    $(InputsWrapper).append('<div><input type="text" name="dest_path[]" id="field_'+ FieldCount +'" style="width:300px" class="short_select" value=""/>&nbsp;&nbsp;<a href="#" rel="external nofollow" rel="external nofollow" rel="external nofollow" class="removeVar"><input type="button" value="删除"></a></div>');
                        x++; //text box increment
                    }
                    return false;
                });
                $("body").on("click",".removeVar", function(e){ //user click on remove text
                     if( x > 1 ) {
                         $(this).parent('div').remove(); //remove text box
                         x--; //decrement textbox
                }
                return false;
            })
        });

		function downloadFiles() {
			var namespace=$("#namespace option:selected");
        	var pod = $("#pod option:selected");
			var container = $("#container option:selected");
			var destPath = $("*[name='dest_path']").val();
            var pathArr = new Array;
            var dest_path=new String;
            $("input[name='dest_path[]']").each(function(index, item){
                pathArr[index] = $(this).val();
                dest_path = pathArr.join('&dest_path=');
            });
			if (container.text() == "") {
				layer.alert('容器名称为空', {skin: 'layui-layer-molv',closeBtn: 1,shadeClose: true,anim: 1,title:"提示",icon: 6});
				return
			}
			var url = "/api/k8s/download?namespace="+namespace.text()+"&pod_name="+pod.text()+"&container_name="+container.text()+"&dest_path="+dest_path;
            var xhr = new XMLHttpRequest();
            xhr.open('GET', url, true);        // 也可以使用POST方式，根据接口
            xhr.responseType = "blob";    // 返回类型blob
            xhr.setRequestHeader("Content-type", "application/json;charset=UTF-8");
            // 定义请求完成的处理函数，请求前也可以增加加载框/禁用下载按钮逻辑
            xhr.onload = function () {
              // 请求完成
              if (this.status === 200) {
                // 返回200
                let content = xhr.response;
                let eLink = document.createElement('a');
                var fileName = this.getResponseHeader('X-File-Name');
                eLink.download = fileName;
                eLink.style.display = 'none';
                let blob = new Blob([content]);
                eLink.href = URL.createObjectURL(blob);
                document.body.appendChild(eLink);
                eLink.click();
                document.body.removeChild(eLink);
              }
            };
            // 发送ajax请求
            xhr.send()
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
