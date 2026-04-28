import request from '@/utils/request'

// 上传图片
export function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传 PDF
export function uploadPDF(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/upload/pdf',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传视频
export function uploadVideo(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/upload/video',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取上传列表
export function getUploadList(page = 1, pageSize = 20, type = '') {
  return request({
    url: '/uploads',
    method: 'get',
    params: { page, page_size: pageSize, type }
  })
}

// 删除上传文件
export function deleteUpload(id) {
  return request({
    url: `/uploads/${id}`,
    method: 'delete'
  })
}
