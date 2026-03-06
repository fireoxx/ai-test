package advertisement

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	advertisementReq "github.com/flipped-aurora/gin-vue-admin/server/model/advertisement/request"
	advertisementRes "github.com/flipped-aurora/gin-vue-admin/server/model/advertisement/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdvertisementApi struct{}

var advertisementService = service.ServiceGroupApp.AdvertisementServiceGroup.AdvertisementService

// CreateAdvertisement
// @Tags      Advertisement
// @Summary   创建广告
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      advertisementReq.CreateAdvertisementRequest  true  "广告信息"
// @Success   200   {object}  response.Response{msg=string}                "创建成功"
// @Router    /advertisement/create [post]
func (api *AdvertisementApi) CreateAdvertisement(c *gin.Context) {
	var req advertisementReq.CreateAdvertisementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	if err := advertisementService.CreateAdvertisement(req); err != nil {
		global.GVA_LOG.Error("创建广告失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	response.OkWithMessage("创建成功", c)
}

// DeleteAdvertisement
// @Tags      Advertisement
// @Summary   删除广告
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdInfo  true  "广告ID"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /advertisement/delete [delete]
func (api *AdvertisementApi) DeleteAdvertisement(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	if err := advertisementService.DeleteAdvertisement(uint(req.ID)); err != nil {
		global.GVA_LOG.Error("删除广告失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	response.OkWithMessage("删除成功", c)
}

// UpdateAdvertisement
// @Tags      Advertisement
// @Summary   更新广告
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      advertisementReq.UpdateAdvertisementRequest  true  "广告信息"
// @Success   200   {object}  response.Response{msg=string}                "更新成功"
// @Router    /advertisement/update [put]
func (api *AdvertisementApi) UpdateAdvertisement(c *gin.Context) {
	var req advertisementReq.UpdateAdvertisementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	if err := advertisementService.UpdateAdvertisement(req); err != nil {
		global.GVA_LOG.Error("更新广告失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	response.OkWithMessage("更新成功", c)
}

// GetAdvertisement
// @Tags      Advertisement
// @Summary   获取单个广告
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id  query     uint                                            true  "广告ID"
// @Success   200  {object}  response.Response{data=advertisementRes.AdvertisementResponse,msg=string}  "获取成功"
// @Router    /advertisement/detail [get]
func (api *AdvertisementApi) GetAdvertisement(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("id不能为空", c)
		return
	}
	
	var idInfo request.GetById
	if err := c.ShouldBindQuery(&idInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	ad, err := advertisementService.GetAdvertisement(uint(idInfo.ID))
	if err != nil {
		global.GVA_LOG.Error("获取广告失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	response.OkWithDetailed(advertisementRes.AdvertisementResponse{Advertisement: ad}, "获取成功", c)
}

// GetAdvertisementList
// @Tags      Advertisement
// @Summary   获取广告列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     advertisementReq.AdvertisementSearch  true  "搜索参数"
// @Success   200   {object}  response.Response{data=advertisementRes.AdvertisementListResponse,msg=string}  "获取成功"
// @Router    /advertisement/list [get]
func (api *AdvertisementApi) GetAdvertisementList(c *gin.Context) {
	var pageInfo advertisementReq.AdvertisementSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	list, total, err := advertisementService.GetAdvertisementList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取广告列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	
	response.OkWithDetailed(advertisementRes.AdvertisementListResponse{
		List:  list,
		Total: total,
	}, "获取成功", c)
}

// GetSuitableAds
// @Tags      Advertisement
// @Summary   获取适合的广告（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     deviceId  query     string                                                              true  "设备ID"
// @Param     position  query     string                                                              true  "广告位置"
// @Param     count     query     int                                                                 false "广告数量"
// @Success   200       {object}  response.Response{data=[]advertisementRes.SuitableAdResponse,msg=string}  "获取成功"
// @Router    /advertisement/suitable [get]
func (api *AdvertisementApi) GetSuitableAds(c *gin.Context) {
	var req advertisementReq.GetSuitableAdRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	// 默认获取3个广告
	if req.Count == 0 {
		req.Count = 3
	}
	
	ads, err := advertisementService.GetSuitableAds(req.DeviceID, req.Position, req.Count)
	if err != nil {
		global.GVA_LOG.Error("获取适合广告失败", zap.Error(err))
		response.FailWithMessage("获取广告失败", c)
		return
	}
	
	// 转换为响应格式
	var suitableAds []advertisementRes.SuitableAdResponse
	for _, ad := range ads {
		suitableAds = append(suitableAds, advertisementRes.SuitableAdResponse{
			ID:          ad.ID,
			Title:       ad.Title,
			Description: ad.Description,
			ImageURL:    ad.ImageURL,
			Link:        ad.Link,
			Position:    ad.Position,
			ForcePopup:  ad.ForcePopup,
		})
	}
	
	response.OkWithDetailed(suitableAds, "获取成功", c)
}

// RecordAdView
// @Tags      Advertisement
// @Summary   记录广告展示（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     data  body      advertisementReq.RecordAdViewRequest  true  "展示记录"
// @Success   200   {object}  response.Response{msg=string}         "记录成功"
// @Router    /advertisement/view [post]
func (api *AdvertisementApi) RecordAdView(c *gin.Context) {
	var req advertisementReq.RecordAdViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	if err := advertisementService.RecordAdView(req.AdID); err != nil {
		global.GVA_LOG.Error("记录广告展示失败", zap.Error(err))
		response.FailWithMessage("记录失败", c)
		return
	}
	
	response.OkWithMessage("记录成功", c)
}

// RecordAdClick
// @Tags      Advertisement
// @Summary   记录广告点击（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     data  body      advertisementReq.RecordAdClickRequest  true  "点击记录"
// @Success   200   {object}  response.Response{msg=string}          "记录成功"
// @Router    /advertisement/click [post]
func (api *AdvertisementApi) RecordAdClick(c *gin.Context) {
	var req advertisementReq.RecordAdClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	if err := advertisementService.RecordAdClick(req.AdID); err != nil {
		global.GVA_LOG.Error("记录广告点击失败", zap.Error(err))
		response.FailWithMessage("记录失败", c)
		return
	}
	
	response.OkWithMessage("记录成功", c)
}

// GetAdvertisementStats
// @Tags      Advertisement
// @Summary   获取广告统计
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=advertisementRes.AdvertisementStatsResponse,msg=string}  "获取成功"
// @Router    /advertisement/stats [get]
func (api *AdvertisementApi) GetAdvertisementStats(c *gin.Context) {
	stats, err := advertisementService.GetAdvertisementStats()
	if err != nil {
		global.GVA_LOG.Error("获取广告统计失败", zap.Error(err))
		response.FailWithMessage("获取统计失败", c)
		return
	}
	
	response.OkWithDetailed(stats, "获取成功", c)
}