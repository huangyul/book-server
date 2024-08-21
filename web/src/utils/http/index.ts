import Axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type CustomParamsSerializer,
} from 'axios'
import { stringify } from 'qs'
import { ElMessage } from 'element-plus'
import NProgress from '../progress'
import type {
  PureHttpError,
  PureHttpRequestConfig,
  PureHttpResponse,
  RequestMethods,
} from './types.d'
import { formatToken, getToken, removeToken } from '@/utils/auth'
import { useUserStoreHook } from '@/store/modules/user'
import router from '@/router'

// 相关配置请参考：www.axios-js.com/zh-cn/docs/#axios-request-config-1
const defaultConfig: AxiosRequestConfig = {
  // baseURL: '/api',
  // 请求超时时间
  timeout: 10000,
  headers: {
    'Accept': 'application/json, text/plain, */*',
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest',
  },
  // 数组格式参数序列化（https://github.com/axios/axios/issues/5142）
  paramsSerializer: {
    serialize: stringify as unknown as CustomParamsSerializer,
  },
}

class Http {
  constructor() {
    this.httpInterceptorsRequest()
    this.httpInterceptorsResponse()
  }

  /** `token`过期后，暂存待执行的请求 */
  private static requests = []

  /** 防止重复刷新`token` */
  private static isRefreshing = false

  /** 初始化配置对象 */
  private static initConfig: PureHttpRequestConfig = {}

  /** 保存当前`Axios`实例对象 */
  private static axiosInstance: AxiosInstance = Axios.create(defaultConfig)

  /** 重连原始请求 */
  private static retryOriginalRequest(config: PureHttpRequestConfig) {
    return new Promise((resolve) => {
      Http.requests.push((token: string) => {
        config.headers.Authorization = formatToken(token)
        resolve(config)
      })
    })
  }

  /** 请求拦截 */
  private httpInterceptorsRequest(): void {
    Http.axiosInstance.interceptors.request.use(
      async (config: PureHttpRequestConfig): Promise<any> => {
        // 开启进度条动画
        NProgress.start()
        // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
        if (typeof config.beforeRequestCallback === 'function') {
          config.beforeRequestCallback(config)
          return config
        }
        if (Http.initConfig.beforeRequestCallback) {
          Http.initConfig.beforeRequestCallback(config)
          return config
        }
        /** 请求白名单，放置一些不需要`token`的接口（通过设置请求白名单，防止`token`过期后再请求造成的死循环问题） */
        const whiteList = ['/refresh-token', '/login']
        return whiteList.some(url => config.url.endsWith(url))
          ? config
          : new Promise((resolve) => {
            const data = getToken()
            if (data) {
              const now = new Date().getTime()
              const expired = Number.parseInt(data.expires) - now <= 0
              if (expired) {
                if (!Http.isRefreshing) {
                  Http.isRefreshing = true
                  // token过期刷新
                  useUserStoreHook()
                    .handRefreshToken({ refreshToken: data.refreshToken })
                    .then((res) => {
                      const token = res.data.accessToken
                      config.headers.Authorization = formatToken(token)
                      Http.requests.forEach(cb => cb(token))
                      Http.requests = []
                    })
                    .finally(() => {
                      Http.isRefreshing = false
                    })
                }
                resolve(Http.retryOriginalRequest(config))
              }
              else {
                config.headers.Authorization = formatToken(
                  data.accessToken,
                )
                resolve(config)
              }
            }
            else {
              resolve(config)
            }
          })
      },
      (error) => {
        return Promise.reject(error)
      },
    )
  }

  /** 响应拦截 */
  private httpInterceptorsResponse(): void {
    const instance = Http.axiosInstance
    instance.interceptors.response.use(
      (response: PureHttpResponse) => {
        const $config = response.config
        // 关闭进度条动画
        NProgress.done()
        // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
        if (typeof $config.beforeResponseCallback === 'function') {
          $config.beforeResponseCallback(response)
          return response.data
        }
        if (Http.initConfig.beforeResponseCallback) {
          Http.initConfig.beforeResponseCallback(response)
          return response.data
        }
        return response.data
      },
      (error: PureHttpError) => {
        const $error = error
        $error.isCancelRequest = Axios.isCancel($error)
        // 关闭进度条动画
        NProgress.done()
        // 所有的响应异常 区分来源为取消请求/非取消请求
        return Promise.reject($error)
      },
    )
  }

  /** 通用请求工具函数 */
  public request<T>(
    method: RequestMethods,
    url: string,
    param?: AxiosRequestConfig,
    axiosConfig?: PureHttpRequestConfig,
  ): Promise<T> {
    const config = {
      method,
      url,
      ...param,
      ...axiosConfig,
    } as PureHttpRequestConfig

    // 单独处理自定义请求/响应回调
    return new Promise((resolve, reject) => {
      Http.axiosInstance
        .request(config)
        .then((response: undefined) => {
          resolve(response)
        })
        .catch((error) => {
          reject(error)
        })
    })
  }

  /** 通用请求工具函数 */
  public httpRequest<T>(
    method: RequestMethods,
    url: string,
    param?: AxiosRequestConfig,
    axiosConfig?: PureHttpRequestConfig,
  ): Promise<T> {
    const config = {
      method,
      url,
      ...param,
      ...axiosConfig,
    } as PureHttpRequestConfig

    // 单独处理自定义请求/响应回调
    return new Promise((resolve, reject) => {
      Http.axiosInstance
        .request(config)
        .then((res: any) => {
          if (res.code !== 0) {
            ElMessage.error(res.msg)
            reject(res.data)
          }
          else {
            resolve(res.data)
          }
        })
        .catch((error) => {
          const { response } = error
          if (response.status === 401) {
            ElMessage.error('登录已过期')
            removeToken()
            router.replace('/login')
            return
          }
          ElMessage.error('接口报错')
          reject(error)
        })
    })
  }

  /** 单独抽离的`post`工具函数 */
  public post<T, P>(
    url: string,
    params?: AxiosRequestConfig<P>,
    config?: PureHttpRequestConfig,
  ): Promise<T> {
    return this.httpRequest<T>('post', url, params, config)
  }

  /** 单独抽离的`get`工具函数 */
  public get<T, P>(
    url: string,
    params?: AxiosRequestConfig<P>,
    config?: PureHttpRequestConfig,
  ): Promise<T> {
    return this.httpRequest<T>('get', url, params, config)
  }
}

export const http = new Http()
