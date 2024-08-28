<template>
  <a-page-header style="border: 0px solid rgb(235, 237, 240); accent-color: red">
    <slot name="title">
      <h2 style="color: whitesmoke">部署服务</h2>
    </slot>
  </a-page-header>
  <a-list :grid="{ gutter: 24, column: 4 }" :data-source="data">
    <template #renderItem="{ item }">
      <a-list-item>
        <a-card :title="item.name">
          <template #actions>
            <DeploymentUnitOutlined
              style="color: darkturquoise"
              key="deploy"
              @click="deploymentClick(item)"
            />
            <a-popconfirm
              title="确认要重启服务吗?"
              ok-text="重启"
              cancel-text="放弃"
              @confirm="confirmReboot(item.ID)"
            >
              <RetweetOutlined style="color: blue" key="reboot" />
            </a-popconfirm>
            <a-popconfirm
              title="确认要关闭服务吗?"
              ok-text="关闭"
              cancel-text="放弃"
              @confirm="confirmClose(item.ID)"
            >
              <CloseCircleOutlined style="color: red" key="close" />
            </a-popconfirm>
            <FolderViewOutlined
              style="color: black"
              key="view"
              @click="viewFolder(item.real_path)"
            />
          </template>
          <template #extra>
            <a-flex horizontal gap="small">
              <a @click="editService(item)">编辑</a>
              <a-popconfirm
              title="确认要删除服务吗?"
              ok-text="删除"
              cancel-text="放弃"
              @confirm="confirmDelService(item.ID)"
            >
            <a @click="delService(item.ID)">删除</a>
          </a-popconfirm>
              
            </a-flex>
          </template>
          <a-flex vertical>
            <a-flex horizontal gap="small">
              <span>服务名: </span><span>{{ item.service_name }}</span>
            </a-flex>
            <!-- <a-flex horizontal gap="small">
              <span>CMD: </span><span>{{ item.cmd_line }}</span>
            </a-flex> -->
            <a-flex horizontal gap="small">
              <span>路径: </span><span>{{ item.real_path }}</span>
            </a-flex>
            <a-flex horizontal gap="small">
              <span>类型: </span><span>{{ formatServiceType(item.service_type) }}</span>
            </a-flex>
            <a-flex horizontal gap="small">
              <span>监听: </span><span>{{ item.listen_url }}</span>
            </a-flex>
            <a-flex horizontal gap="small">
              <span>状态: </span><span :style="item.color">{{ item.status }}</span>
            </a-flex>
          </a-flex>
        </a-card>
      </a-list-item>
    </template>
  </a-list>



  <a-float-button-group trigger="click" type="primary" :style="{ right: '24px' }">
    <template #icon>
      <MenuUnfoldOutlined />
    </template>
    <a-float-button tooltip="切换用户" @click="switchUserClick">
      <template #icon>
        <UserSwitchOutlined />
      </template>
    </a-float-button>
    <a-float-button tooltip="服务管理" @click="serviceManageClick">
      <template #icon>
        <AppstoreOutlined />
      </template>
    </a-float-button>
    <a-float-button tooltip="用户管理" @click="userManageClick">
      <template #icon>
        <UserOutlined />
      </template>
    </a-float-button>
  </a-float-button-group>


  <a-modal v-model:open="viewFolderOpen" :title="viewFolderTitle" footer="" okText="确认" cancelText="取消">
    <a-flex vertical gap="middle">
      <a-tag>{{ curPath }}</a-tag>
      <a-tag> {{ selectedPath }}</a-tag>
      <a-directory-tree multiple :tree-data="treeData" :load-data="onLoadData" @select="treeSelect">
      </a-directory-tree>
      <a-button @click="selPath" :disabled="selBtnDisabled" >选择</a-button>
    </a-flex>
  </a-modal>

  <a-modal v-model:open="viewIISListOpen" title="IIS服务列表" footer="" okText="确认" cancelText="取消">
    <a-flex vertical gap="middle">
      <a-tag> 你选择了: {{ selectedIIS.Name }}</a-tag>
      <a-list item-layout="horizontal" :data-source="iisList" selectable>
        <template #renderItem="{ item }">
          <a-list-item @click="iisListClick(item)">
            <a-list-item-meta
              :description="item.PhysicalPath"
            >
              <template #title>
                <span>{{ item.Name }} &lt;== {{ item.Binding }}</span>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>

      <a-button @click="selIIS" :disabled="selBtnDisabled">选择</a-button>
    </a-flex>
  </a-modal>

  <a-modal v-model:open="serviceModalOpen" :title="tt" @ok="updateService" okText="确认" cancelText="取消">
    <a-flex horizontal gap="middle">
      <a-flex vertical gap="middle">
        <a-input v-model:value="model.name" placeholder="请输入服务名称"> </a-input>
        <a-input v-model:value="model.service_name" placeholder="请输入服务名"> </a-input>
        <a-input v-model:value="model.arg" placeholder="请输入参数"> </a-input>
        <a-input v-model:value="model.real_path" placeholder="请输入路径"> </a-input>
        <a-input v-model:value="model.listen_url" placeholder="请输入监听地址"> </a-input>
        <a-select v-model:value="model.service_type" size="middle" :options="serviceOptions">

        </a-select>
      </a-flex>
      <a-flex vertical gap="middle">
        <a-flex horizontal>
          <span>上次部署:</span><span>{{ format(model.last_deploy_time) }}</span>
        </a-flex>
        <a-flex horizontal>
          <span>加入时间:</span><span>{{ format(model.CreatedAt) }}</span>
        </a-flex>
        <a-flex horizontal>
          <span>修改时间:</span><span>{{ format(model.UpdatedAt) }}</span>
        </a-flex>
        <a-button @click="selServerPath">选择服务器路径</a-button>
        <a-button @click="selServerIIS">选择IIS服务</a-button>
      </a-flex>
    </a-flex>
  </a-modal>

  <a-modal v-model:open="deployModalOpen" :title="modelDeploy.deployTitle" :ok-button-props="{ disabled: true }" cancelText="取消">
    <a-flex vertical gap="middle">
      <a-flex horizontal gap="middle">
        <a-select v-model:value="modelDeploy.deployType" size="middle" style="min-width: 150px;">
          <a-select-option value="文件">文件</a-select-option>
          <a-select-option value="目录">目录</a-select-option>
          <a-select-option value="zip压缩包">压缩包</a-select-option>
        </a-select>
        <a-button @click="deploy">选择{{modelDeploy.deployType}}部署</a-button>
        <input type="file" id="dir-selector" hidden @change="selectA" webkitdirectory directory multiple />
        <input type="file" id="file-selector" hidden @change="selectB" webkitfile />
        <input type="file" id="multifile-selector" hidden @change="selectA" webkitfile multiple />
      </a-flex>
      <a-textarea v-model:value="modelDeploy.paths" :rows="5"></a-textarea>
      
      <TerminalComponent ref="terminal" v-if="showTerminal" />
      
    </a-flex>
  
  </a-modal>

 
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import { showError,showSuccess,msgSuccess, msgInfo } from '@/utils/notification'
import type { TreeProps,SelectProps  } from 'ant-design-vue'
import dayjs from 'dayjs'
import {
  DeploymentUnitOutlined,
  RetweetOutlined,
  CloseCircleOutlined,
  FolderViewOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  UserSwitchOutlined,
  AppstoreOutlined
} from '@ant-design/icons-vue'
import JSZip from 'jszip'
import {SparkMD5} from 'spark-md5';
import iconv from 'iconv-lite'
import {apiFetch} from '@/utils/ofetch'
import TerminalComponent from '@/components/TerminalComponent.vue'
import { Terminal } from '@xterm/xterm'
const terminal = ref<any>(null)
const showTerminal = ref(false)
import { saveAs } from 'file-saver'



const zip = new JSZip();
const data = ref([])

const router = useRouter()
const treeData = ref<TreeProps['treeData']>([])

const onLoadData: TreeProps['loadData'] = (treeNode) => {
  return new Promise<void>((resolve) => {
    if (treeNode.isLeaf) {
      resolve()
      return
    }
    apiFetch('/api/service/dir_tree', {
      method: 'GET',
      query: {
        home: treeNode.key
      },
    }).then((resp) => {
      //插入子目录到对应的位置
      
      let a = resp.data.map((r:any) => {
        return {
          title: r.path,
          key: r.path,
          pathType: r.info.is_dir,
          isLeaf: !r.info.is_dir
        }
      })
      treeNode.dataRef.children = a
      treeData.value = [...treeData.value]
      resolve()
    })
  })
}

const viewFolderOpen = ref(false)
const viewFolderTitle = ref('目录查看')
const curPath = ref('')
const selBtnDisabled = ref(true)

const serviceOptions = ref<SelectProps['options']>([
  {
    value: 0,
    label: 'IIS网站'
  },
  {
    value: 1,
    label: '终端程序'
  },
  {
    value: 2,
    label: 'Windows服务'
  },
  {
    value: 3,
    label: 'Systemd服务'
  },
  {
    value: 4,
    label: '静态目录'
  },
])

const serviceModalOpen = ref(false)
const tt = ref('修改服务')
const model = ref<any>({})
const isCanSelect = ref(false)



const selectedPath = ref('未选择')
const selectedIIS = ref<any>({})
const iisList = ref([])
const viewIISListOpen = ref(false)
const deployModalOpen = ref(false)
const modelDeploy = ref<any>({
  deployType:'文件',
  deployTypeName: '文件',
  paths:'',
  deployTitle:'部署服务-',
  model:{}
})


function writeToTerminal(data: any) {
  if (terminal.value) {
    var xterm: Terminal = terminal.value?.xterm
    xterm.writeln(data)
  }
}
async function selectB(e:any) {
  console.log(e.target.files)
  const files:FileList = e.target.files
  modelDeploy.value.paths = files[0].name
  doDeploy(files[0])

}

async function startDeploy(id:number) {
  showTerminal.value = true
  let eventSource = new EventSource('http://10.10.0.21:3058/api/deploy/sse?id='+id)
  eventSource.onmessage = (event)=>{
    writeToTerminal(event.data)
    if (event.data==='done') {
       setTimeout(()=>{
        showTerminal.value =  false
        modelDeploy.value.paths = ''
        showSuccess('部署完成')
       },2000)
    }
  }
  eventSource.onerror = (event:any) =>{
    console.log(event)
    eventSource.close()
    writeToTerminal(event)
  }

}
function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
async function selectA(e:any) {
  console.log(e)
  const files:FileList = e.target.files
  
  if (files.length >=1) {
    msgInfo('正在压缩...')
    for (const file of files) {
      modelDeploy.value.paths += file.name + '\n'
      console.log(file.name)
      console.log(iconv.encode(file.name,'utf8').toString())
      zip.file( iconv.encode(file.name,'utf8').toString(), file);
    }
    await sleep(2000)
    zip.generateAsync({
      type:"blob"
    })
      .then((content)=> {
        //saveAs(content,'hello.zip')
        
        doDeploy(content)
      });

  }

  
}
function convertToUtf8String(str: string): string {
  // Encode the string into UTF-8 bytes
  const utf8Bytes = new TextEncoder().encode(str);

  // Decode the UTF-8 bytes back into a string
  const utf8String = new TextDecoder('utf-8').decode(utf8Bytes);

  return utf8String;
}
function formatServiceType(st: number) {
  switch (st) {
    case 0:
        return 'IIS网站'
    case 1:
      return 'console终端'
    case 2:
      return 'windows 服务'
    case 3:
      return 'systemd服务'
    case 4:
      return '静态目录'

    default:
      break
  }
}

async function doDeploy(content:any) {
  const formData = new FormData();
  //let h = SparkMD5.hash(content)
  formData.append("file", content);
  //formData.append("hash", h);
  //formData.append("size", content.length);
  formData.append("serviceId", modelDeploy.value.model.ID);
   apiFetch('/api/deploy',{
    method: 'POST',
     
    body: formData,
   }).then(resp=>{
    console.log(resp)
    startDeploy(resp.data)
  
   })
}

function delService(serviceId:any) {

}

function serviceManageClick() {
  router.push('/services')
}

function deploy() {
  console.log('doOpenSelect',modelDeploy.value)
  if (modelDeploy.value.deployType=='文件') {
     document.getElementById('multifile-selector')?.click();
  } else if (modelDeploy.value.deployType == '目录') {
    document.getElementById('dir-selector')?.click();
  } else if (modelDeploy.value.deployType=='zip压缩包') {
    document.getElementById('file-selector')?.setAttribute('accept','.zip')
    document.getElementById('file-selector')?.click();
  }
}


//{selected: bool, selectedNodes, node, event}
function treeSelect(selectedKeys: any, e: { node: { isLeaf: any; key: string }; selected: boolean }) {
  if (!e.node.isLeaf) {
    selectedPath.value = e.node.key
    if (!isCanSelect.value) {
      console.log('set sel btn enable')
      selBtnDisabled.value = false
    }
  } else {
    e.selected = false
    selectedPath.value = "未选择"
    selBtnDisabled.value = true
  }
}

function format(inParam: string | number | Date | dayjs.Dayjs | null | undefined) {
  return dayjs(inParam).format('YYYY-MM-DD HH:mm:ss')
}

function showEditOrAdd(isEdit: boolean, model1: any = undefined) {
  if (isEdit) {
    tt.value = '修改服务'
    model.value = model1
    console.log(model.value)
  } else {
    tt.value = '添加服务'
    model.value = {
      name: '',
      service_name: '',
      arg: '',
      real_path: '',
      listen_url: '',
      service_type: '',
      last_deploy_time: '',
      CreatedAt: '',
      UpdatedAt: ''
    }
  }
  serviceModalOpen.value = true
}
function updateService() {
  apiFetch('/api/service',{
    method: 'POST',
    body: model.value
  })
  .then(data=>{
    console.log(data)
    if (data.code===200) {
      showSuccess(data.message)
      serviceModalOpen.value = false
    } else {
      showError(data.message)
    }
  }).catch(err=>{
    console.log(err)
  })

}





function editService(item: any) {
  console.log('editService', JSON.stringify(item))
  showEditOrAdd(true, item)
  refreshData()
}
function refreshStatus() {
  data.value.forEach((item:any,_) =>{
      apiFetch('/api/service/status',{
        query: {
          'service_id': item.ID
        }
      }).then((resp)=>{
        //color: green;
        if (resp.data.status==0) {
          item.color = 'color: green;'
          item.status = '运行中'
        } else if (resp.data.status==1) {
          item.color = 'color: red;'
          item.status = '已停止'
        } else {
          item.color = 'color: black;'
          item.status = '未知'
        }
        item.cmd_line = resp.data.cmdline
      })

    })
}

onMounted(() => {
  refreshData()
  //refreshStatus()
  
})
function refreshData() {
  apiFetch('/api/service/list', {
      query: { page: 1,size:10 }
    }).then((resp) => {
      data.value = resp.data.list
      //console.log(data.value)
      refreshStatus()
    })
}

function switchUserClick() {
  apiFetch('/api/user/signOut',{
    method: 'POST'
  })
  router.push('/')
}
function userManageClick() {
  router.push('/users')
}
function deploymentClick(item: any) {
  console.log(item)
  //showInfo(`部署被点击:${JSON.stringify(item)}`)
  modelDeploy.value.deployTitle = `部署服务- ${item.name} ${item.service_name}`
  modelDeploy.value.model = item
  deployModalOpen.value = true
}

function selectFolder() {
  
}


function confirmReboot(serviceId: any) {
  apiFetch('/api/service/action', {
    method: 'POST',
    body: {
      service_id: serviceId,
      action_type: 0
    },

  }).then((resp) => {
    if (resp.code == 200) {
      msgSuccess('重启成功')
    }
  })
}
function confirmClose(serviceId: any) {
  apiFetch('/api/service/action', {
    method: 'POST',
    body: {
      service_id: serviceId,
      action_type: 1
    },
 
  }).then((resp) => {
    if (resp.code == 200) {
       msgSuccess('关闭成功')
    }
  })
}

function selServerPath() {
  if (model.value.real_path) {
    viewFolder(model.value.real_path, true)
  } else {
    viewFolder(".", true)
  }
  
}
function selPath() {
  model.value.realPath = selectedPath.value
  viewFolderOpen.value = false
}
function selServerIIS() {
  viewIIS()
}
function selIIS() {
  //将IIS信息填写到当前model
  model.value.name = selectedIIS.value.Name
  model.value.service_name = selectedIIS.value.Name
  model.value.arg = selectedIIS.value.ApplicationPool
  model.value.real_path = selectedIIS.value.PhysicalPath
  model.value.listen_url = selectedIIS.value.binding
  model.value.service_type = 0
  console.log(model.value)
  viewIISListOpen.value = false
}

function viewFolder(path: string, isSel: boolean = false) {
  isCanSelect.value = !isSel
  if (isSel) {
    viewFolderTitle.value = "选择目录"
  } else {
    viewFolderTitle.value = "目录查看"
  }

  apiFetch('/api/service/dir_tree', {
    method: 'GET',
    query: {
      home: path
    },
  }).then((resp) => {
    curPath.value = path
    treeData.value = resp.data.map((r: any) => {
      return {
        title: r.path,
        key: r.path,
        pathType: r.info.is_dir,
        isLeaf: !r.info.is_dir
      }
    })
    viewFolderOpen.value = true
  })
}
function viewIIS() {
  isCanSelect.value = false
  apiFetch('/api/service/IISServiceList', {
    method: 'GET'
  }).then((resp) => {
    iisList.value = resp.data.list
    viewIISListOpen.value = true
  })
}
function iisListClick(item: {}) {
  selectedIIS.value = item
  selBtnDisabled.value = false
}
function confirmDelService(serviceId: any) {
  apiFetch('/api/service/delete', {
    method: 'POST',

    body: {
      ids: [serviceId],
    },

  }).then((resp) => {
    if (resp.code == 200) {
      msgSuccess('删除成功')
      refreshData()
    } else {
      showError(resp.message)
    }
  })
}
</script>
