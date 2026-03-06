package advertisement

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/advertisement"
	advertisementReq "github.com/flipped-aurora/gin-vue-admin/server/model/advertisement/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdvertisementService struct{}

// CreateAdvertisement 创建广告
func (adService *AdvertisementService) CreateAdvertisement(req advertisementReq.CreateAdvertisementRequest) error {
	ad := advertisement.Advertisement{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Link:        req.Link,
		Position:    req.Position,
		Status:      req.Status,
		Sort:        req.Sort,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		DeviceType:  req.DeviceType,
		ForcePopup:  req.ForcePopup,
	}

	if err := global.GVA_DB.Create(&ad).Error; err != nil {
		global.GVA_LOG.Error("创建广告失败", zap.Error(err))
		return errors.New("创建广告失败")
	}
	return nil
}

// DeleteAdvertisement 删除广告
func (adService *AdvertisementService) DeleteAdvertisement(id uint) error {
	if err := global.GVA_DB.Where("id = ?", id).Delete(&advertisement.Advertisement{}).Error; err != nil {
		global.GVA_LOG.Error("删除广告失败", zap.Error(err))
		return errors.New("删除广告失败")
	}
	return nil
}

// UpdateAdvertisement 更新广告
func (adService *AdvertisementService) UpdateAdvertisement(req advertisementReq.UpdateAdvertisementRequest) error {
	var ad advertisement.Advertisement
	if err := global.GVA_DB.First(&ad, req.ID).Error; err != nil {
		return errors.New("广告不存在")
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.Link != "" {
		updates["link"] = req.Link
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.StartTime != "" {
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		updates["end_time"] = req.EndTime
	}
	if req.DeviceType != "" {
		updates["device_type"] = req.DeviceType
	}
	if req.ForcePopup != nil {
		updates["force_popup"] = *req.ForcePopup
	}

	if err := global.GVA_DB.Model(&ad).Updates(updates).Error; err != nil {
		global.GVA_LOG.Error("更新广告失败", zap.Error(err))
		return errors.New("更新广告失败")
	}
	return nil
}

// GetAdvertisement 获取单个广告
func (adService *AdvertisementService) GetAdvertisement(id uint) (advertisement.Advertisement, error) {
	var ad advertisement.Advertisement
	if err := global.GVA_DB.First(&ad, id).Error; err != nil {
		return ad, errors.New("广告不存在")
	}
	return ad, nil
}

// GetAdvertisementList 获取广告列表
func (adService *AdvertisementService) GetAdvertisementList(info advertisementReq.AdvertisementSearch) (list []advertisement.Advertisement, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	
	db := global.GVA_DB.Model(&advertisement.Advertisement{})
	
	// 条件查询
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.Position != "" {
		db = db.Where("position = ?", info.Position)
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	if info.DeviceType != "" {
		db = db.Where("device_type = ?", info.DeviceType)
	}
	
	// 排序
	db = db.Order("sort DESC, created_at DESC")
	
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return
}

// GetSuitableAds 获取适合的广告（包含强制弹窗广告和指定位置的广告）
func (adService *AdvertisementService) GetSuitableAds(deviceID, position string, count int) ([]advertisement.Advertisement, error) {
	var ads []advertisement.Advertisement
	
	now := time.Now().Format("2006-01-02 15:04:05")
	
	// 查询条件：启用状态 + (强制弹窗 OR 指定位置) + 时间范围 + 设备类型
	db := global.GVA_DB.Model(&advertisement.Advertisement{}).
		Where("status = ?", 1). // 启用状态
		Where("(force_popup = ? OR position = ?)", true, position). // 强制弹窗或指定位置
		Where("(start_time = '' OR start_time <= ?)", now).
		Where("(end_time = '' OR end_time >= ?)", now)
	
	// 设备类型过滤
	db = db.Where("device_type IN (?, 'all')", getDeviceType(deviceID))
	
	// 排序：强制弹窗优先，然后按排序值降序，创建时间降序
	db = db.Order("force_popup DESC, sort DESC, created_at DESC")
	
	// 限制数量（强制弹窗不计入数量限制）
	if count > 0 {
		db = db.Limit(count + 10) // 多查询一些，前端会过滤
	}
	
	if err := db.Find(&ads).Error; err != nil {
		return nil, err
	}
	
	return ads, nil
}

// RecordAdView 记录广告展示
func (adService *AdvertisementService) RecordAdView(adID uint) error {
	return global.GVA_DB.Model(&advertisement.Advertisement{}).
		Where("id = ?", adID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// RecordAdClick 记录广告点击
func (adService *AdvertisementService) RecordAdClick(adID uint) error {
	return global.GVA_DB.Model(&advertisement.Advertisement{}).
		Where("id = ?", adID).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error
}

// GetAdvertisementStats 获取广告统计
func (adService *AdvertisementService) GetAdvertisementStats() (map[string]interface{}, error) {
	var stats struct {
		TotalAds    int64
		ActiveAds   int64
		TotalClicks int64
		TotalViews  int64
	}
	
	// 总广告数
	global.GVA_DB.Model(&advertisement.Advertisement{}).Count(&stats.TotalAds)
	
	// 活跃广告数（启用状态且在有效期内）
	global.GVA_DB.Model(&advertisement.Advertisement{}).
		Where("status = ?", 1).
		Where("(start_time = '' OR start_time <= ?)", time.Now().Format("2006-01-02 15:04:05")).
		Where("(end_time = '' OR end_time >= ?)", time.Now().Format("2006-01-02 15:04:05")).
		Count(&stats.ActiveAds)
	
	// 总点击和展示
	global.GVA_DB.Model(&advertisement.Advertisement{}).
		Select("SUM(click_count) as total_clicks, SUM(view_count) as total_views").
		Scan(&stats)
	
	// 今日点击和展示
	var todayStats struct {
		TodayClicks int64
		TodayViews  int64
	}
	
	// 这里需要广告点击记录表，暂时返回0
	todayStats.TodayClicks = 0
	todayStats.TodayViews = 0
	
	result := map[string]interface{}{
		"totalAds":     stats.TotalAds,
		"activeAds":    stats.ActiveAds,
		"totalClicks":  stats.TotalClicks,
		"totalViews":   stats.TotalViews,
		"todayClicks":  todayStats.TodayClicks,
		"todayViews":   todayStats.TodayViews,
	}
	
	return result, nil
}

// getDeviceType 根据设备ID判断设备类型（简化实现）
func getDeviceType(deviceID string) string {
	// 这里可以根据实际需求实现设备类型判断
	// 暂时返回mobile作为默认值
	return "mobile"
}