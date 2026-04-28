import request from '@/utils/request'

// 发送短信验证码
export function sendSMSCode(phone) {
  return request({
    url: '/auth/send-sms',
    method: 'post',
    data: { phone }
  })
}

// 管理员登录
export function login(username, password, phone, code) {
  return request({
    url: '/auth/login',
    method: 'post',
    data: { username, password, phone, code }
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

// 登出
export function logout() {
  return request({
    url: '/user/logout',
    method: 'post'
  })
}

// 注册管理员
export function registerAdmin(username, password, phone) {
  return request({
    url: '/auth/register',
    method: 'post',
    data: { username, password, phone }
  })
}
