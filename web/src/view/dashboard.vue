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
                <el-button @click.native="openFilebrowser(scope.row)">{{ $t('filebrowser') }}</el-button>
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
      <span>未实现</span>
      <div id="terminal-container"></div>
    </div>
    </el-dialog>
    <el-dialog
        :title="$t('filebrowser')"
        :visible.sync="dialogFileBrowserVisible"
        @close="dialogFileBrowserVisible = false"
        :before-close="handleClose">
      <div>
        <span>未实现</span>
        <div id="filebrowser-container"></div>
      </div>
    </el-dialog>
  </div>
</template>


<script>
import { GetStatus } from '@/api/status'
import { GetNamespace } from '@/api/namespaces'
import { GetDeployment } from "@/api/deployment";

export default {
  data() {
    return {
      namespace: "",
      deployment: [],
      namespaces:[],
      tableLoading: false,
      tableData: [],
      deployments: [],
      dialogTerminalVisible: false,
      dialogFileBrowserVisible: false,
      ws: null,
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
          var data = res.items
          for(var key in data){
            this.namespaces.push(data[key].metadata.name)
          };
          console.log(this.namespaces)
        }
      })
    },
    selectedNamespace(options) {
      console.log(options);
      GetDeployment({namespace: options}).then(res => {
        if (res) {
          var deployments=[]
          this.deployments = []
          this.tableData = []
          var data = res.items
          for(var key in data){
            var _d = {label:data[key].metadata.name, value:data[key].metadata.name}
            deployments.push(_d)
          };
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
      var deployment = this.deployment[0]
      if (this.deployment[0] === "all") {
        deployment = []
        for (var key in this.deployments) {
          deployment.push(this.deployments[key].value)
        }
      }
      GetStatus({namespace: this.namespace, deployment:deployment}).then(res => {
        if (res) {
          console.log(res)
          this.tableData = []
          for (var i in res) {
            var pod_name =res[i].pod_name;
            var cData = res[i].containers;
            for (var j in cData) {
              var tr = {
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
    openFilebrowser(options) {
      this.dialogFileBrowserVisible = true
      console.log(options)
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