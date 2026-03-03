package driftbottle

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	driftbottleReq "github.com/flipped-aurora/gin-vue-admin/server/model/driftbottle/request"
	driftbottleRes "github.com/flipped-aurora/gin-vue-admin/server/model/driftbottle/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DriftBottleApi struct{}

// ThrowBottle
// @Tags      DriftBottle
// @Summary   扔瓶子（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     data  body      driftbottleReq.ThrowBottleRequest  true  "投瓶信息"
// @Success   200   {object}  response.Response{msg=string}      "扔瓶子成功"
// @Router    /driftBottle/throw [post]
func (d *DriftBottleApi) ThrowBottle(c *gin.Context) {
	var req driftbottleReq.ThrowBottleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := driftBottleService.ThrowBottle(req); err != nil {
		global.GVA_LOG.Error("扔瓶子失败", zap.Error(err))
		response.FailWithMessage("扔瓶子失败", c)
		return
	}
	response.OkWithMessage("扔瓶子成功", c)
}

// PickBottle
// @Tags      DriftBottle
// @Summary   捞瓶子（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     deviceId  query     string                                                              true  "设备ID"
// @Success   200       {object}  response.Response{data=driftbottleRes.BottleDetailResponse,msg=string}  "捞瓶子成功"
// @Router    /driftBottle/pick [get]
func (d *DriftBottleApi) PickBottle(c *gin.Context) {
	deviceID := c.Query("deviceId")
	if deviceID == "" {
		response.FailWithMessage("deviceId不能为空", c)
		return
	}
	bottle, err := driftBottleService.PickRandomBottle(deviceID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(bottle, "捞瓶子成功", c)
}

// ReplyBottle
// @Tags      DriftBottle
// @Summary   回复瓶子（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     data  body      driftbottleReq.ReplyBottleRequest  true  "回复信息"
// @Success   200   {object}  response.Response{msg=string}      "回复成功"
// @Router    /driftBottle/reply [post]
func (d *DriftBottleApi) ReplyBottle(c *gin.Context) {
	var req driftbottleReq.ReplyBottleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := driftBottleService.ReplyBottle(req); err != nil {
		global.GVA_LOG.Error("回复瓶子失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("回复成功", c)
}

// GetMyBottles
// @Tags      DriftBottle
// @Summary   获取我的瓶子列表（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     deviceId  query     string                                                true  "设备ID"
// @Param     page      query     int                                                   true  "页码"
// @Param     pageSize  query     int                                                   true  "每页数量"
// @Success   200       {object}  response.Response{data=response.PageResult,msg=string}  "获取成功"
// @Router    /driftBottle/myBottles [get]
func (d *DriftBottleApi) GetMyBottles(c *gin.Context) {
	var req driftbottleReq.MyBottlesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := driftBottleService.GetMyBottles(req.DeviceID, req.PageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取我的瓶子失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetBottleDetail
// @Tags      DriftBottle
// @Summary   获取瓶子详情及回复（H5接口）
// @accept    application/json
// @Produce   application/json
// @Param     id  query     int                                                                        true  "瓶子ID"
// @Success   200 {object}  response.Response{data=driftbottleRes.BottleDetailResponse,msg=string}    "获取成功"
// @Router    /driftBottle/detail [get]
func (d *DriftBottleApi) GetBottleDetail(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bottle, replies, err := driftBottleService.GetBottleDetail(req.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取瓶子详情失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(driftbottleRes.BottleDetailResponse{
		Bottle:  bottle,
		Replies: replies,
	}, "获取成功", c)
}

// AdminGetBottleList
// @Tags      DriftBottle
// @Summary   后台获取漂流瓶列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     driftbottleReq.BottleListRequest                           true  "分页参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}     "获取成功"
// @Router    /driftBottle/admin/bottleList [get]
func (d *DriftBottleApi) AdminGetBottleList(c *gin.Context) {
	var req driftbottleReq.BottleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := driftBottleService.AdminGetBottleList(req)
	if err != nil {
		global.GVA_LOG.Error("获取漂流瓶列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// AdminGetReplyList
// @Tags      DriftBottle
// @Summary   后台获取回复列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     driftbottleReq.ReplyListRequest                            true  "分页参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}     "获取成功"
// @Router    /driftBottle/admin/replyList [get]
func (d *DriftBottleApi) AdminGetReplyList(c *gin.Context) {
	var req driftbottleReq.ReplyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := driftBottleService.AdminGetReplyList(req)
	if err != nil {
		global.GVA_LOG.Error("获取回复列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// AdminDeleteBottle
// @Tags      DriftBottle
// @Summary   后台删除漂流瓶
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "漂流瓶ID"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /driftBottle/admin/deleteBottle [delete]
func (d *DriftBottleApi) AdminDeleteBottle(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := driftBottleService.AdminDeleteBottle(req.Uint()); err != nil {
		global.GVA_LOG.Error("删除漂流瓶失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// AdminDeleteReply
// @Tags      DriftBottle
// @Summary   后台删除回复
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "回复ID"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /driftBottle/admin/deleteReply [delete]
func (d *DriftBottleApi) AdminDeleteReply(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := driftBottleService.AdminDeleteReply(req.Uint()); err != nil {
		global.GVA_LOG.Error("删除回复失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
