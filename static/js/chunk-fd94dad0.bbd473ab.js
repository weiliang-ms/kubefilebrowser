(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-fd94dad0"],{"46c9":function(e,t,a){"use strict";a.r(t);var n=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",[a("el-card",{attrs:{shadow:"never"}},[a("div",[a("el-select",{attrs:{filterable:"",placeholder:e.$t("please_select_namespace")},on:{change:e.selectedNamespace},nativeOn:{click:function(t){return e.getNamespace(t)}},model:{value:e.namespace,callback:function(t){e.namespace=t},expression:"namespace"}},e._l(e.namespaces,(function(e){return a("el-option",{key:e,attrs:{label:e,value:e}})})),1),e._v(" "),a("el-select",{attrs:{placeholder:e.$t("please_select_pod")},on:{change:e.selectedPod},model:{value:e.pod,callback:function(t){e.pod=t},expression:"pod"}},e._l(e.pods,(function(e){return a("el-option",{key:e.value,attrs:{label:e.label,value:e.value}})})),1),e._v(" "),a("el-select",{attrs:{placeholder:e.$t("please_select_container"),multiple:""},model:{value:e.container,callback:function(t){e.container=t},expression:"container"}},[a("el-option",{attrs:{label:e.$t("check_all"),value:"all"}}),e._l(e.containers,(function(e){return a("el-option",{key:e.value,attrs:{label:e.label,value:e.value}})}))],2),e._v(" "),a("el-input",{staticStyle:{width:"193px",height:"40px"},attrs:{autocomplete:"off",placeholder:e.$t("please_input_dest_path")},model:{value:e.destPath,callback:function(t){e.destPath=t},expression:"destPath"}})],1),e._v(" "),a("div",[a("el-upload",{staticStyle:{display:"inline-block"},attrs:{action:"",multiple:"","file-list":e.fileList,"on-change":e.fileChange,"on-remove":e.fileRemove,"auto-upload":!1}},[a("el-button",{attrs:{type:"primary",plain:""}},[a("i",{staticClass:"el-icon-upload el-icon--right"}),e._v(e._s(e.$t("select_file")))])],1),a("el-button",{staticStyle:{"margin-left":"10px","vertical-align":"top"},attrs:{type:"success",plain:""},on:{click:e.submitUpload}},[e._v(e._s(e.$t("upload_all")))])],1)])],1)},s=[],l=a("1764"),i=a("f492"),o=a("ead3");function c(e,t,a){return Object(o["b"])("/k8s/upload",e,t,a)}var r={data(){return{namespace:"",pod:"",container:[],namespaces:[],pods:[],containers:[],destPath:"",fileList:[]}},methods:{getNamespace(){Object(l["a"])().then(e=>{if(e){this.namespaces=[],this.pod="",this.pods=[],this.container=[],this.containers=[];const t=e.items;for(const e in t)this.namespaces.push(t[e].metadata.name);console.log(this.namespaces)}})},selectedNamespace(){Object(i["a"])({namespace:this.namespace}).then(e=>{if(e){const t=[];this.pod="",this.pods=[],this.container=[],this.containers=[];const a=e.items;for(const e in a){const n={label:a[e].metadata.name,value:a[e].metadata.name};t.push(n)}this.pods=t,console.log(this.pods)}})},selectedPod(){Object(i["a"])({namespace:this.namespace,pod:this.pod}).then(e=>{if(e){console.log(e);const t=[];this.container=[],this.containers=[];const a=e.spec.containers;for(const e in a){const n={label:a[e].name,value:a[e].name};t.push(n)}this.containers=t,console.log(this.containers)}})},fileChange(e,t){this.fileList=t},fileRemove(e,t){this.fileList=t},submitUpload(){console.log("this.fileList",this.fileList);let e=new FormData;this.fileList.forEach(t=>{e.append("files",t.raw)});var t=this.container[0];if("all"===this.container[0]){t=[];for(const e in this.containers)t.push(this.containers[e].value)}c(e,{namespace:this.namespace,pod_name:this.pod,container_name:t,dest_path:this.destPath},{"Content-Type":"multipart/form-data"}).then(e=>{void 0!==e.failure?alert(e.failure):alert(e.success)},e=>{alert(e.info.message)})}},watch:{container:function(e,t){let a=e.indexOf("all"),n=t.indexOf("all");-1!==a&&-1===n&&e.length>1?this.container=["all"]:-1!==a&&-1!==n&&e.length>1&&this.container.splice(e.indexOf("all"),1)}}},p=r,d=a("2877"),h=Object(d["a"])(p,n,s,!1,null,null,null);t["default"]=h.exports},f492:function(e,t,a){"use strict";a.d(t,"a",(function(){return s}));var n=a("ead3");function s(e){return Object(n["a"])("/k8s/pods",e)}}}]);
//# sourceMappingURL=chunk-fd94dad0.bbd473ab.js.map