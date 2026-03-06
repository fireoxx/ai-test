import service from './request'

/**
 * 获取广告列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getAdList = (params) => service({
  url: '/advertisement/list',
  method: 'get',
  params
})

/**
 * 获取单个广告
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getAdDetail = (params) => service({
  url: '/advertisement/detail',
  method: 'get',
  params
})

/**
 * 记录广告点击
 * @param {Object} data 点击数据
 * @returns {Promise}
 */
export const recordAdClick = (data) => {
  return service({
    url: '/advertisement/click',
    method: 'post',
    data
  })
}

/**
 * 记录广告展示
 * @param {Object} data 展示数据
 * @returns {Promise}
 */
export const recordAdView = (data) => {
  return service({
    url: '/advertisement/view',
    method: 'post',
    data
  })
}

/**
 * 获取适合当前用户的广告
 * @param {Object} params 用户参数
 * @returns {Promise}
 */
export const getSuitableAd = (params) => {
  return service({
    url: '/advertisement/suitable',
    method: 'get',
    params: {
      deviceId: params.deviceId,
      position: 'bottom', // 广告位置：bottom底部
      count: 3, // 默认获取3个广告用于轮播
      ...params
    }
  })
}