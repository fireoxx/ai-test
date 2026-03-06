package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/advertisement"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// AdvertisementSearch 广告搜索参数
type AdvertisementSearch struct {
	advertisement.Advertisement
	request.PageInfo
}

// CreateAdvertisementRequest 创建广告请求
type CreateAdvertisementRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"imageUrl"`
	Link        string `json:"link" binding:"required,url"`
	Position    string `json:"position" binding:"required,oneof=bottom top popup"`
	Status      int    `json:"status" binding:"oneof=1 2"`
	Sort        int    `json:"sort"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	DeviceType  string `json:"deviceType" binding:"oneof=all mobile pc"`
	ForcePopup  bool   `json:"forcePopup"`
}

// UpdateAdvertisementRequest 更新广告请求
type UpdateAdvertisementRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Link        string `json:"link"`
	Position    string `json:"position" binding:"omitempty,oneof=bottom top popup"`
	Status      int    `json:"status" binding:"omitempty,oneof=1 2"`
	Sort        int    `json:"sort"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	DeviceType  string `json:"deviceType" binding:"omitempty,oneof=all mobile pc"`
	ForcePopup  *bool  `json:"forcePopup"`
}

// GetSuitableAdRequest 获取适合广告请求
type GetSuitableAdRequest struct {
	DeviceID string `json:"deviceId" form:"deviceId" binding:"required"`
	Position string `json:"position" form:"position" binding:"required,oneof=bottom top popup"`
	Count    int    `json:"count" form:"count"`
}

// RecordAdViewRequest 记录广告展示请求
type RecordAdViewRequest struct {
	AdID     uint   `json:"adId" binding:"required"`
	DeviceID string `json:"deviceId" binding:"required"`
	Position string `json:"position" binding:"required"`
}

// RecordAdClickRequest 记录广告点击请求
type RecordAdClickRequest struct {
	AdID     uint   `json:"adId" binding:"required"`
	DeviceID string `json:"deviceId" binding:"required"`
	Position string `json:"position" binding:"required"`
}