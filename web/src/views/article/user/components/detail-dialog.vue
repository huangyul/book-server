<script setup lang="ts">
import { ref } from 'vue'
import LikeFilled from '@iconify-icons/ant-design/like-filled'
import EyeFilled from '@iconify-icons/ant-design/eye-filled'
import StartFilled from '@iconify-icons/ant-design/star-filled'
import { ElMessage } from 'element-plus'
import {
  cancelCollcetApi,
  collectApi,
  getPubDetail,
  likeApi,
} from '@/api/article'

interface Article {
  id: number
  title: string
  content: string
  readCnt: number
  likeCnt: number
  collectCnt: number
  authorName: string
  liked: boolean
  collected: boolean
}

const dialogVisible = ref(false)

const data = ref<Article>({
  id: 0,
  title: '',
  content: '',
  readCnt: 0,
  likeCnt: 0,
  collectCnt: 0,
  authorName: '',
  liked: false,
  collected: false,
})

async function init() {
  const res = await getPubDetail(data.value.id)
  data.value.title = res.title
  data.value.content = res.content
  data.value.readCnt = res.read_cnt
  data.value.likeCnt = res.like_cnt
  data.value.collectCnt = res.collect_cnt
  data.value.authorName = res.author_name
  data.value.liked = res.liked
  data.value.collected = res.collected
}

async function handleClose() {
  dialogVisible.value = false
}

async function handleOpen(id: number) {
  data.value.id = id
  await init()
  dialogVisible.value = true
}

async function handleLike() {
  await likeApi(data.value.id, !data.value.liked)
  ElMessage.success('点赞成功')
  init()
}

async function handleCollect() {
  if (data.value.collected) {
    await cancelCollcetApi(data.value.id)
    ElMessage.success('取消成功')
  }
  else {
    await collectApi(data.value.id)
    ElMessage.success('收藏成功')
  }

  init()
}

defineExpose({
  handleOpen,
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    title="文章详情"
    width="500"
    :before-close="handleClose"
  >
    <div>标题：{{ data.title }}</div>
    <div>作者: {{ data.authorName }}</div>
    <div>
      内容:
      <p>{{ data.content }}</p>
    </div>
    <div class="flex items-center">
      <div class="flex items-center mr-2">
        <IconifyIconOffline :icon="EyeFilled" />{{ data.readCnt }}
      </div>
      <div class="flex items-center mr-2">
        <IconifyIconOffline
          class="cursor-pointer"
          :class="{ 'text-red-500': data.liked }"
          :icon="LikeFilled"
          @click="handleLike"
        />{{ data.likeCnt }}
      </div>
      <div class="flex items-center mr-2">
        <IconifyIconOffline
          class="cursor-pointer"
          :class="{ 'text-red-500': data.collected }"
          :icon="StartFilled"
          @click="handleCollect"
        />{{ data.collectCnt }}
      </div>
    </div>
  </el-dialog>
</template>
