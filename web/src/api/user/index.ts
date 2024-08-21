import type { LoginEmailReq, LoginResponse } from './types'
import { http } from '@/utils/http'

export interface UserResult {
  success: boolean
  data: {
    /** 头像 */
    avatar: string
    /** 用户名 */
    username: string
    /** 昵称 */
    nickname: string
    /** 当前登录用户的角色 */
    roles: Array<string>
    /** `token` */
    accessToken: string
    /** 用于调用刷新`accessToken`的接口时所需的`token` */
    refreshToken: string
    /** `accessToken`的过期时间（格式'xxxx/xx/xx xx:xx:xx'） */
    expires: Date
  }
}

export interface RefreshTokenResult {
  success: boolean
  data: {
    /** `token` */
    accessToken: string
    /** 用于调用刷新`accessToken`的接口时所需的`token` */
    refreshToken: string
    /** `accessToken`的过期时间（格式'xxxx/xx/xx xx:xx:xx'） */
    expires: Date
  }
}

/** 登录 */
export function getLogin(data?: object) {
  return http.request<UserResult>('post', '/login', { data })
}

/** 刷新`token` */
export function refreshTokenApi(data?: object) {
  return http.request<RefreshTokenResult>('post', '/refresh-token', { data })
}

/* 使用邮箱登录 */
export function loginByEmailApi(data: LoginEmailReq) {
  return http.httpRequest<LoginResponse>('post', '/api/user/login', { data })
}
