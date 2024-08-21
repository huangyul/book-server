import type { ListResponse } from 'types/global'
import type { Article, CreateForm, ListParams } from './types'
import { http } from '@/utils/http'

export function editArticleApi(data: CreateForm) {
  return http.httpRequest('post', '/api/articles/edit', { data })
}

export function getListApi(data: ListParams) {
  return http.httpRequest<ListResponse>('post', '/api/articles/list', { data })
}

/**
 * 获取文章详情
 * @param id
 * @returns
 */
export function getDetailByAuthorApi(id: number) {
  return http.httpRequest<Article>('get', `/api/articles/detail/${id}`)
}

/**
 * 获取已发表文章详情
 * @param id
 * @returns
 */
export function getPubDetail(id: number) {
  return http.httpRequest<Article>('get', `/api/articles/pub/detail/${id}`)
}

/**
 * 删除文章
 * @param id
 * @returns
 */
export function deleteApi(id: number) {
  return http.httpRequest('get', `/api/articles/delete/${id}`)
}

export function getListByAuthor(data: ListParams) {
  return http.httpRequest<
    ListResponse<{
      id: number
      title: string
      content: string
      author_id: number
      author_name: string
      status: number
      created_at: string
      updated_at: string
    }>
  >('post', '/api/articles/list-by-author', { data })
}

// 发布文章
export function publishArticleApi(id: number) {
  return http.httpRequest('get', `/api/articles/publish/${id}`)
}

/**
 * 点赞接口
 * @param id
 * @param status true 点赞   false 取消点赞
 * @returns
 */
export function likeApi(id: number, status: boolean) {
  return http.httpRequest('post', '/api/articles/pub/like', {
    data: {
      id,
      status,
    },
  })
}

/**
 * 收藏文章接口
 * @param id 文章id
 * @returns
 */
export function collectApi(id: number) {
  return http.httpRequest('get', `/api/articles/pub/collect/${id}`)
}

/**
 * 取消收藏
 * @param id
 * @returns
 */
export function cancelCollcetApi(id: number) {
  return http.httpRequest('get', `/api/articles/pub/cancel-collect/${id}`)
}
