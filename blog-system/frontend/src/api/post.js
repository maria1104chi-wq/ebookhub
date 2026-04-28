import request from '@/utils/request'

// 获取博文列表
export function getPostList(page = 1, pageSize = 10) {
  return request({
    url: '/posts',
    method: 'get',
    params: { page, page_size: pageSize }
  })
}

// 获取博文详情
export function getPost(slug) {
  return request({
    url: `/posts/${slug}`,
    method: 'get'
  })
}

// 创建博文
export function createPost(data) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}

// 更新博文
export function updatePost(id, data) {
  return request({
    url: `/posts/${id}`,
    method: 'put',
    data
  })
}

// 删除博文
export function deletePost(id) {
  return request({
    url: `/posts/${id}`,
    method: 'delete'
  })
}

// 搜索博文
export function searchPosts(keyword, page = 1, pageSize = 10) {
  return request({
    url: '/posts/search',
    method: 'get',
    params: { keyword, page, page_size: pageSize }
  })
}

// 获取热门博文
export function getHotPosts(limit = 10) {
  return request({
    url: '/posts/hot',
    method: 'get',
    params: { limit }
  })
}

// 点赞博文
export function likePost(id) {
  return request({
    url: `/posts/${id}/like`,
    method: 'post'
  })
}

// 分享博文
export function sharePost(id) {
  return request({
    url: `/posts/${id}/share`,
    method: 'post'
  })
}
