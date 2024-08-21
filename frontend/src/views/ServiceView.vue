<template>
  <a-page-header
    style="border: 0px solid rgb(235, 237, 240);color: azure;"
    title="服务管理"
    sub-title="管理服务"
    @back="() => router.push('/home')"
  >
  </a-page-header>
  <a-flex vertical gap="middle">
    <a-flex horizontal gap="middle">
      <a-button type="primary" @click="newServiceClick">新服务</a-button>
      <!-- <a-checkbox v-model:checked="onlyChecked" @change="onlyChange">仅显示未审核</a-checkbox> -->
    </a-flex>
    <a-table :dataSource="data" :columns="columns" style="width: 1800px">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'service_type'">
          <a-tag>{{ formatServiceType(record.service_type) }}</a-tag>
        </template>
        <template v-if="column.key === 'CreatedAt'">
          <span>{{ format(record.CreatedAt) }}</span>
        </template>
        <template v-if="column.key === 'UpdatedAt'">
          <span>{{ format(record.UpdatedAt) }}</span>
        </template>
        <template v-if="column.key === 'last_deploy_time'">
          <span>{{ format(record.last_deploy_time) }}</span>
        </template>
        <!-- <template v-if="column.key === 'action'">
          <span>
            <a>修改</a>
            <a v-if="record.Enabled === true">停用</a>

            <a>删除</a>
            <a v-if="record.Enabled === false">审核</a>
          </span>
        </template> -->
      </template>
    </a-table>
  </a-flex>

  <a-modal v-model:open="newServiceModel" title="新增服务" @ok="addService" okText="确认" cancelText="取消">
    <!-- <a-flex horizontal gap="middle"> -->
      <a-flex vertical gap="middle">
        <a-input v-model:value="model.name" size="middle" placeholder="请输入服务名称"> </a-input>
        <a-select v-model:value="model.service_type" size="middle" :options="serviceOptions" />

        <a-input v-model:value="model.service_name" placeholder="请输入服务名"> </a-input>

        <a-input v-model:value="model.arg" placeholder="请输入参数"> </a-input>
        <a-input v-model:value="model.real_path" placeholder="请输入路径"> </a-input>
        <a-input v-model:value="model.listen_url" placeholder="请输入监听地址"> </a-input>
        
      </a-flex>
      <!-- <a-flex vertical gap="middle">
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
      </a-flex> -->
    <!-- </a-flex> -->
  </a-modal>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import {apiFetch} from '@/utils/ofetch'
import type { SelectProps } from 'ant-design-vue';
import { msgError, showError, showSuccess } from '@/utils/notification';
const data = ref<ServiceModel[]>([])

const router = useRouter()

interface ServiceModel {
    name:string,
    real_path: string,
    service_name:string,
    service_type: number,
    arg:string,
    listen_url:string
}

const model = ref<ServiceModel>({
    name: '',
    real_path: '',
    service_name: '',
    service_type: 0,
    arg: '',
    listen_url: ''
})
const newServiceModel = ref(false)
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

const columns = [
  {
    title: '服务名称',
    dataIndex: 'name',
    width: '10%'
  },
  {
    title: '路径',
    dataIndex: 'real_path',
    width: '10%'
  },
  {
    title: '服务名',
    dataIndex: 'service_name',
    width: '15%'
  },
  {
    title: '类型',
    dataIndex: 'service_type',
    width: '10%',
    key: 'service_type'
  },
  {
    title: '加入时间',
    dataIndex: 'CreatedAt',
    width: '15%',
    key: 'CreatedAt'
  },
  {
    title: '修改时间',
    dataIndex: 'UpdatedAt',
    width: '15%',
    key: 'UpdatedAt'
  },
  {
    title: '上次部署',
    dataIndex: 'last_deploy_time',
    width: '15%',
    key: 'last_deploy_time'
  }
]

onMounted(() => {
  getServices()
})

function addService() {
    console.log(model.value)
    if (model.value.service_type===3 && !model.value.service_name.endsWith('.service')) {
        msgError('systemd服务名需要以.service结尾')
        model.value.service_name += '.service'
        return
    }
    if (model.value.name.trim()==='') {
        msgError('名称不能为空')
        return
    }
    if (model.value.service_name.trim()==='') {
        msgError('服务名不能为空')
        return
    }
    if (model.value.service_name.trim()==='') {
        msgError('路径不能为空')
        return
    }
    apiFetch('/api/service',{
        method: 'POST',
        body: model.value
    }).then(data=>{
        if (data.code==200) {
            showSuccess('新增成功')
            newServiceModel.value = false
            getServices()
        } else {
            showError(data.message)
        }
    })

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

function getServices() {
  apiFetch('/api/service/list', {
    query: { page: 1, size: 10 }
  }).then((resp) => {
    data.value = resp.data.list
  })
}

function format(inParam:string) {
    if (inParam==='') {
        return ''
    }
  return dayjs(inParam).format('YYYY-MM-DD HH:mm:ss')
}

function newServiceClick() {
    newServiceModel.value = true

}
</script>

<style></style>
