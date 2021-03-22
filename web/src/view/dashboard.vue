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
        <el-select :placeholder="$t('please_select_deployment')" multiple v-model="deployment" @change="selecteddeployment">
          <el-option :label="$t('check_all')" value="all"></el-option>
          <el-option v-for="item in deployments" :label="item.label" :value="item.value" :key="item.value"></el-option>
        </el-select>
        &nbsp;&nbsp;
        <el-button @click.native="getStatus">{{ $t('enter') }}</el-button>
        <br>
        <el-table
                class="app-table"
                size="medium"
                :data="tableData">
            <el-table-column prop="Pods" :label="$t('name')"></el-table-column>
            <el-table-column prop="Container" :label="$t('container')"></el-table-column>
            <el-table-column prop="Image" :label="$t('image')"></el-table-column>
            <el-table-column prop="Tag" :label="$t('tag')"></el-table-column>
            <el-table-column prop="ImagePullSecrets" :label="$t('image_pull_secrets')"></el-table-column>
            <el-table-column prop="State" :label="$t('state')"></el-table-column>
            <el-table-column prop="CPU" :label="$t('cpu')"></el-table-column>
            <el-table-column prop="RAM" :label="$t('ram')"></el-table-column>
            <el-table-column prop="OS" :label="$t('os')"></el-table-column>
            <el-table-column prop="Operate" :label="$t('operate')" width="160" align="center">
              <template slot-scope="scope">
                <el-button v-if="scope.row.OS==='unix'" @click.native="openTerminal(scope.row, 'sh')">{{ $t('sh') }}</el-button>
                <el-button v-if="scope.row.OS==='windows'" @click.native="openTerminal(scope.row, 'cmd')">{{ $t('cmd') }}</el-button>
                <span></span>
                <el-button @click.native="openFileBrowser(scope.row, '/')">{{ $t('file_browser') }}</el-button>
              </template>
          </el-table-column>
        </el-table>
    </el-card>
    <el-dialog
      :title="$t('terminal')"
      :visible.sync="dialogTerminalVisible"
      @close="dialogTerminalVisible = false"
      :before-close="handleClose">
    <div>
      <span>前端实现中</span>
      <div id="terminal-container"></div>
    </div>
    </el-dialog>
    <el-dialog
        :title="$t('file_browser')"
        :visible.sync="dialogFileBrowserVisible"
        @close="dialogFileBrowserVisible = false"
        :before-close="handleClose">
      <span class="el-icon-folder-opened"></span>
      <span class="el-icon-refresh"></span>
        <el-table
            class="app-table"
            size="max"
            :data="fileBrowserData">
          <el-header></el-header>
          <el-table-column prop="Name" :label="$t('name')">
            <template slot-scope="scope">
              <span class="el-icon-folder" v-if="scope.row.IsDir" @click="openFileBrowser('', scope.row.Path+scope.row.Name)">{{scope.row.Name}}</span>
              <span class="el-icon-files" v-else>{{scope.row.Name}}</span>
            </template>
          </el-table-column>
          <el-table-column prop="Size" :label="$t('size')"></el-table-column>
          <el-table-column prop="Mode" :label="$t('mode')"></el-table-column>
          <el-table-column prop="ModTime" :label="$t('mod_time')"></el-table-column>
          <el-table-column prop="Download" :label="$t('operate')">
            <template slot-scope="scope">
              <span class="el-icon-download" @click="download(scope.row.Path+scope.row.Name)"></span>
            </template>
          </el-table-column>
        </el-table>
    </el-dialog>
  </div>
</template>


<script>
import { GetStatus } from '../api/status'
import { GetNamespace } from '../api/namespaces'
import { GetDeployment } from "../api/deployment";
import {FileBrowser} from "../api/filebrowser";

export default {
  data() {
    return {
      namespace: "",
      deployment: [],
      namespaces:[],
      tableLoading: false,
      tableData: [],
      deployments: [],
      fileBrowserData: [],
      pods: "",
      container:"",
      path: "",
      oldPath: '',
      dialogTerminalVisible: false,
      dialogFileBrowserVisible: false,
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res =>{
        if (res) {
          this.namespaces=[]
          this.deployment = []
          this.deployments = []
          this.tableData = []
          const data = res.items
          for(var key in data){
            this.namespaces.push(data[key].metadata.name)
          }
          console.log(this.namespaces)
        }
      })
    },
    selectedNamespace(options) {
      console.log(options);
      GetDeployment({namespace: options}).then(res => {
        if (res) {
          const deployments=[]
          this.deployments = []
          this.tableData = []
          const data = res.items
          for(const key in data){
            const _d = {label:data[key].metadata.name, value:data[key].metadata.name}
            deployments.push(_d)
          }
          this.deployments = deployments
          console.log(this.deployments)
        }
      })
    },
    selecteddeployment(options) {
      console.log(options);
    },
    getStatus() {
      console.log(this.deployment[0], this.namespace);
      let deployment = this.deployment
      if (this.deployment[0] === "all") {
        deployment = []
        for (const key in this.deployments) {
          deployment.push(this.deployments[key].value)
        }
      }
      GetStatus({namespace: this.namespace, deployment:deployment}).then(res => {
        if (res) {
          console.log(res)
          this.tableData = []
          for (const i in res) {
            const pod_name =res[i].pod_name;
            const cData = res[i].containers;
            for (const j in cData) {
              let tr = {
                Pods:"",
                Container:cData[j].name,
                Image:cData[j].image,
                Tag:cData[j].version,
                ImagePullSecrets:cData[j].image_pull_policy,
                State:cData[j].state,
                CPU:cData[j].cpu,
                RAM:cData[j].ram,
                OS: cData[j].os,
              }
              if (j === "0") {
                tr = {
                  Pods:pod_name,
                  Container:cData[j].name,
                  Image:cData[j].image,
                  Tag:cData[j].version,
                  ImagePullSecrets:cData[j].image_pull_policy,
                  State:cData[j].state,
                  CPU:cData[j].cpu,
                  RAM:cData[j].ram,
                  OS: cData[j].os,
                }
              }
              this.tableData.push(tr)
            }
          }
        }
      })
    },
    openTerminal(options, shell) {
      this.dialogTerminalVisible = true
      console.log(options, this.namespace,shell)
    },
    openFileBrowser(options, path) {
      this.dialogFileBrowserVisible = true
      console.log(options, path)
      if (path === undefined) {
        path = "/"
      }
      if (path === "/") {
        this.pods = options.Pods
        this.container = options.Container
      }
      this.oldPath = this.path
      this.path = path
      this.fileBrowserData = []
      FileBrowser({
        namespace: this.namespace,
        pods: this.pods,
        container: this.container,
        path: path,
      }).then(res => {
        console.log(res)
        this.fileBrowserData = []
        if (res !== undefined) {
          this.fileBrowserData = res
        }
      }, err => {
        console.log(err)
        alert(err.info.message)
      })
    },
    download(path) {
      console.log(this.namespace)
      console.log(this.pods)
      console.log(this.container)
      console.log(path)
    },
    handleClose(done) {
      this.$confirm('确认关闭？')
          .then(_ => {
            done();
          })
          .catch(_ => {});
    }
  },
  watch:{
    deployment:function(val,oldVal){
      let index =  val.indexOf('all'),oldIndex =  oldVal.indexOf('all');
      if(index!==-1 && oldIndex===-1 && val.length>1)
        this.deployment=['all'];
      else if(index!==-1 && oldIndex!==-1 && val.length>1)
        this.deployment.splice(val.indexOf('all'),1)
    }
  },
}
</script>