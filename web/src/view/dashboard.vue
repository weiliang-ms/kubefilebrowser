<template>
  <div>
    <el-card shadow="never">
      <div>
        <el-select v-model="namespace" @click.native="getNamespace" @change="selectedNamespace" style="width: 100%" filterable :placeholder="$t('please_select_namespace')">
          <el-option
              v-for="item in namespaces"
              :label="item"
              :value="item"
              :key="item"
          ></el-option>
        </el-select>
      </div>
      <div style="margin-top: 5px">
        <el-select :placeholder="$t('please_select_deployment')" style="width: 100%" filterable multiple v-model="deployment" @change="selecteddeployment">
          <el-option :label="$t('check_all')" value="all"></el-option>
          <el-option v-for="item in deployments" :label="item.label" :value="item.value" :key="item.value"></el-option>
        </el-select>
      </div>
      <div style="margin-top: 5px">
        <el-button @click.native="getStatus" style="float: right;margin-bottom: 10px">{{ $t('enter') }}</el-button>
      </div>
      <div>
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
          <el-table-column prop="Arch" :label="$t('arch')"></el-table-column>
          <el-table-column prop="Operate" :label="$t('operate')" width="160" align="center">
            <template slot-scope="scope">
              <el-button class="el-icon-s-fold" v-if="scope.row.OS==='windows'" @click.native="openTerminal(scope.row, 'cmd')">{{ $t('cmd') }}</el-button>
              <el-button class="el-icon-s-fold" v-else @click.native="openTerminal(scope.row, 'sh')">{{ $t('sh') }}</el-button>
              <span></span>
              <el-button class="el-icon-files" @click.native="openFileBrowser(scope.row, '/')">{{ $t('file_browser') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <el-dialog
        :visible.sync="dialogTerminalVisible"
        :title="$t('terminal')"
        center
        fullscreen
        :modal="false"
        :destroy-on-close="true"
        @opened="doOpened"
        @close="doClose"
    >
      <div style="margin-top: -25px;">
        <div ref="terminal" />
      </div>
    </el-dialog>
    <el-dialog
        center
        fullscreen
        :title="$t('file_browser')"
        :visible.sync="dialogFileBrowserVisible"
        @close="dialogFileBrowserVisible = false">
      <div style="margin-top: -25px;">
        <el-table-header store>
          <el-dropdown  type="info" class="avatar-container" trigger="click" style="height: 36px; float: right; margin-bottom: 10px;">
            <div class="avatar-wrapper">
              <el-button style="width: 90px; height: 30px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-s-tools" size="medium">
                {{ $t('create') }}
                <i class="el-icon-caret-bottom" />
              </el-button>
            </div>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <span style="display:block;" @click="openFileDialog(globalPath, 'create')">{{ $t('create_file') }}</span>
              </el-dropdown-item>
              <el-dropdown-item>
                <span style="display:block;" @click="createDir()">{{ $t('create_dir') }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
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
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, globalPath)" name="files" multiple="true">
                </span>
              </el-dropdown-item>
              <el-dropdown-item divided>
                <span class="fake-file-btn">
                  {{ $t('upload_dir') }}
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, globalPath)" name="files" webkitdirectory mozdirectory accept="*/*">
                </span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <el-dropdown type="primary" class="el-upload avatar-container" trigger="click" style="height: 36px;float: right;margin-bottom: 10px;">
            <div class="avatar-wrapper">
              <el-button style="width: 120px; height: 30px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="primary" round class="el-icon-download" size="medium">
                {{ $t('bulk_download') }}
                <i class="el-icon-caret-bottom" />
              </el-button>
            </div>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <span style="display:block;" @click="bulkDownload(bulkPath, 'tar')">TAR{{ $t('download') }}</span>
              </el-dropdown-item>
              <el-dropdown-item divided>
                <span style="display:block;" @click="bulkDownload(bulkPath, 'zip')">ZIP{{ $t('download') }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <ul>
            <li style="float: left; margin-top: 10px; list-style: none;" v-for="(item) in headerPaths">
              <a style="margin-right: 5px; font-size: 16px" class="el-icon-folder-opened" @click="openFileBrowser(null, item.path)">{{item.name}}</a>
            </li>
          </ul>
          &nbsp;&nbsp;&nbsp;&nbsp;
          <span style="float: left;">
            <el-button type="info" style="padding: 3px;margin-top: 8px;" icon="el-icon-refresh" circle @click="openFileBrowser(null, path)"></el-button>
          </span>
        </el-table-header>
        <el-table
            id="tableData"
            class="app-table"
            border
            style="width: 100%"
            size="100%"
            :cell-style="{padding:'6px 0'}"
            :data="fileBrowserData"
            @selection-change="handleSelectionChange"
            :default-sort="{prop: 'Name', order: 'ascending'}">
          <el-table-column type="selection"></el-table-column>
          <el-table-column
              min-width="80px"
              prop="Name"
              :label="$t('name')"
              sortable
              :sort-orders="['ascending', 'descending']"
          >
            <template slot-scope="scope">
              <div class="el-icon-folder"  v-if="scope.row.IsDir" @click="openFileBrowser(null, scope.row.Path)">&nbsp;&nbsp;{{scope.row.Name}}</div>
              <div class="el-icon-files"  v-else @click="openFileDialog(scope.row.Path, 'open')">{{scope.row.Name}}</div>
            </template>
          </el-table-column>
          <el-table-column
              prop="Size"
              :label="$t('size')"
              sortable
              :sort-orders="['ascending', 'descending']"
          >
          </el-table-column>
          <el-table-column
              prop="Mode"
              :label="$t('mode')"
          >
          </el-table-column>
          <el-table-column
              prop="ModTime"
              :label="$t('mod_time')"
              sortable
              :sort-orders="['ascending', 'descending']"
          >
          </el-table-column>
          <el-table-column
              prop="Download"
              :label="$t('operate')" align="center"
          >
            <template slot-scope="scope">
              <el-dropdown type="info" class="avatar-container" trigger="click" style="height: 36px;font-size: 9px">
                <div class="avatar-wrapper">
                  <el-button style="width: 90px; height: 30px; margin-top: 4px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-s-tools" size="medium">
                    {{ $t('operate') }}
                    <i class="el-icon-caret-bottom" />
                  </el-button>
                </div>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item v-if="!scope.row.IsDir">
                    <span class="fake-file-btn" @click="openFileDialog(scope.row.Path, 'open')">{{ $t('change') }}</span>
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <span class="fake-file-btn" @click="openRenameDialog(scope.row.Name)">{{ $t('rename') }}</span>
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <span class="fake-file-btn" @click="removeFileOrDir(scope.row.Path)">{{ $t('remove') }}</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
              <span>
                &nbsp;&nbsp;
              </span>
              <el-dropdown v-if="scope.row.IsDir" type="success" class="avatar-container" trigger="click" style="height: 36px;font-size: 9px">
                <div class="avatar-wrapper">
                  <el-button style="width: 90px; height: 30px; margin-top: 4px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-upload" size="medium">
                    {{ $t('upload') }}
                    <i class="el-icon-caret-bottom" />
                  </el-button>
                </div>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item>
                    <span class="fake-file-btn">
                      {{ $t('upload_file') }}
                      <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, scope.row.Path)" name="files" multiple="true">
                    </span>
                  </el-dropdown-item>
                  <el-dropdown-item divided>
                    <span class="fake-file-btn">
                      {{ $t('upload_dir') }}
                      <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, scope.row.Path)" name="files" webkitdirectory mozdirectory accept="*/*">
                    </span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
              <span>
                &nbsp;&nbsp;
              </span>
              <el-dropdown type="primary" class="avatar-container" trigger="click" style="height: 36px;font-size: 9px">
                <div class="avatar-wrapper">
                  <el-button style="width: 90px; height: 30px; margin-top: 4px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="primary" round class="el-icon-download" size="medium">
                    {{ $t('download') }}
                    <i class="el-icon-caret-bottom" />
                  </el-button>
                </div>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item>
                    <span style="display:block;" @click="download(scope.row.Path, 'tar')">TAR{{ $t('download') }}</span>
                  </el-dropdown-item>
                  <el-dropdown-item divided>
                    <span style="display:block;" @click="download(scope.row.Path, 'zip')">ZIP{{ $t('download') }}</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
        <el-table-footer store>
        </el-table-footer>
      </div>
    </el-dialog>
    <el-dialog
        :append-to-body="true"
        :with-header="false"
        width="80%"
        :title="$t('file')"
        :visible.sync="dialogFileVisible"
        @close="dialogFileVisible = false">
      <div style="margin-top: -45px">
        <span style="display:block; float: left;font-size: 25px;margin-top: 20px; margin-left: 12px">{{createForPath}}</span>
        <el-input v-if="isNewFile" v-model="createName" size="small" style="margin-top: 17px; margin-left: 6px;float: left; width: auto;" autocomplete="off" :placeholder="$t('please_input_name')"></el-input>
        <el-button @click.native="saveFile" style="float: right;margin-right: 12px;margin-top: 12px">{{ $t('enter') }}</el-button>
        <el-input v-model="fileContent" rows="15" type="textarea" style="margin-top: 15px" :placeholder="$t('please_input_content')"></el-input>
        <!--        style="margin-top: 15px;margin-bottom: 60px;padding-bottom: 70px;height: 600px;"-->
        <!--        <quill-editor-->
        <!--            style="margin-top: 15px;height: 100%"-->
        <!--            v-model="fileContent"-->
        <!--            :options="editorOption">-->
        <!--        </quill-editor>-->
      </div>
    </el-dialog>
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
import {GetStatus} from '../api/status'
import {GetNamespace} from '../api/namespaces'
import {GetDeployment} from "../api/deployment";
import {
  FileBrowserCreateDir,
  FileBrowserCreateFile,
  FileBrowserList,
  FileBrowserOpen,
  FileBrowserRemove, FileBrowserRename
} from "../api/filebrowser";
import {FileOrDirUpload} from "../api/upload";
import { Terminal } from 'xterm'
import * as fit from 'xterm/lib/addons/fit/fit'
import { Base64 } from 'js-base64'
import * as webLinks from 'xterm/lib/addons/webLinks/webLinks'
import * as search from 'xterm/lib/addons/search/search'
import 'xterm/lib/addons/fullscreen/fullscreen.css'
import 'xterm/dist/xterm.css'

const toolbarOptions = [
  ["bold", "italic", "underline", "strike"], // 加粗 斜体 下划线 删除线
  ["blockquote", "code-block"], // 引用  代码块
  [{ header: 1 }, { header: 2 }], // 1、2 级标题
  [{ list: "ordered" }, { list: "bullet" }], // 有序、无序列表
  [{ script: "sub" }, { script: "super" }], // 上标/下标
  [{ indent: "-1" }, { indent: "+1" }], // 缩进
  // [{'direction': 'rtl'}],                         // 文本方向
  [{ size: ["small", false, "large", "huge"] }], // 字体大小
  [{ header: [1, 2, 3, 4, 5, 6, false] }], // 标题
  [{ color: [] }, { background: [] }], // 字体颜色、字体背景颜色
  [{ font: [] }], // 字体种类
  [{ align: [] }], // 对齐方式
  ["clean"], // 清除文本格式
  ["link"] // 链接、图片
];

const defaultTheme = {
  foreground: '#ffffff', // 字体
  background: '#1b212f', // 背景色
  cursor: '#ffffff', // 设置光标
  selection: 'rgba(255, 255, 255, 0.3)',
  black: '#000000',
  brightBlack: '#808080',
  red: '#ce2f2b',
  brightRed: '#f44a47',
  green: '#00b976',
  brightGreen: '#05d289',
  yellow: '#e0d500',
  brightYellow: '#f4f628',
  magenta: '#bd37bc',
  brightMagenta: '#d86cd8',
  blue: '#1d6fca',
  brightBlue: '#358bed',
  cyan: '#00a8cf',
  brightCyan: '#19b8dd',
  white: '#e5e5e5',
  brightWhite: '#ffffff'
}
const bindTerminalResize = (term, websocket) => {
  const onTermResize = size => {
    websocket.send(
        Base64.encode(
            JSON.stringify({
              type: 'resize',
              rows: size.rows,
              cols: size.cols
            })
        )
    )
  }
  // register resize event.
  term.on('resize', onTermResize)
  // unregister resize event when WebSocket closed.
  websocket.addEventListener('close', function() {
    term.off('resize', onTermResize)
  })
}
const bindTerminal = (term, websocket, bidirectional, bufferedTime) => {
  term.socket = websocket
  let messageBuffer = null
  const handleWebSocketMessage = function(ev) {
    if (bufferedTime && bufferedTime > 0) {
      if (messageBuffer) {
        messageBuffer += Base64.decode(ev.data)
      } else {
        messageBuffer = Base64.decode(ev.data)
        setTimeout(function() {
          term.write(messageBuffer)
        }, bufferedTime)
      }
    } else {
      term.write(Base64.decode(ev.data))
    }
  }
  const handleTerminalData = function(data) {
    websocket.send(
        Base64.encode(
            JSON.stringify({
              type: 'input',
              input: data
            })
        )
    )
  }
  websocket.onmessage = handleWebSocketMessage
  if (bidirectional) {
    term.on('data', handleTerminalData)
  }
  // send heartbeat package to avoid closing webSocket connection in some proxy environmental such as nginx.
  const heartBeatTimer = setInterval(function() {
    websocket.send(
        Base64.encode(
            JSON.stringify({
              type: 'heartbeat',
              data: ''
            })
        )
    )
    // websocket.send('1')
  }, 20 * 1000)
  websocket.addEventListener('close', function() {
    websocket.removeEventListener('message', handleWebSocketMessage)
    term.off('data', handleTerminalData)
    delete term.socket
    clearInterval(heartBeatTimer)
  })
}

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
      bulkPath: [],
      globalPath: "",
      headerPaths: [],
      dialogTerminalVisible: false,
      dialogFileBrowserVisible: false,
      dialogFileVisible: false,
      wsUrl: "",
      isFullScreen: false,
      searchKey: '',
      v: this.visible,
      ws: null,
      term: null,
      thisV: this.visible,
      createForPath: "",
      createName: "",
      fileContent: "",
      isNewFile: false,
      editorOption: {
        theme: "snow", // or 'bubble'
        modules: {
          toolbar: {
            container: toolbarOptions,
          }
        }
      },
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
          for(const key in data){
            this.namespaces.push(data[key].metadata.name)
          }
        }
      })
    },
    selectedNamespace(options) {
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
        }
      })
    },
    selecteddeployment(options) {
    },
    getStatus() {
      let deployment = this.deployment
      if (this.deployment[0] === "all") {
        deployment = []
        for (const key in this.deployments) {
          deployment.push(this.deployments[key].value)
        }
      }
      GetStatus({namespace: this.namespace, deployment:deployment}).then(res => {
        if (res) {
          this.tableData = []
          for (const i in res) {
            const pod_name =res[i].pod_name;
            const cData = res[i].containers;
            for (const j in cData) {
              let tr = {
                Pods:pod_name,
                Container:cData[j].name,
                Image:cData[j].image,
                Tag:cData[j].version,
                ImagePullSecrets:cData[j].image_pull_policy,
                State:cData[j].state,
                CPU:cData[j].cpu,
                RAM:cData[j].ram,
                OS: cData[j].os,
                Arch: cData[j].arch,
              }
              this.tableData.push(tr)
            }
          }
        }
      })
    },
    openTerminal(options, shell) {
      this.dialogTerminalVisible = true
      this.wsUrl = "ws://"+window.location.host+"/api/k8s/terminal?namespace="+this.namespace+"&pods="+options.Pods+"&container="+options.Container+"&shell="+shell;
    },
    openFileBrowser(options, path) {
      if (path === undefined) {
        path = "/"
      }
      if (path === "/" && options !== null) {
        this.pods = options.Pods
        this.container = options.Container
      }
      this.headerPaths = []
      this.globalPath=path
      this.headerPaths.push(path)
      if (path !== undefined) {
        let _p = path.split('/')
        let _pa = ""
        this.headerPaths = []
        _p.forEach((item,index) => {
          if (index === 0) {
            _pa = "/"
            item = "/"
            this.headerPaths.push({
              name: item,
              path: _pa,
            })
          }
          if (index !== 0 && item !== "") {
            _pa += item + "/"
            this.headerPaths.push({
              name: item,
              path: _pa,
            })
          }
        })
      }
      this.path = path
      this.fileBrowserData = []
      FileBrowserList({
        namespace: this.namespace,
        pods: this.pods,
        container: this.container,
        path: path,
      }).then(res => {
        this.dialogFileBrowserVisible = true
        this.fileBrowserData = []
        if (res !== undefined) {
          this.fileBrowserData = res
        }
      }, err => {
        this.$message.error(err.info.message)
      })
    },
    openFileDialog(path, type) {
      this.fileContent = ""
      this.dialogFileVisible=true
      this.createForPath = path
      this.isNewFile = true
      if (type === "open") {
        this.isNewFile = false
        this.createName = ""
        FileBrowserOpen({
          namespace: this.namespace,
          pods: this.pods,
          container: this.container,
          path: path,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.fileContent = res
          }
        }, err => {
          this.dialogFileVisible=false
          this.$message.error(err.info.message)
        })
      }
    },
    saveFile() {
      let path = this.createForPath+"/"+this.createName
      if (this.createName === "") {
        path = this.createForPath
      }
      FileBrowserCreateFile(this.fileContent, {
        namespace: this.namespace,
        pods: this.pods,
        container: this.container,
        path: path,
      }).then(res => {
        console.log(res)
        if (res !== undefined) {
          this.dialogFileVisible = false
          this.fileContent = ""
          this.createName = ""
          this.$message.success(res)
          this.openFileBrowser(null, this.path)
        }
      }, err => {
        this.$message.error(err.info.message)
      })
    },
    createDir() {
      this.$prompt(this.$t('please_input_name'), this.$t('tips'), {
        confirmButtonText: this.$t('enter'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning'
      }).then(({value}) => {
        if(!value) {//对输入内容校验
          return this.$t('please_input_name');
        }
        FileBrowserCreateDir({
          namespace: this.namespace,
          pods: this.pods,
          container: this.container,
          path: this.path+"/"+value,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.$message.success(res)
            this.openFileBrowser(null, this.path)
          }
        }, err => {
          this.$message.error(err.info.message)
        })
      }).catch(() => {
        this.$message.info(this.$t('cancel'))
      });
    },
    openRenameDialog(oldName) {
      this.$prompt(this.$t('please_input_name')+"<br>"+oldName, this.$t('tips'), {
        confirmButtonText: this.$t('enter'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning'
      }).then(({value}) => {
        if(!value) {//对输入内容校验
          return this.$t('please_input_name');
        }
        FileBrowserRename({
          namespace: this.namespace,
          pods: this.pods,
          container: this.container,
          old_path: this.path+"/"+oldName,
          path: this.path+"/"+value,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.$message.success(res)
            this.openFileBrowser(null, this.path)
          }
        }, err => {
          this.$message.error(err.info.message)
        })
      }).catch(() => {
        this.$message.info(this.$t('cancel'))
      });
    },
    removeFileOrDir(path) {
      this.$confirm(this.$t('tips_msg')+"<br>"+path, this.$t('tips'), {
        confirmButtonText: this.$t('enter'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning'
      }).then(() => {
        FileBrowserRemove({
          namespace: this.namespace,
          pods: this.pods,
          container: this.container,
          path: path,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.$message.success(res)
            this.openFileBrowser(null, this.path)
          }
        }, err => {
          this.$message.error(err.info.message)
        })
      }).catch(() => {
        this.$message.info(this.$t('cancel'))
      });
    },
    handleSelectionChange(val) {
      this.bulkPath = []
      val.forEach((item) => {
        this.bulkPath.push(item.Path)
      })
    },
    download(path, style) {
      const url = "/api/k8s/download?namespace="+this.namespace+"&pod_name="+this.pods+"&container_name="+this.container+"&dest_path="+path+"&style="+style;
      const xhr = new XMLHttpRequest();
      xhr.open('GET', url, true);        // 也可以使用POST方式，根据接口
      xhr.responseType = "blob";    // 返回类型blob
      xhr.setRequestHeader("Content-type", "application/json;charset=UTF-8");
      // 定义请求完成的处理函数，请求前也可以增加加载框/禁用下载按钮逻辑
      xhr.onload = function () {
        // 请求完成
        if (this.status === 200) {
          // 返回200
          let content = xhr.response;
          let eLink = document.createElement('a');
          eLink.download = this.getResponseHeader('X-File-Name');
          eLink.style.display = 'none';
          let blob = new Blob([content]);
          eLink.href = URL.createObjectURL(blob);
          document.body.appendChild(eLink);
          eLink.click();
          document.body.removeChild(eLink);
        }
      };
      // 发送ajax请求
      xhr.send()
    },
    bulkDownload(paths, style) {
      if (paths.length === 0) {
        this.$message.error(this.$t('cannot_empty'))
        return
      }
      let path = ""
      paths.forEach(item => {
        path += "&dest_path="+item
      })
      const url = "/api/k8s/download?namespace="+this.namespace+"&pod_name="+this.pods+"&container_name="+this.container+"&dest_path="+path+"&style="+style;
      const xhr = new XMLHttpRequest();
      xhr.open('GET', url, true);        // 也可以使用POST方式，根据接口
      xhr.responseType = "blob";    // 返回类型blob
      xhr.setRequestHeader("Content-type", "application/json;charset=UTF-8");
      // 定义请求完成的处理函数，请求前也可以增加加载框/禁用下载按钮逻辑
      xhr.onload = function () {
        // 请求完成
        if (this.status === 200) {
          // 返回200
          let content = xhr.response;
          let eLink = document.createElement('a');
          eLink.download = this.getResponseHeader('X-File-Name');
          eLink.style.display = 'none';
          let blob = new Blob([content]);
          eLink.href = URL.createObjectURL(blob);
          document.body.appendChild(eLink);
          eLink.click();
          document.body.removeChild(eLink);
        }
      };
      // 发送ajax请求
      xhr.send()
    },
    uploadFileOrDir(e, path) {
      const files = e.target.files;
      if (files.length === 0 ) {
        e.target.value = ""
        return
      }
      const formData = new FormData();
      //追加文件数据
      for (let i = 0; i < files.length; i++) {
        formData.append("files", files[i]);
      }
      FileOrDirUpload(formData, {
        namespace:this.namespace,
        pod_name:this.pods,
        container_name:this.container,
        dest_path:path},{"Content-Type":"multipart/form-data"}).then((res) => {
        if (res.failure !== undefined) {
          this.$message.warning(res.failure)
        }else {
          this.$message.success(res.success)
        }
      }, (err) => {
        this.$message.error(err.info.message)
      })
      e.target.value = ""
    },

    onWindowResize() {
      // console.log("resize")
      // this.term.fit() // it will make terminal resized.
      // this.term.scrollToBottom();
      let height = document.body.clientHeight;
      let rows = height/23;
      this.term.fit();
      this.term.resize(this.term.cols,parseInt(rows))//终端窗口重新设置大小 并触发term.on("resize"
      this.term.scrollToBottom();
    },
    doLink(ev, url) {
      if (ev.type === 'click') {
        window.open(url)
      }
    },
    doClose() {
      window.removeEventListener('resize', this.onWindowResize)
      // term.off("resize", this.onTerminalResize);
      if (this.ws) {
        this.ws.close()
      }
      if (this.term) {
        this.term.dispose()
      }
      this.$emit('pclose', false)// 子组件对openStatus修改后向父组件发送事件通知
    },
    doOpened() {
      Terminal.applyAddon(fit)
      Terminal.applyAddon(webLinks)
      Terminal.applyAddon(search)
      this.term = new Terminal({
        rendererType: 'canvas', // 渲染类型
        rows: parseInt(document.body.clientHeight/23),
        cols: parseInt(document.body.clientWidth),
        convertEol: true, // 启用时，光标将设置为下一行的开头
        // scrollback: 10, // 终端中的回滚量
        disableStdin: false, // 是否应禁用输入
        fontSize: 18,
        cursorBlink: true, // 光标闪烁
        cursorStyle: 'bar', // 光标样式 underline
        bellStyle: 'sound',
        theme: defaultTheme
      })
      this.term._initialized = true
      this.term.prompt = () => {
        this.term.write('\r\n')
      }
      this.term.prompt()
      this.term.on('key', function(key, ev) {
        console.log(key, ev, ev.keyCode)
      })
      this.term.open(this.$refs.terminal)
      this.term.webLinksInit(this.doLink)
      // term.on("resize", this.onTerminalResize);
      this.term.on('resize', this.onWindowResize)
      window.addEventListener('resize', this.onWindowResize)
      this.term.fit() // first resizing
      this.ws = new WebSocket(this.wsUrl)
      this.ws.onerror = () => {
        this.$message.error(this.$t('web_socker_connection_failed'))
      }
      this.ws.onclose = () => {
        this.term.setOption('cursorBlink', false)
        this.$message(this.$t('web_socket_disconnect'))
      }
      bindTerminal(this.term, this.ws, true, -1)
      bindTerminalResize(this.term, this.ws)
    }
  },
  watch:{
    visible(val) {
      this.v = val// 新增result的watch，监听变更并同步到myResult上
    },
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