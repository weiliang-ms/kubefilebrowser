(function(e){function t(t){for(var n,s,c=t[0],d=t[1],l=t[2],u=0,i=[];u<c.length;u++)s=c[u],Object.prototype.hasOwnProperty.call(o,s)&&o[s]&&i.push(o[s][0]),o[s]=0;for(n in d)Object.prototype.hasOwnProperty.call(d,n)&&(e[n]=d[n]);f&&f(t);while(i.length)i.shift()();return r.push.apply(r,l||[]),a()}function a(){for(var e,t=0;t<r.length;t++){for(var a=r[t],n=!0,s=1;s<a.length;s++){var c=a[s];0!==o[c]&&(n=!1)}n&&(r.splice(t--,1),e=d(d.s=a[0]))}return e}var n={},s={index:0},o={index:0},r=[];function c(e){return d.p+"js/"+({}[e]||e)+"."+{"chunk-119f7568":"b77e3f49","chunk-3c5be004":"93ca7cc7","chunk-100baf8e":"e0672207","chunk-2d21f4b8":"4db647b5","chunk-74d03424":"ef89a284","chunk-d6e2ec14":"8e11d2b0","chunk-fd94dad0":"d0d2451f"}[e]+".js"}function d(t){if(n[t])return n[t].exports;var a=n[t]={i:t,l:!1,exports:{}};return e[t].call(a.exports,a,a.exports,d),a.l=!0,a.exports}d.e=function(e){var t=[],a={"chunk-119f7568":1};s[e]?t.push(s[e]):0!==s[e]&&a[e]&&t.push(s[e]=new Promise((function(t,a){for(var n="css/"+({}[e]||e)+"."+{"chunk-119f7568":"aca39d08","chunk-3c5be004":"31d6cfe0","chunk-100baf8e":"31d6cfe0","chunk-2d21f4b8":"31d6cfe0","chunk-74d03424":"31d6cfe0","chunk-d6e2ec14":"31d6cfe0","chunk-fd94dad0":"31d6cfe0"}[e]+".css",o=d.p+n,r=document.getElementsByTagName("link"),c=0;c<r.length;c++){var l=r[c],u=l.getAttribute("data-href")||l.getAttribute("href");if("stylesheet"===l.rel&&(u===n||u===o))return t()}var i=document.getElementsByTagName("style");for(c=0;c<i.length;c++){l=i[c],u=l.getAttribute("data-href");if(u===n||u===o)return t()}var f=document.createElement("link");f.rel="stylesheet",f.type="text/css",f.onload=t,f.onerror=function(t){var n=t&&t.target&&t.target.src||o,r=new Error("Loading CSS chunk "+e+" failed.\n("+n+")");r.code="CSS_CHUNK_LOAD_FAILED",r.request=n,delete s[e],f.parentNode.removeChild(f),a(r)},f.href=o;var b=document.getElementsByTagName("head")[0];b.appendChild(f)})).then((function(){s[e]=0})));var n=o[e];if(0!==n)if(n)t.push(n[2]);else{var r=new Promise((function(t,a){n=o[e]=[t,a]}));t.push(n[2]=r);var l,u=document.createElement("script");u.charset="utf-8",u.timeout=120,d.nc&&u.setAttribute("nonce",d.nc),u.src=c(e);var i=new Error;l=function(t){u.onerror=u.onload=null,clearTimeout(f);var a=o[e];if(0!==a){if(a){var n=t&&("load"===t.type?"missing":t.type),s=t&&t.target&&t.target.src;i.message="Loading chunk "+e+" failed.\n("+n+": "+s+")",i.name="ChunkLoadError",i.type=n,i.request=s,a[1](i)}o[e]=void 0}};var f=setTimeout((function(){l({type:"timeout",target:u})}),12e4);u.onerror=u.onload=l,document.head.appendChild(u)}return Promise.all(t)},d.m=e,d.c=n,d.d=function(e,t,a){d.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:a})},d.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},d.t=function(e,t){if(1&t&&(e=d(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var a=Object.create(null);if(d.r(a),Object.defineProperty(a,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var n in e)d.d(a,n,function(t){return e[t]}.bind(null,n));return a},d.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return d.d(t,"a",t),t},d.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},d.p="/static/",d.oe=function(e){throw console.error(e),e};var l=window["webpackJsonp"]=window["webpackJsonp"]||[],u=l.push.bind(l);l.push=t,l=l.slice();for(var i=0;i<l.length;i++)t(l[i]);var f=u;r.push([0,"chunk-vendors"]),a()})({0:function(e,t,a){e.exports=a("56d7")},"034f":function(e,t,a){"use strict";a("85ec")},"0351":function(e,t,a){var n={"./dashboard.vue":["d8e3","chunk-3c5be004","chunk-2d21f4b8"],"./download.vue":["67a0","chunk-3c5be004","chunk-d6e2ec14"],"./layer.vue":["ad86","chunk-119f7568"],"./multiupload.vue":["5667","chunk-3c5be004","chunk-74d03424"],"./terminal.vue":["fe80","chunk-3c5be004","chunk-100baf8e"],"./upload.vue":["46c9","chunk-3c5be004","chunk-fd94dad0"]};function s(e){if(!a.o(n,e))return Promise.resolve().then((function(){var t=new Error("Cannot find module '"+e+"'");throw t.code="MODULE_NOT_FOUND",t}));var t=n[e],s=t[0];return Promise.all(t.slice(1).map(a.e)).then((function(){return a(s)}))}s.keys=function(){return Object.keys(n)},s.id="0351",e.exports=s},"3dfd":function(e,t,a){"use strict";var n=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"global-container"},[a("router-view")],1)},s=[],o=(a("034f"),a("2877")),r={},c=Object(o["a"])(r,n,s,!1,null,null,null);t["a"]=c.exports},"3fab":function(e,t,a){"use strict";t["a"]={}},4678:function(e,t,a){var n={"./af":"2bfb","./af.js":"2bfb","./ar":"8e73","./ar-dz":"a356","./ar-dz.js":"a356","./ar-kw":"423e","./ar-kw.js":"423e","./ar-ly":"1cfd","./ar-ly.js":"1cfd","./ar-ma":"0a84","./ar-ma.js":"0a84","./ar-sa":"8230","./ar-sa.js":"8230","./ar-tn":"6d83","./ar-tn.js":"6d83","./ar.js":"8e73","./az":"485c","./az.js":"485c","./be":"1fc1","./be.js":"1fc1","./bg":"84aa","./bg.js":"84aa","./bm":"a7fa","./bm.js":"a7fa","./bn":"9043","./bn-bd":"9686","./bn-bd.js":"9686","./bn.js":"9043","./bo":"d26a","./bo.js":"d26a","./br":"6887","./br.js":"6887","./bs":"2554","./bs.js":"2554","./ca":"d716","./ca.js":"d716","./cs":"3c0d","./cs.js":"3c0d","./cv":"03ec","./cv.js":"03ec","./cy":"9797","./cy.js":"9797","./da":"0f14","./da.js":"0f14","./de":"b469","./de-at":"b3eb","./de-at.js":"b3eb","./de-ch":"bb71","./de-ch.js":"bb71","./de.js":"b469","./dv":"598a","./dv.js":"598a","./el":"8d47","./el.js":"8d47","./en-au":"0e6b","./en-au.js":"0e6b","./en-ca":"3886","./en-ca.js":"3886","./en-gb":"39a6","./en-gb.js":"39a6","./en-ie":"e1d3","./en-ie.js":"e1d3","./en-il":"7333","./en-il.js":"7333","./en-in":"ec2e","./en-in.js":"ec2e","./en-nz":"6f50","./en-nz.js":"6f50","./en-sg":"b7e9","./en-sg.js":"b7e9","./eo":"65db","./eo.js":"65db","./es":"898b","./es-do":"0a3c","./es-do.js":"0a3c","./es-mx":"b5b7","./es-mx.js":"b5b7","./es-us":"55c9","./es-us.js":"55c9","./es.js":"898b","./et":"ec18","./et.js":"ec18","./eu":"0ff2","./eu.js":"0ff2","./fa":"8df4","./fa.js":"8df4","./fi":"81e9","./fi.js":"81e9","./fil":"d69a","./fil.js":"d69a","./fo":"0721","./fo.js":"0721","./fr":"9f26","./fr-ca":"d9f8","./fr-ca.js":"d9f8","./fr-ch":"0e49","./fr-ch.js":"0e49","./fr.js":"9f26","./fy":"7118","./fy.js":"7118","./ga":"5120","./ga.js":"5120","./gd":"f6b4","./gd.js":"f6b4","./gl":"8840","./gl.js":"8840","./gom-deva":"aaf2","./gom-deva.js":"aaf2","./gom-latn":"0caa","./gom-latn.js":"0caa","./gu":"e0c5","./gu.js":"e0c5","./he":"c7aa","./he.js":"c7aa","./hi":"dc4d","./hi.js":"dc4d","./hr":"4ba9","./hr.js":"4ba9","./hu":"5b14","./hu.js":"5b14","./hy-am":"d6b6","./hy-am.js":"d6b6","./id":"5038","./id.js":"5038","./is":"0558","./is.js":"0558","./it":"6e98","./it-ch":"6f12","./it-ch.js":"6f12","./it.js":"6e98","./ja":"079e","./ja.js":"079e","./jv":"b540","./jv.js":"b540","./ka":"201b","./ka.js":"201b","./kk":"6d79","./kk.js":"6d79","./km":"e81d","./km.js":"e81d","./kn":"3e92","./kn.js":"3e92","./ko":"22f8","./ko.js":"22f8","./ku":"2421","./ku.js":"2421","./ky":"9609","./ky.js":"9609","./lb":"440c","./lb.js":"440c","./lo":"b29d","./lo.js":"b29d","./lt":"26f9","./lt.js":"26f9","./lv":"b97c","./lv.js":"b97c","./me":"293c","./me.js":"293c","./mi":"688b","./mi.js":"688b","./mk":"6909","./mk.js":"6909","./ml":"02fb","./ml.js":"02fb","./mn":"958b","./mn.js":"958b","./mr":"39bd","./mr.js":"39bd","./ms":"ebe4","./ms-my":"6403","./ms-my.js":"6403","./ms.js":"ebe4","./mt":"1b45","./mt.js":"1b45","./my":"8689","./my.js":"8689","./nb":"6ce3","./nb.js":"6ce3","./ne":"3a39","./ne.js":"3a39","./nl":"facd","./nl-be":"db29","./nl-be.js":"db29","./nl.js":"facd","./nn":"b84c","./nn.js":"b84c","./oc-lnc":"167b","./oc-lnc.js":"167b","./pa-in":"f3ff","./pa-in.js":"f3ff","./pl":"8d57","./pl.js":"8d57","./pt":"f260","./pt-br":"d2d4","./pt-br.js":"d2d4","./pt.js":"f260","./ro":"972c","./ro.js":"972c","./ru":"957c","./ru.js":"957c","./sd":"6784","./sd.js":"6784","./se":"ffff","./se.js":"ffff","./si":"eda5","./si.js":"eda5","./sk":"7be6","./sk.js":"7be6","./sl":"8155","./sl.js":"8155","./sq":"c8f3","./sq.js":"c8f3","./sr":"cf1e","./sr-cyrl":"13e9","./sr-cyrl.js":"13e9","./sr.js":"cf1e","./ss":"52bd","./ss.js":"52bd","./sv":"5fbd","./sv.js":"5fbd","./sw":"74dc","./sw.js":"74dc","./ta":"3de5","./ta.js":"3de5","./te":"5cbb","./te.js":"5cbb","./tet":"576c","./tet.js":"576c","./tg":"3b1b","./tg.js":"3b1b","./th":"10e8","./th.js":"10e8","./tk":"5aff","./tk.js":"5aff","./tl-ph":"0f38","./tl-ph.js":"0f38","./tlh":"cf75","./tlh.js":"cf75","./tr":"0e81","./tr.js":"0e81","./tzl":"cf51","./tzl.js":"cf51","./tzm":"c109","./tzm-latn":"b53d","./tzm-latn.js":"b53d","./tzm.js":"c109","./ug-cn":"6117","./ug-cn.js":"6117","./uk":"ada2","./uk.js":"ada2","./ur":"5294","./ur.js":"5294","./uz":"2e8c","./uz-latn":"010e","./uz-latn.js":"010e","./uz.js":"2e8c","./vi":"2921","./vi.js":"2921","./x-pseudo":"fd7e","./x-pseudo.js":"fd7e","./yo":"7f33","./yo.js":"7f33","./zh-cn":"5c3a","./zh-cn.js":"5c3a","./zh-hk":"49ab","./zh-hk.js":"49ab","./zh-mo":"3a6c","./zh-mo.js":"3a6c","./zh-tw":"90ea","./zh-tw.js":"90ea"};function s(e){var t=o(e);return a(t)}function o(e){if(!a.o(n,e)){var t=new Error("Cannot find module '"+e+"'");throw t.code="MODULE_NOT_FOUND",t}return n[e]}s.keys=function(){return Object.keys(n)},s.resolve=o,e.exports=s,s.id="4678"},"4f87":function(e,t,a){},"56d7":function(e,t,a){"use strict";a.r(t),function(e){var t=a("2b0e"),n=a("5c96"),s=a.n(n),o=(a("0fae"),a("c1df")),r=a.n(o),c=a("3dfd"),d=a("a18c"),l=a("9923"),u=a("f1ed"),i=a("9607");a("9368"),a("43dd"),a("4f87");let f;e.navigator.language&&(f=e.navigator.language,f=f.toLowerCase()),0!==f.indexOf("en")&&(f="zh-cn"),r.a.locale(f),t["default"].config.debug=!0,t["default"].config.productionTip=!1,t["default"].use(s.a),new t["default"]({i18n:l["a"],router:d["a"],methods:u["a"],data:i["a"],render:e=>e(c["a"])}).$mount("#app")}.call(this,a("c8ba"))},"6e41":function(e,t,a){"use strict";t["a"]={add:"添加",cannot_empty:"不能为空",cancel:"取消",cancel_all:"全部取消",check_all:"全选",contribute_to_kube_file_browser:"为 KubeFileBrowser 做贡献",console:"控制台",container:"容器",copy:"上传下载",cpu:"处理器",dashboard:"仪表盘",delete:"删除",deleted:"已删除",deployment:"无状态集",description:"描述信息",download:"文件/目录下载",download_file:"下载",edit:"编辑",enter:"确定",file_browser:"文件浏览器",help:"帮助",image:"镜像",image_pull_secrets:"镜像拉取策略",keyword_search:"关键字搜索",name:"名称",namespace:"命名空间",network_error:"网络错误",more:"更多",mode:"权限",mod_time:"最后修改时间",multiupload:"多Pod文件上传",operate_success:"操作成功",operate:"操作",os:"系统类型",please_input_dest_path:"请输入目标路径",please_select_container:"请选择Container",please_select_deployment:"请选择Deployment",please_select_namespace:"请选择Namespace",please_select_pod:"请选择Pod",please_select_shell:"请选择Shell",prompt_message:"提示信息",ram:"内存",remove:"移除",refresh:"刷新",state:"状态",select_dir:"选取目录",select_file:"选取文件",size:"大小",tag:"版本",terminal:"终端命令行",unknown_error:"未知错误",upload:"单Pod文件上传",upload_all:"全部上传"}},"85ec":function(e,t,a){},9607:function(e,t,a){"use strict";t["a"]={PageSize:0,Page:0,Total:0,DialogSmallWidth:"500px",DialogNormalWidth:"750px",DialogLargeWidth:"900px",DialogNormalTop:"5vh"}},9923:function(e,t,a){"use strict";(function(e){var n=a("2b0e"),s=a("a925"),o=a("6e41"),r=a("3fab");n["default"].use(s["a"]);const c={"zh-cn":o["a"],en:r["a"]};let d;e.navigator.language&&(d=e.navigator.language,d=d.toLowerCase()),0!=d.indexOf("en")&&(d="zh-cn");const l=new s["a"]({locale:d,messages:c});t["a"]=l}).call(this,a("c8ba"))},a18c:function(e,t,a){"use strict";(function(e){a.d(t,"b",(function(){return c}));var n=a("2b0e"),s=a("8c4f"),o=a("9923");n["default"].use(s["a"]);const r=e=>()=>a("0351")("./"+e+".vue"),c=[{path:"/",component:r("layer"),name:"dashboard",title:o["a"].t("dashboard"),redirect:{name:"dashboard"},children:[{path:"dashboard",name:"dashboard",meta:{title:o["a"].t("dashboard"),icon:"el-icon-monitor",single:!0},component:r("dashboard")}]},{path:"/copy",name:"copy",title:o["a"].t("copy"),component:r("layer"),meta:{title:o["a"].t("copy"),icon:"el-icon-upload"},children:[{path:"upload",name:"upload",meta:{title:o["a"].t("upload")},component:r("upload")},{path:"multiupload",name:"multiupload",meta:{title:o["a"].t("multiupload")},component:r("multiupload")},{path:"download",name:"download",meta:{title:o["a"].t("download")},component:r("download")}]}],d=new s["a"]({routes:c,base:e,scrollBehavior:()=>({y:0}),mode:"history"});t["a"]=d}).call(this,"/")},f1ed:function(e,t,a){"use strict";var n=a("c1df"),s=a.n(n);t["a"]={MessageSuccess(e){this.$message({message:this.$t("operate_success"),type:"success",duration:1200,onClose:e})},PageInit(){this.$root.PageSize=7,this.$root.Page=1,this.$root.Total=0},PageReset(){this.$root.Total--;let e=Math.ceil(this.$root.Total/this.$root.PageSize);this.$root.Page>e&&(this.$root.Page=e),this.$root.Page<1&&(this.$root.Page=1)},PageOffset(){return this.$root.PageSize*(this.$root.Page-1)},ConfirmDelete(e,t){t||(t="此操作将永久删除该数据, 是否继续?"),this.$confirm(t,"提示",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then(()=>{e()}).catch(e=>{console.log(e)})},FormatDateTime(e,t){return e?(t||(t="YYYY-MM-DD HH:mm:ss"),s.a.unix(e).format(t)):"--"},FormatDateDuration(e){return s.a.duration(e).humanize(!1)},FormatDateFromNow(e){return e?s.a.unix(e).fromNow():"--"},Substr(e,t){if("[object String]"!=Object.prototype.toString.call(e))return"";let a="";return e.length>t&&(a="..."),e.substr(0,t)+a}}}});
//# sourceMappingURL=index.5eab3938.js.map