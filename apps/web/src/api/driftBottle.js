import service from '@/utils/request'

// H5 接口

export const throwBottle = (data) => {
  return service({
    url: '/driftBottle/throw',
    method: 'post',
    data
  })
}

export const pickBottle = (params) => {
  return service({
    url: '/driftBottle/pick',
    method: 'get',
    params
  })
}

export const replyBottle = (data) => {
  return service({
    url: '/driftBottle/reply',
    method: 'post',
    data
  })
}

export const getMyBottles = (params) => {
  return service({
    url: '/driftBottle/myBottles',
    method: 'get',
    params
  })
}

export const getBottleDetail = (params) => {
  return service({
    url: '/driftBottle/detail',
    method: 'get',
    params
  })
}

// 后台管理接口

export const adminGetBottleList = (params) => {
  return service({
    url: '/driftBottle/admin/bottleList',
    method: 'get',
    params
  })
}

export const adminGetReplyList = (params) => {
  return service({
    url: '/driftBottle/admin/replyList',
    method: 'get',
    params
  })
}

export const adminDeleteBottle = (data) => {
  return service({
    url: '/driftBottle/admin/deleteBottle',
    method: 'delete',
    data
  })
}

export const adminDeleteReply = (data) => {
  return service({
    url: '/driftBottle/admin/deleteReply',
    method: 'delete',
    data
  })
}
