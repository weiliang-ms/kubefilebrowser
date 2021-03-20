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
      <el-input v-model="destPath" style="width: 217px;height: 40px" autocomplete="off" :placeholder="$t('please_input_dest_path')"></el-input>
      &nbsp;&nbsp;
      <el-button @click.native="downloadFile">{{ $t('download_file') }}</el-button>
    </el-card>
    <div id="terminal-container" style="width: 100%;height: 100%"></div>
  </div>
</template>

<script>
import { GetNamespace } from '@/api/namespaces'
import {GetPods} from "@/api/pods";

export default {
  data() {
    return {
      namespace: "",
      pod: "",
      container: "",
      namespaces:[],
      pods: [],
      containers:[],
      destPath:""
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res => {
        if (res) {
          this.namespaces = []
          this.pod = ""
          this.pods = []
          this.container = ""
          this.containers = []
          this.shell = ""
          var data = res.items
          for (var key in data) {
            this.namespaces.push(data[key].metadata.name)
          }
          ;
          console.log(this.namespaces)
        }
      })
    },
    selectedNamespace() {
      GetPods({namespace: this.namespace}).then(res => {
        if (res) {
          var pods = []
          this.pod = ""
          this.pods = []
          this.container = ""
          this.containers = []
          var data = res.items
          for (var key in data) {
            var _d = {label: data[key].metadata.name, value: data[key].metadata.name}
            pods.push(_d)
          }
          ;
          this.pods = pods
          console.log(this.pods)
        }
      })
    },
    selectedPod() {
      GetPods({namespace: this.namespace, pod: this.pod}).then(res => {
        if (res) {
          console.log(res)
          var containers = []
          this.container = ""
          this.containers = []
          var data = res.spec.containers
          for (var key in data) {
            var _d = {label: data[key].name, value: data[key].name}
            containers.push(_d)
          }
          ;
          this.containers = containers
          console.log(this.containers)
        }
      })
    },
    downloadFile() {

    }
  },
}
</script>