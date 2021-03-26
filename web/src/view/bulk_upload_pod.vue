<template>
  <div>
    <el-card shadow="never">
      <div>
        <el-select v-model="namespace" @click.native="getNamespace" @change="selectedNamespace" filterable :placeholder="$t('please_select_namespace')">
          <el-option
              v-for="item in namespaces"
              :label="item"
              :value="item"
              :key="item"
          ></el-option>
        </el-select>
        &nbsp;&nbsp;
        <el-select :placeholder="$t('please_select_pod')" filterable multiple v-model="pod">
          <el-option :label="$t('check_all')" value="all"></el-option>
          <el-option v-for="item in pods" :label="item.label" :value="item.value" :key="item.value"></el-option>
        </el-select>
        &nbsp;&nbsp;
        <el-input v-model="destPath" style="width: 193px;height: 40px" autocomplete="off" :placeholder="$t('please_input_dest_path')"></el-input>
      </div>
      &nbsp;&nbsp;
      <div>
        <el-upload
            action=""
            multiple
            :file-list="fileList"
            :on-change="fileChange"
            :on-remove="fileRemove"
            :auto-upload="false"
            style="display: inline-block">
          <el-button type="primary" plain><i class="el-icon-upload el-icon--right"></i>{{$t('select_file')}}</el-button>
        </el-upload>
        <!--        <el-upload-->
        <!--            action=""-->
        <!--            webkitdirector-->
        <!--            :file-list="fileList"-->
        <!--            :on-change="fileChange"-->
        <!--            :on-remove="fileRemove"-->
        <!--            :auto-upload="false"-->
        <!--            style="display: inline-block">-->
        <!--          <el-button type="primary" plain><i class="el-icon-upload el-icon&#45;&#45;right"></i>{{$t('select_dir')}}</el-button>-->
        <!--        </el-upload>-->
        <el-button style="margin-left: 10px;vertical-align: top;" type="success" plain @click="submitUpload">{{$t('upload_all')}}</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { GetNamespace } from '../api/namespaces'
import {GetPods} from "../api/pods"
import {MultiUpload} from "../api/multiupload"

export default {
  data() {
    return {
      namespace: "",
      pod: [],
      namespaces:[],
      pods: [],
      destPath:"",
      fileList:[],
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res => {
        if (res) {
          this.namespaces = []
          this.pod = []
          this.pods = []
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
          this.pod = []
          this.pods = []
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
    fileChange(file, fileList) {
      this.fileList = fileList
    },
    fileRemove(file, fileList) {
      this.fileList = fileList
    },
    submitUpload(){
      console.log("this.fileList", this.fileList)
      let formData = new FormData
      this.fileList.forEach(file => {
        formData.append("files", file.raw)
      })
      var pod = this.pod[0]
      if (this.pod[0] === "all") {
        pod = ""
        for (const key in this.pods) {
          pod.push(this.pods[key].value)
        }
      }else {
        for (const  key in this.pod) {
          pod.push(this.pod[key].value)
        }
      }
      MultiUpload(formData, {
        namespace:this.namespace,
        pod_name:pod,
        dest_path:this.destPath
      },{"Content-Type":"multipart/form-data"}).then(res => {
        if (res.failure !== undefined) {
          alert(res.failure)
        }else {
          alert(res.success)
        }
      }, (err) =>{
        alert(err.info.message)
      })
    },
  },
  watch:{
    pod:function(val,oldVal){
      let index =  val.indexOf('all'),oldIndex =  oldVal.indexOf('all');
      if(index!==-1 && oldIndex===-1 && val.length>1)
        this.pod=['all'];
      else if(index!==-1 && oldIndex!==-1 && val.length>1)
        this.pod.splice(val.indexOf('all'),1)
    }
  },
}
</script>