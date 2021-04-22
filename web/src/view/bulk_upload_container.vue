<template>
  <div>
    <el-card shadow="never">
      <div>
        <el-select style="width: 100%" v-model="namespace" @click.native="getNamespace" @change="selectedNamespace" filterable :placeholder="$t('please_select_namespace')">
          <el-option
              v-for="item in namespaces"
              :label="item"
              :value="item"
              :key="item"
          ></el-option>
        </el-select>
      </div>
      <div style="margin-top: 15px">
      <el-select style="width: 100%" :placeholder="$t('please_select_pod')" v-model="pod" filterable @change="selectedPod">
        <el-option v-for="item in pods" :label="item.label" :value="item.value" :key="item.value"></el-option>
      </el-select>
      </div>
      <div style="margin-top: 15px">
      <el-select style="width: 100%" :placeholder="$t('please_select_container')" multiple filterable v-model="container">
        <el-option :label="$t('check_all')" value="all"></el-option>
        <el-option v-for="item in containers" :label="item.label" :value="item.value" :key="item.value"></el-option>
      </el-select>
      </div>
      <div style="margin-top: 15px">
      <el-input v-model="destPath" style="width: 100%;height: 40px" autocomplete="off" :placeholder="$t('please_input_dest_path')"></el-input>
      </div>
      <div style="margin-top: 15px">
        <el-dropdown  type="success" class="avatar-container" trigger="click" style="height: 36px;float: right;margin-bottom: 10px;">
          <div class="avatar-wrapper">
            <el-button style="width: 90px; height: 30px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-upload" size="medium">
              {{ $t('upload') }}
              <i class="el-icon-caret-bottom" />
            </el-button>
          </div>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item>
                <span class="fake-file-btn">
                  {{ $t('upload_file') }}
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event)" name="files" multiple="true">
                </span>
            </el-dropdown-item>
            <el-dropdown-item divided>
                <span class="fake-file-btn">
                  {{ $t('upload_dir') }}
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event)" name="files" webkitdirectory mozdirectory accept="*/*">
                </span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </el-card>
  </div>
</template>

<style>
.fake-file-btn {
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

<script>
import { GetNamespace } from '../api/namespaces'
import {GetPods} from "../api/pods"
import {FileOrDirUpload} from "../api/upload"

export default {
  data() {
    return {
      namespace: "",
      pod: "",
      container: [],
      namespaces:[],
      pods: [],
      containers:[],
      destPath:"",
      fileList:[],
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res => {
        if (res) {
          this.namespaces = []
          this.pod = ""
          this.pods = []
          this.container = []
          this.containers = []
          const data = res.items;
          for (const key in data) {
            this.namespaces.push(data[key].metadata.name)
          }
          console.log(this.namespaces)
        }
      })
    },
    selectedNamespace() {
      GetPods({namespace: this.namespace}).then(res => {
        if (res) {
          const pods = [];
          this.pod = ""
          this.pods = []
          this.container = []
          this.containers = []
          const data = res.items;
          for (const key in data) {
            const _d = {label: data[key].metadata.name, value: data[key].metadata.name}
            pods.push(_d)
          }
          this.pods = pods
          console.log(this.pods)
        }
      })
    },
    selectedPod() {
      this.container = []
      this.containers = []
      GetPods({namespace: this.namespace, pod: this.pod}).then(res => {
        if (res) {
          console.log(res)
          const containers = []
          this.container = []
          this.containers = []
          const data = res.spec.containers
          for (const key in data) {
            const _d = {label: data[key].name, value: data[key].name}
            containers.push(_d)
          }
          this.containers = containers
          console.log(this.containers)
        }
      })
    },
    uploadFileOrDir(e) {
      const files = e.target.files;
      if (files.length === 0 ) {
        e.target.value = ""
        return
      }
      let container = this.container
      if (container[0] === "all") {
        container = this.containers
      }
      const formData = new FormData();
      //追加文件数据
      for (let i = 0; i < files.length; i++) {
        formData.append("files", files[i]);
      }
      FileOrDirUpload(formData, {
        namespace:this.namespace,
        pod_name:this.pod,
        container_name:container,
        dest_path:this.destPath},{"Content-Type":"multipart/form-data"}).then((res) => {
        if (res.failure !== undefined) {
          alert(res.failure)
        }else {
          alert(res.success)
        }
      }, (err) => {
        alert(err.info.message)
      })
      e.target.value = ""
    },
  },
  watch:{
    container:function(val,oldVal){
      let index =  val.indexOf('all'),oldIndex =  oldVal.indexOf('all');
      if(index!==-1 && oldIndex===-1 && val.length>1)
        this.container=['all'];
      else if(index!==-1 && oldIndex!==-1 && val.length>1)
        this.container.splice(val.indexOf('all'),1)
    }
  },
}
</script>