export default {
  path: '/article',
  redirect: '/article/user/edit',
  meta: {
    icon: 'ri:information-line',
    // showLink: false,
    title: '文章',
    rank: 9,
  },
  children: [
    {
      path: '/article/user/edit',
      name: 'edit',
      component: () => import('@/views/article/user/index.vue'),
      meta: {
        title: '用户页',
      },
    },
    {
      path: '/article/admin/list',
      name: 'list',
      component: () => import('@/views/article/list/index.vue'),
      meta: {
        title: '文章管理',
      },
    },
  ],
} satisfies RouteConfigsTable
