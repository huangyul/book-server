<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Article from './article.vue'
import EditDialog from './edit-dialog.vue'
import { getListByAuthor } from '@/api/article'

const pageIndex = ref(1)
const pageSize = ref(10)

const list = ref([])

const editRef = ref<InstanceType<typeof EditDialog> | null>(null)
function openDialog() {
  editRef.value.handleOpen()
}

async function handleSearch() {
  const res = await getListByAuthor({
    page_index: pageIndex.value,
    page_size: pageSize.value,
  })
  res.data.forEach((d) => {
    list.value.push({
      id: d.id,
      title: d.title,
      content: d.content,
      createdAt: d.created_at,
      status: d.status,
    })
  })
}

async function init() {
  pageIndex.value = 1
  list.value = []
  handleSearch()
}

onMounted(() => {
  init()
})
</script>

<template>
  <div>
    <el-card>
      <el-button type="primary" @click="openDialog">
        新增
      </el-button>
      <div v-for="data in list" :key="data.id">
        <Article :data @refresh="init()" />
        <el-divider />
      </div>
    </el-card>

    <EditDialog ref="editRef" @refresh="init()" />
  </div>
</template>
