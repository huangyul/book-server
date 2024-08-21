<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { editArticleApi, getDetailByAuthorApi } from '@/api/article'

// withDefaults(defineProps<{
//   type?: 'create' | 'view'
// }>(), {
//   type: 'create',
// })

const emits = defineEmits(['refresh'])

const type = ref<'create' | 'edit'>('create')

const formRef = ref<FormInstance>()

const form = ref<{
  id?: number
  title: string
  content: string
}>({
  title: '',
  content: '',
})

const rules = ref({
  title: [{ required: true, message: 'Please input title', trigger: 'blur' }],
  content: [
    { required: true, message: 'Please input content', trigger: 'blur' },
  ],
})

async function onSubmit() {
  if (!formRef.value)
    return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      await editArticleApi(form.value)
      ElMessage.success('编辑成功')
      form.value.content = ''
      form.value.title = ''
      emits('refresh')
      handleClose()
    }
    else {
      ElMessage.warning('请按要求输入')
    }
  })
}
const dialogVisible = ref(false)

async function init(id: number) {
  const res = await getDetailByAuthorApi(id)
  form.value.id = res.id
  form.value.content = res.content
  form.value.title = res.title
}

function handleClose() {
  dialogVisible.value = false
}

async function handleOpen(id?: number) {
  type.value = id ? 'edit' : 'create'

  if (id) {
    await init(id)
  }

  dialogVisible.value = true
}

defineExpose({
  handleOpen,
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    :title="type === 'create' ? '新建' : '编辑'"
    width="500"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="auto"
      style="max-width: 600px"
    >
      <el-form-item label="title" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="content" prop="content">
        <el-input v-model="form.content" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">
          Cancel
        </el-button>
        <el-button type="primary" @click="onSubmit">
          Confirm
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>
