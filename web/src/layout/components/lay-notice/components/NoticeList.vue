<script setup lang="ts">
import type { PropType } from 'vue'
import type { ListItem } from '../data'
import NoticeItem from './NoticeItem.vue'

defineProps({
  list: {
    type: Array as PropType<Array<ListItem>>,
    default: () => [],
  },
  emptyText: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div v-if="list.length">
    <NoticeItem v-for="(item, index) in list" :key="index" :notice-item="item" />
  </div>
  <el-empty v-else :description="emptyText" />
</template>
