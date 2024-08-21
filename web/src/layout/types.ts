import type { IconifyIcon } from '@iconify/vue'

const { VITE_HIDE_HOME } = import.meta.env

export const routerArrays: Array<RouteConfigs>
  = VITE_HIDE_HOME === 'false'
    ? [
        {
          path: '/welcome',
          meta: {
            title: '首页',
            icon: 'ep:home-filled',
          },
        },
      ]
    : []

export interface routeMetaType {
  title?: string
  icon?: string | IconifyIcon
  showLink?: boolean
  savedPosition?: boolean
  auths?: Array<string>
}

export interface RouteConfigs {
  path?: string
  query?: object
  params?: object
  meta?: routeMetaType
  children?: RouteConfigs[]
  name?: string
}

export interface multiTagsType {
  tags: Array<RouteConfigs>
}

export interface tagsViewsType {
  icon: string | IconifyIcon
  text: string
  divided: boolean
  disabled: boolean
  show: boolean
}

export interface setType {
  sidebar: {
    opened: boolean
    withoutAnimation: boolean
    isClickCollapse: boolean
  }
  device: string
  fixedHeader: boolean
  classes: {
    hideSidebar: boolean
    openSidebar: boolean
    withoutAnimation: boolean
    mobile: boolean
  }
  hideTabs: boolean
}

export interface menuType {
  id?: number
  name?: string
  path?: string
  noShowingChildren?: boolean
  children?: menuType[]
  value: unknown
  meta?: {
    icon?: string
    title?: string
    rank?: number
    showParent?: boolean
    extraIcon?: string
  }
  showTooltip?: boolean
  parentId?: number
  pathList?: number[]
  redirect?: string
}

export interface themeColorsType {
  color: string
  themeColor: string
}

export interface scrollbarDomType extends HTMLElement {
  wrap?: {
    offsetWidth: number
  }
}
