(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-4a1c06fe"],{2134:function(t,e,a){},"7b90":function(t,e,a){"use strict";a("2134")},"91b6":function(t,e,a){"use strict";a.d(e,"a",(function(){return s}));var n=a("ead3");function s(t,e,a){return Object(n["b"])("/k8s/upload",t,e,a)}},c928:function(t,e,a){"use strict";a.r(e);var n=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("el-card",{attrs:{shadow:"never"}},[a("div",[a("el-select",{staticStyle:{width:"100%"},attrs:{filterable:"",placeholder:t.$t("please_select_namespace")},on:{change:t.selectedNamespace},nativeOn:{click:function(e){return t.getNamespace(e)}},model:{value:t.namespace,callback:function(e){t.namespace=e},expression:"namespace"}},t._l(t.namespaces,(function(t){return a("el-option",{key:t,attrs:{label:t,value:t}})})),1)],1),a("div",{staticStyle:{"margin-top":"15px"}},[a("el-select",{staticStyle:{width:"100%"},attrs:{placeholder:t.$t("please_select_pod"),filterable:""},on:{change:t.selectedPod},model:{value:t.pod,callback:function(e){t.pod=e},expression:"pod"}},t._l(t.pods,(function(t){return a("el-option",{key:t.value,attrs:{label:t.label,value:t.value}})})),1)],1),a("div",{staticStyle:{"margin-top":"15px"}},[a("el-select",{staticStyle:{width:"100%"},attrs:{placeholder:t.$t("please_select_container"),multiple:"",filterable:""},model:{value:t.container,callback:function(e){t.container=e},expression:"container"}},[a("el-option",{attrs:{label:t.$t("check_all"),value:"all"}}),t._l(t.containers,(function(t){return a("el-option",{key:t.value,attrs:{label:t.label,value:t.value}})}))],2)],1),a("div",{staticStyle:{"margin-top":"15px"}},[a("el-input",{staticStyle:{width:"100%",height:"40px"},attrs:{autocomplete:"off",placeholder:t.$t("please_input_dest_path")},model:{value:t.destPath,callback:function(e){t.destPath=e},expression:"destPath"}})],1),a("div",{staticStyle:{"margin-top":"15px"}},[a("el-dropdown",{staticClass:"avatar-container",staticStyle:{height:"36px",float:"right","margin-bottom":"10px"},attrs:{type:"success",trigger:"click"}},[a("div",{staticClass:"avatar-wrapper"},[a("el-button",{staticClass:"el-icon-upload",staticStyle:{width:"90px",height:"30px","margin-right":"6px","padding-top":"7px","padding-left":"14px"},attrs:{type:"success",round:"",size:"medium"}},[t._v(" "+t._s(t.$t("upload"))+" "),a("i",{staticClass:"el-icon-caret-bottom"})])],1),a("el-dropdown-menu",{attrs:{slot:"dropdown"},slot:"dropdown"},[a("el-dropdown-item",[a("span",{staticClass:"fake-file-btn"},[t._v(" "+t._s(t.$t("upload_file"))+" "),a("input",{staticStyle:{display:"block"},attrs:{type:"file",name:"files",multiple:"true"},on:{change:function(e){return t.uploadFileOrDir(e)}}})])]),a("el-dropdown-item",{attrs:{divided:""}},[a("span",{staticClass:"fake-file-btn"},[t._v(" "+t._s(t.$t("upload_dir"))+" "),a("input",{staticStyle:{display:"block"},attrs:{type:"file",name:"files",webkitdirectory:"",mozdirectory:"",accept:"*/*"},on:{change:function(e){return t.uploadFileOrDir(e)}}})])])],1)],1)],1)])],1)},s=[],i=a("1764"),l=a("f492"),o=a("91b6"),c={data(){return{namespace:"",pod:"",container:[],namespaces:[],pods:[],containers:[],destPath:"",fileList:[]}},methods:{getNamespace(){Object(i["a"])().then(t=>{if(t){this.namespaces=[],this.pod="",this.pods=[],this.container=[],this.containers=[];const e=t.items;for(const t in e)this.namespaces.push(e[t].metadata.name);console.log(this.namespaces)}})},selectedNamespace(){Object(l["a"])({namespace:this.namespace}).then(t=>{if(t){const e=[];this.pod="",this.pods=[],this.container=[],this.containers=[];const a=t.items;for(const t in a){const n={label:a[t].metadata.name,value:a[t].metadata.name};e.push(n)}this.pods=e,console.log(this.pods)}})},selectedPod(){this.container=[],this.containers=[],Object(l["a"])({namespace:this.namespace,pod:this.pod}).then(t=>{if(t){console.log(t);const e=[];this.container=[],this.containers=[];const a=t.spec.containers;for(const t in a){const n={label:a[t].name,value:a[t].name};e.push(n)}this.containers=e,console.log(this.containers)}})},uploadFileOrDir(t){const e=t.target.files;if(0===e.length)return void(t.target.value="");let a=this.container;"all"===a[0]&&(a=this.containers);const n=new FormData;for(let s=0;s<e.length;s++)n.append("files",e[s]);Object(o["a"])(n,{namespace:this.namespace,pod_name:this.pod,container_name:a,dest_path:this.destPath},{"Content-Type":"multipart/form-data"}).then(t=>{void 0!==t.failure?alert(t.failure):alert(t.success)},t=>{alert(t.info.message)}),t.target.value=""}},watch:{container:function(t,e){let a=t.indexOf("all"),n=e.indexOf("all");-1!==a&&-1===n&&t.length>1?this.container=["all"]:-1!==a&&-1!==n&&t.length>1&&this.container.splice(t.indexOf("all"),1)}}},r=c,p=(a("7b90"),a("2877")),d=Object(p["a"])(r,n,s,!1,null,null,null);e["default"]=d.exports},f492:function(t,e,a){"use strict";a.d(e,"a",(function(){return s}));var n=a("ead3");function s(t){return Object(n["a"])("/k8s/pods",t)}}}]);
//# sourceMappingURL=chunk-4a1c06fe.144aac8d.js.map