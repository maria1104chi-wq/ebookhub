import request from '@/utils/request'

// 获取仪表盘数据
export function getDashboard() {
  return request({
    url: '/admin/dashboard',
    method: 'get'
  })
}

// 获取敏感词列表
export function getSensitiveList(page = 1, pageSize = 20) {
  return request({
    url: '/admin/sensitive',
    method: 'get',
    params: { page, page_size: pageSize }
  })
}

// 添加敏感词
export function addSensitiveWord(word, category = 'default', level = 1) {
  return request({
    url: '/admin/sensitive',
    method: 'post',
    data: { word, category, level }
  })
}

// 删除敏感词
export function deleteSensitiveWord(id) {
  return request({
    url: `/admin/sensitive/${id}`,
    method: 'delete'
  })
}

// 更新敏感词状态
export function updateSensitiveStatus(id, status) {
  return request({
    url: '/admin/sensitive/status',
    method: 'put',
    data: { id, status }
  })
}

// 导入敏感词
export function importSensitiveWords() {
  return request({
    url: '/admin/sensitive/import',
    method: 'post'
  })
}

// 获取统计数据
export function getStatistics(startDate = '', endDate = '') {
  return request({
    url: '/admin/statistics',
    method: 'get',
    params: { start_date: startDate, end_date: endDate }
  })
}

// 备份数据库
export function backupDatabase() {
  return request({
    url: '/admin/backup',
    method: 'post'
  })
}
