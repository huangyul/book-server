import { http } from '@/utils/http'

interface Result {
  success: boolean
  data: Array<any>
}

export function getAsyncRoutes() {
  return http.request<Result>('get', '/get-async-routes')
}
