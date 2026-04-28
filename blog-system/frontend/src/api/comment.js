import request from '@/utils/request'

// 获取评论列表
export function getComments(postId, page = 1, pageSize = 50) {
  return request({
    url: `/posts/${postId}/comments`,
    method: 'get',
    params: { page, page_size: pageSize }
  })
}

// 创建评论
export function createComment(data) {
  return request({
    url: '/comments',
    method: 'post',
    data
  })
}

// 点赞评论
export function likeComment(id) {
  return request({
    url: `/comments/${id}/like`,
    method: 'post'
  })
}

// 删除评论
export function deleteComment(id) {
  return request({
    url: `/comments/${id}`,
    method: 'delete'
  })
}

// 获取待审核评论
export function getPendingComments(page = 1, pageSize = 20) {
  return request({
    url: '/admin/comments/pending',
    method: 'get',
    params: { page, page_size: pageSize }
  })
}

// 审核评论
export function auditComment(id, status) {
  return request({
    url: '/admin/comments/audit',
    method: 'post',
    data: { id, status }
  })
}
