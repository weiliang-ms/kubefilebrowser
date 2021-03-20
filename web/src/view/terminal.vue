<template>
  <div>
    <el-card shadow="never">
      <el-select v-model="namespace" @click.native="getNamespace" @change="selectedNamespace" filterable :placeholder="$t('please_select_namespace')">
        <el-option
            v-for="item in namespaces"
            :label="item"
            :value="item"
            :key="item"
        ></el-option>
      </el-select>
      &nbsp;&nbsp;
      <el-select :placeholder="$t('please_select_pod')" v-model="pod" @change="selectedPod">
        <el-option v-for="item in pods" :label="item.label" :value="item.value" :key="item.value"></el-option>
      </el-select>
      &nbsp;&nbsp;
      <el-select :placeholder="$t('please_select_container')" v-model="container">
        <el-option v-for="item in containers" :label="item.label" :value="item.value" :key="item.value"></el-option>
      </el-select>
      &nbsp;&nbsp;
      <el-select :placeholder="$t('please_select_shell')" v-model="shell">
        <el-option v-for="item in shells" :label="item.label" :value="item.value" :key="item.value"></el-option>
      </el-select>
      &nbsp;&nbsp;
      <el-button @click.native="openTerminal">{{ $t('enter') }}</el-button>
    </el-card>
    <div id="terminal-container" style="width: 100%;height: 100%"></div>
  </div>
</template>


<script>
import { GetNamespace } from '@/api/namespaces'
import {GetPods} from "@/api/pods";
import { Terminal } from 'xterm';
import * as fit from 'xterm/lib/addons/fit/fit'
import * as attach from 'xterm/lib/addons/attach/attach'
import * as winptyCompat from 'xterm/lib/addons/winptyCompat/winptyCompat'
import * as webLinks from 'xterm/lib/addons/webLinks/webLinks'

export default {
  data() {
    return {
      namespace: "",
      pod: "",
      container: "",
      shell: "",
      namespaces:[],
      pods: [],
      containers:[],
      order:"",
      shells : [
        {
          label: "sh",
          value: "sh"
        },
        {
          label: "bash",
          value: "bash"
        },
        {
          label: "cmd",
          value: "cmd"
        }
      ]
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res =>{
        if (res) {
          this.namespaces=[]
          this.pod = ""
          this.pods = []
          this.container = ""
          this.containers = []
          this.shell = ""
          const data = res.items
          for(const key in data){
            this.namespaces.push(data[key].metadata.name)
          };
          console.log(this.namespaces)
        }
      })
    },
    selectedNamespace() {
      GetPods({namespace: this.namespace}).then(res => {
        if (res) {
          const pods=[]
          this.pod = ""
          this.pods = []
          this.container = ""
          this.containers = []
          const data = res.items
          for(const key in data){
            const _d = {label:data[key].metadata.name, value:data[key].metadata.name}
            pods.push(_d)
          };
          this.pods = pods
          console.log(this.pods)
        }
      })
    },
    selectedPod() {
      GetPods({namespace: this.namespace, pod: this.pod}).then(res => {
        if (res) {
          console.log(res)
          const containers=[]
          this.container = ""
          this.containers = []
          const data = res.spec.containers
          for(const key in data){
            const _d = {label:data[key].name, value:data[key].name}
            containers.push(_d)
          };
          this.containers = containers
          console.log(this.containers)
        }
      })
    },
    openTerminal() {
      console.log(this.ws)
      document.getElementById("terminal-container").innerHTML="";
      document.getElementById('terminal-container').style.height = window.innerHeight + 'px';
      document.getElementById('terminal-container').style.width = window.innerWidth + 'px';
      const url = "ws://"+window.location.host+"/api/k8s/terminal?namespace="+this.namespace+"&pods="+this.pod+"&container="+this.container+"&shell="+this.shell;

      // xterm配置自适应大小插件
      Terminal.applyAddon(fit);
      Terminal.applyAddon(attach)

      // 这俩插件不知道干嘛的, 用总比不用好
      Terminal.applyAddon(winptyCompat)
      Terminal.applyAddon(webLinks)

      // 创建终端
      const term = new Terminal({
        cursorBlink: true
      });
      term.open(document.getElementById("terminal-container"));
      // 使用fit插件自适应terminal size
      term.fit();
      term.winptyCompatInit()
      term.webLinksInit()

      // 取得输入焦点
      term.focus();

      // 连接websocket
      const ws = new WebSocket(url);

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
        const msg = {type: "resize", rows: term.rows, cols: term.cols}
        ws.send(JSON.stringify(msg))

        // console.log(term.rows + "," + term.cols)
      })

      // 当向web终端敲入字符时候的回调
      term.on('data', function(input) {
        // 写给服务端, 由服务端发给container
        const msg = {type: "input", input: input}
        ws.send(JSON.stringify(msg))
      })
      // 向web终端敲入回车时候的回调

    }
  }
}
</script>