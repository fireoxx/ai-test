package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// ThrowBottleRequest 扔瓶子请求
type ThrowBottleRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Content  string `json:"content" binding:"required,max=500"`
	DeviceID string `json:"deviceId" binding:"required"`
}

// ReplyBottleRequest 回复瓶子请求
type ReplyBottleRequest struct {
	BottleID uint   `json:"bottleId" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Content  string `json:"content" binding:"required,max=500"`
	DeviceID string `json:"deviceId" binding:"required"`
}

// PickBottleRequest 捞瓶子请求
type PickBottleRequest struct {
	DeviceID string `json:"deviceId" form:"deviceId" binding:"required"`
}

// MyBottlesRequest 我的瓶子请求
type MyBottlesRequest struct {
	request.PageInfo
	DeviceID string `json:"deviceId" form:"deviceId" binding:"required"`
}

// BottleListRequest 后台管理列表请求
type BottleListRequest struct {
	request.PageInfo
	Status int `json:"status" form:"status"`
}

// ReplyListRequest 后台回复列表请求
type ReplyListRequest struct {
	request.PageInfo
	BottleID uint `json:"bottleId" form:"bottleId"`
}
