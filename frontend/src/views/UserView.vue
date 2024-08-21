<template>
    <a-page-header
        style="border: 0px solid rgb(235, 237, 240)"
        title="用户管理"
        sub-title="管理用户和审核"
        @back="() => router.push('/home')"
  >
</a-page-header>
  <a-flex vertical gap="middle">
    <a-flex horizontal gap="middle">
      <a-button type="primary" @click="newUserClick">新用户</a-button>
      <a-checkbox v-model:checked="onlyChecked" @change="onlyChange">仅显示未审核</a-checkbox>
    </a-flex>
    <a-table :dataSource="data" :columns="columns" style="width: 1800px">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'IsAdmin'">
          <a-checkbox :checked="record.IsAdmin" />
        </template>
        <template v-if="column.key === 'Enabled'">
          <a-checkbox :checked="record.Enabled" />
        </template>
        <template v-if="column.key === 'CreatedAt'">
          <span>{{ format(record.CreatedAt) }}</span>
        </template>
        <template v-if="column.key === 'UpdatedAt'">
          <span>{{ format(record.UpdatedAt) }}</span>
        </template>
        <template v-if="column.key === 'LastLoginTime'">
          <span>{{ format(record.LastLoginTime) }}</span>
        </template>
        <template v-if="column.key === 'action'">
          <span>
            <a>修改</a>
            <a v-if="record.Enabled === true">停用</a>
            
            <a>删除</a>
            <a v-if="record.Enabled === false">审核</a>

          </span>
        </template>
      </template>
    </a-table>
    
  </a-flex>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import {apiFetch} from '@/utils/ofetch'
import dayjs from 'dayjs'
const data = ref([])
const router = useRouter()
const onlyChecked = ref(false)


const columns = [
  {
    title: '用户名',
    dataIndex: 'Name',
    width: '10%'
  },
  {
    title: '邮箱',
    dataIndex: 'Email',
    width: '10%'
  },
  {
    title: '管理员',
    dataIndex: 'IsAdmin',
    width: '15%',
    key: 'IsAdmin'
  },
  {
    title: '启用',
    dataIndex: 'Enabled',
    width: '10%',
    key: 'Enabled'
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
    title: '上次登录',
    dataIndex: 'LastLoginTime',
    width: '15%',
    key: 'LastLoginTime'
  },
  {
    title: '操作',
    width:'30%',
    key: 'action',
  },
]

onMounted(() => {
    getUsers()
})

function getUsers() {
    apiFetch('/api/user/list', {
    query: { page:1,size: 10 }
  }).then((resp) => {
    data.value = resp.data.list
  })
}

function onlyChange() {
    getUsers()
}
function format(inParam:string) {
  return dayjs(inParam).format('YYYY-MM-DD HH:mm:ss')
}

function newUserClick() {

}
</script>

<style>

</style>