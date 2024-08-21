export interface CreateForm {
  id?: number
  content: string
  title: string
}

export interface ListParams {
  page_index: number
  page_size: number
}

export interface Article {
  liked: boolean
  collected: boolean
  id: number
  title: string
  content: string
  author_id: number
  author_name: string
  status: number
  created_at: string
  updated_at: string
  read_cnt?: number
  like_cnt?: number
  collect_cnt?: number
}
