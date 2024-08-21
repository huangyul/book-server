<script setup lang="ts">
import LikeFilled from '@iconify-icons/ant-design/like-filled'
import EyeFilled from '@iconify-icons/ant-design/eye-filled'
import StartFilled from '@iconify-icons/ant-design/star-filled'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import DetailDialog from './detail-dialog.vue'
import type { Article } from '@/types/article'
import { publishArticleApi } from '@/api/article'

interface Props {
  data: Article
}

const props = withDefaults(defineProps<Props>(), {
  data: () => ({
    id: 0,
    title: '',
    content: '',
    authorID: 0,
    authorName: '',
    status: 0,
    createdAt: '',
    updatedAt: '',
  }),
})

const emits = defineEmits(['refresh'])

async function handlePublish() {
  await publishArticleApi(props.data.id)
  ElMessage.success('发布成功')
  emits('refresh')
}

const detailRef = ref<InstanceType<typeof DetailDialog> | null>(null)
function openDetial() {
  detailRef.value.handleOpen(props.data.id)
}
</script>

<template>
  <div>
    <div>
      <p>{{ data.authorName }}</p>
      <p>{{ data.createdAt }}</p>
      <p>{{ data.updatedAt }}</p>
    </div>
    <div>
      <div class="font-bold text-lg">
        {{ data.title }}
      </div>
      <div class="text-#303030">
        {{ data.content }}
      </div>
    </div>
    <div>
      <el-tag v-if="data.status === 2" type="success">
        已发布
      </el-tag>
      <el-tag v-if="data.status === 1" type="warning">
        未发布
      </el-tag>
      <el-button v-if="data.status === 1" type="success" @click="handlePublish">
        发布
      </el-button>
      <el-button v-if="data.status === 2" type="primary" @click="openDetial">
        查看
      </el-button>
    </div>
    <div class="flex items-center">
      <div class="flex items-center mr-2">
        <IconifyIconOffline :icon="EyeFilled" />123
      </div>
      <div class="flex items-center mr-2">
        <IconifyIconOffline class="cursor-pointer" :icon="LikeFilled" />123
      </div>
      <div class="flex items-center mr-2">
        <IconifyIconOffline class="cursor-pointer" :icon="StartFilled" />123
      </div>
    </div>

    <DetailDialog ref="detailRef" />
  </div>
</template>
