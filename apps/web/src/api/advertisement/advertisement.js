import request from '@/utils/request'

// 获取广告列表
export function getAdvertisementList(params) {
  return request({
    url: '/advertisement/list',
    method: 'get',
    params
  })
}

// 创建广告
export function createAdvertisement(data) {
  return request({
    url: '/advertisement/create',
    method: 'post',
    data
  })
}

// 更新广告
export function updateAdvertisement(data) {
  return request({
    url: '/advertisement/update',
    method: 'put',
    data
  })
}

// 删除广告
export function deleteAdvertisement(data) {
  return request({
    url: '/advertisement/delete',
    method: 'delete',
    data
  })
}

// 获取广告详情
export function getAdvertisement(params) {
  return request({
    url: '/advertisement/detail',
    method: 'get',
    params
  })
}

// 获取广告统计
export function getAdvertisementStats() {
  return request({
    url: '/advertisement/stats',
    method: 'get'
  })
}

// 获取适合广告（H5接口）
export function getSuitableAds(params) {
  return request({
    url: '/advertisement/suitable',
    method: 'get',
    params
  })
}

// 记录广告展示（H5接口）
export function recordAdView(data) {
  return request({
    url: '/advertisement/view',
    method: 'post',
    data
  })
}

// 记录广告点击（H5接口）
export function recordAdClick(data) {
  return request({
    url: '/advertisement/click',
    method: 'post',
    data
  })
}