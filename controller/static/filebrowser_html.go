package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const fileBrowserHtml = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>FileBrowser</title>
	<script src="https://code.jquery.com/jquery-3.4.1.min.js" crossorigin="anonymous"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/css/select2.min.css" rel="stylesheet" />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/js/select2.min.js"></script>
    <link rel="stylesheet" href="http://www.yfdou.com/layui/css/layui.css"  media="all">
    <script src="http://lib.h-ui.net/layer/3.1.1/layer.js"></script>
    <script src="https://www.yfdou.com/layui/layui.js" charset="utf-8"></script>
    <style>
        a:link {text-decoration:none;}
        p{height:30px;}
        label{height:30px;}
		.hide_border{border: 0px;padding-left: 5px;line-height: 26px;width:205px;height: 26px;}
		.ge_input{border:1px solid #ccc;padding-left:5px;line-height:26px;width:205px; height:26px;}
		.short_select{
    		background:#fafdfe;
    		height:30px;
    		width:60px;
    		line-height:28px;
    		border:1px solid #9bc0dd;
    		-moz-border-radius:2px;
    		-webkit-border-radius:2px;
    		border-radius:4px;
		}
		#containers{
            /* width:680px; */
            margin:5px auto;
            padding:5px;
            background-color:#eee;
            border-radius: 15px;
        }
        		.fake-file-btn {
			position: relative;
			/* display: inline-block; */
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
  <div id="containers">
      <p>
        <label>命名空间: </label>
        <label>
          <select id="namespace" class="short_select"></select>
          <script type="text/javascript">
            $("#namespace").select2({placeholder: '请选择命名空间',allowClear: true,width: '160px'});
          </script>
        </label>
        &nbsp;&nbsp;
        <label>POD名称: </label>
        <label>
          <select id="pod" class="short_select"></select>
          <script type="text/javascript">
            $("#pod").select2({placeholder: '请选择pod',allowClear: true,width: '260px'});
          </script>
        </label>
        <label>容器名字: </label>
        <label>
          <select id="container" class="short_select"></select>
          <script type="text/javascript">
            $("#container").select2({width: '260px'});
          </script>
		</label>
        &nbsp;&nbsp;&nbsp;&nbsp;
        <label>
	      <span class="fake-file-btn" id="fake-btn">
            打开文件浏览器
	      </span>
          &nbsp;&nbsp;&nbsp;&nbsp;
          <span class="fake-file-btn">
	        <a href="/">首页</a>
	      </span>
        </label>
    </p>
    <table class="layui-hide" id="demo"></table>
    <script type="text/html" id="toolbarDemo">
      <div class="layui-btn-container">
        <button class="layui-btn layui-btn-sm" lay-event="getCheckData">获取选中行数据</button>
        <button class="layui-btn layui-btn-sm" lay-event="removeSelected">批量删除</button>
        <button class="layui-btn layui-btn-sm" lay-event="isAll">全选</button>
      </div>
    </script>
    <script type="text/html" id="barDemo">
      <a class="layui-btn layui-btn-xs" lay-event="edit">编辑</a>
      <a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del">删除</a>
    </script>
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
          $("#pod").empty();
          $("#container").empty();
		  //获取数据
		  $.ajax({
		  	type:"GET",
		  	contentType:"application/json",
		  	dataType:"json",
		  	url:"/api/k8s/namespace",
		  	dataType: "json",
		  	success:function(res){
		  	  if (res.code != 0) {
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
          $("#container").empty();
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
		  	  if (res.code != 0) {
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
		  	  if (res.code != 0) {
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

        function getPath(path) {
          if (path === undefined) {
            path = "/";
          };
          var namespace=$("#namespace option:selected");
          if (namespace.text() == "") {
		  	layer.alert('namespace为空');
		  	return
		  }
		  var pod=$("#pod option:selected");
          if (pod.text() == "") {
		  	layer.alert('pod为空');
		  	return
		  }
          var container=$("#container option:selected");
		  if (container.text() == "") {
		  	layer.alert('container为空');
		  	return
          }

          layui.use('table', function(){
            var table = layui.table;
            table.render({
              elem: '#demo',
              url:"/api/k8s/file_browser?namespace="+namespace.text()+"&pods="+pod.text()+"&container="+container.text()+"&path="+path,
              toolbar: '#toolbarDemo',
              totalRow: true,
              cellMinWidth: 80,
              cols: [[ //标题栏
                {type:'checkbox',fixed: 'left'},
                {field: 'Name', title: '名称',fixed: 'left', unresize: true, sort: true, totalRowText: '合计'},
                {field: 'Mode', title: '权限'},
                {field: 'Size', title: '大小',sort: true},
                {field: 'ModTime', title: '最后修改时间',sort: true},
                {field: 'IsDir', title: '是否目录',sort: true,totalRow: true},
                {fixed: 'right', title:'操作', toolbar: '#barDemo', width:150}
              ]]
            });
            //工具栏事件
            table.on('toolbar(demo)', function(obj){
              var checkStatus = table.checkStatus(obj.config.id);
              switch(obj.event){
                case 'getCheckData':
                  var data = checkStatus.data;
                  layer.alert(JSON.stringify(data));
                break;
                case 'removeSelected':
                  layer.msg('选中了：'+ data.length + ' 个');
                break;
                case 'isAll':
                  layer.msg(checkStatus.isAll ? '全选': '未全选')
                break;
              };
            });
            //监听行单击事件（双击事件为：rowDouble）
            table.on('row(demo)', function(obj){
              var data = obj.data;
    
              layer.alert(JSON.stringify(data), {
                title: '当前行数据：'
              });
    
              //标注选中样式
              obj.tr.addClass('layui-table-click').siblings().removeClass('layui-table-click');
              });
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
        $("#fake-btn").unbind("click").click(function() {
            // layer.alert("打开文件浏览器");
            getPath("/")
		});
    </script>
  </div>
</body>
</html>`

func FileBrowserHtml(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fileBrowserHtml)
}
