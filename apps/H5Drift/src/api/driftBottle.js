import service from './request'

export const throwBottle = (data) => service({ url: '/driftBottle/throw', method: 'post', data })

export const pickBottle = (params) => service({ url: '/driftBottle/pick', method: 'get', params })

export const replyBottle = (data) => service({ url: '/driftBottle/reply', method: 'post', data })

export const getMyBottles = (params) => service({ url: '/driftBottle/myBottles', method: 'get', params })

export const getBottleDetail = (params) => service({ url: '/driftBottle/detail', method: 'get', params })
