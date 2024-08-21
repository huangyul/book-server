<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Edit from './components/edit.vue'
import { deleteApi, getListApi } from '@/api/article'

const tableData = ref([{}])

const pageSize = ref(10)
const pageIndex = ref(1)
const total = ref(0)

async function doSearch() {
  const res = await getListApi({
    page_index: pageIndex.value,
    page_size: pageSize.value,
  })
  tableData.value = res.data
  total.value = res.total
}

function handleDelete(id: number) {
  ElMessageBox.confirm('确定要删除文章吗？')
    .then(async () => {
      await deleteApi(id)
      ElMessage.success()
      doSearch()
    })
}

const editRef = ref<InstanceType<typeof Edit> | null>(null)
function openDialog(id?: number) {
  editRef.value.handleOpen(id)
}

onMounted(() => {
  doSearch()
})
</script>

<template>
  <div>
    <el-card>
      <el-button @click="openDialog(null)">
        新增
      </el-button>
    </el-card>
    <el-card class="mt-2">
      <div class="h-full">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column prop="title" label="title" width="180" />
          <el-table-column prop="content" label="content" />
          <el-table-column label="created at" prop="created_at" />
          <el-table-column label="updated at" prop="updated_at" />
          <el-table-column label="action">
            <template #default="{ row }">
              <el-button @click="openDialog(row.id)">
                编辑
              </el-button>
              <el-button type="danger" @click="handleDelete(row.id)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <el-card class="mt-2 flex justify-end">
      <el-pagination
        v-model:current-page="pageIndex" v-model:page-size="pageSize" bcckground
        :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="doSearch" @current-change="doSearch"
      />
    </el-card>

    <Edit ref="editRef" @refresh="doSearch" />
  </div>
</template>
