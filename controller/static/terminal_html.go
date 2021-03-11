package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const terminalHtml = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Container Terminal</title>
	<script src="https://code.jquery.com/jquery-3.4.1.min.js" crossorigin="anonymous"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/css/select2.min.css" rel="stylesheet" />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/select2/4.0.2/js/select2.min.js"></script>
    <script src="http://lib.h-ui.net/layer/3.1.1/layer.js"></script>
    <link rel="stylesheet" href="https://www.yfdou.com/xterm/dist/xterm.css" />
    <script src="https://www.yfdou.com/xterm/dist/xterm.js"></script>
    <script src="https://www.yfdou.com/xterm/dist/addons/fit/fit.js"></script>
    <script src="https://www.yfdou.com/xterm/dist/addons/winptyCompat/winptyCompat.js"></script>
    <script src="https://www.yfdou.com/xterm/dist/addons/webLinks/webLinks.js"></script>
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
            $("#container").select2({
              placeholder: '请选择容器，最多选择1个',
              allowClear: true,
              maximumSelectionLength:1,
              width: '260px',
            });
          </script>
		</label>
          &nbsp;&nbsp;
          <label >shell: </label>
        <label>
          <select id="shell" name="shell">
            <option value="sh">sh</option>
            <option value="bash">bash</option>
            <option value="cmd">cmd</option>
          </select >
          <script type="text/javascript">
            $("#shell").select2({
              placeholder: '选择',
              allowClear: true,
              width: '80px',
            });
          </script>
        </label>
    &nbsp;&nbsp;&nbsp;&nbsp;
    <label>
	  <span class="fake-file-btn" id="fake-btn">
		连接
	  </span>
      &nbsp;&nbsp;&nbsp;&nbsp;
      <span class="fake-file-btn">
	    <a href="/">首页</a>
	  </span>
    </label>
    </p>
    <div id="terminal"></div>
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
        function connTerminal() {
          document.getElementById('terminal').innerHTML = ""
          // 获取要连接的容器信息
          var namespace=$("#namespace option:selected");
          var pod = $("#pod option:selected");
          var container = $("#container option:selected");
          var shell = $("#shell option:selected");
          if (container.text() == "") {
		  	layer.alert('容器不能为空');
		  	return
		  }
          // xterm配置自适应大小插件
          Terminal.applyAddon(fit);
          
          // 这俩插件不知道干嘛的, 用总比不用好
          Terminal.applyAddon(winptyCompat)
          Terminal.applyAddon(webLinks)
          
          // 创建终端
          var term = new Terminal();
          term.open(document.getElementById('terminal'));
          
          // 使用fit插件自适应terminal size
          term.fit();
          term.winptyCompatInit()
          term.webLinksInit()
          
          // 取得输入焦点
          term.focus();
          
          // 连接websocket
          ws = new WebSocket("ws://"+window.location.host+"/api/k8s/terminal?namespace=" + namespace.text() + "&pods=" + pod.text() + "&container=" + container.text() + "&shell=" + shell.text());
          
          ws.onopen = function(event) {
            console.log(event)
            console.log("onopen")
          }
          ws.onclose = function(event) {
              console.log(event)
              console.log("onclose")
          }
          ws.onmessage = function(event) {
              console.log(event)
              // 服务端ssh输出, 写到web shell展示
              term.write(event.data)
          }
          ws.onerror = function(event) {
              console.log(event)
              console.log("onerror")
          }
          
          // 当浏览器窗口变化时, 重新适配终端
          window.addEventListener("resize", function () {
              term.fit()
          
              // 把web终端的尺寸term.rows和term.cols发给服务端, 通知sshd调整输出宽度
              var msg = {type: "resize", rows: term.rows, cols: term.cols}
              ws.send(JSON.stringify(msg))
          
              // console.log(term.rows + "," + term.cols)
          })
          
          // 当向web终端敲入字符时候的回调
          term.on('data', function(input) {
              // 写给服务端, 由服务端发给container
              var msg = {type: "input", input: input}
              ws.send(JSON.stringify(msg))
          })
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
            //console.log("");
			connTerminal();
		});
    </script>
</body>
`

func TerminalHtml(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, terminalHtml)
}
