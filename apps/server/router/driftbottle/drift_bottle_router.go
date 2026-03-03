package driftbottle

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DriftBottleRouter struct{}

// InitDriftBottleRouter 注册漂流瓶路由
// H5公开接口挂 PublicGroup，后台管理接口挂 PrivateGroup（需JWT鉴权）
func (r *DriftBottleRouter) InitDriftBottleRouter(PrivateGroup *gin.RouterGroup, PublicGroup *gin.RouterGroup) {
	// H5 公开接口（无需登录）
	h5Router := PublicGroup.Group("driftBottle")
	{
		h5Router.POST("throw", driftBottleApi.ThrowBottle)      // 扔瓶子
		h5Router.GET("pick", driftBottleApi.PickBottle)         // 捞瓶子
		h5Router.POST("reply", driftBottleApi.ReplyBottle)      // 回复瓶子
		h5Router.GET("myBottles", driftBottleApi.GetMyBottles)  // 我的瓶子
		h5Router.GET("detail", driftBottleApi.GetBottleDetail)  // 瓶子详情
	}

	// 后台管理接口（需JWT鉴权）
	adminRouter := PrivateGroup.Group("driftBottle/admin").Use(middleware.OperationRecord())
	adminRouterWithoutRecord := PrivateGroup.Group("driftBottle/admin")
	{
		adminRouter.DELETE("deleteBottle", driftBottleApi.AdminDeleteBottle) // 删除漂流瓶
		adminRouter.DELETE("deleteReply", driftBottleApi.AdminDeleteReply)   // 删除回复
	}
	{
		adminRouterWithoutRecord.GET("bottleList", driftBottleApi.AdminGetBottleList) // 漂流瓶列表
		adminRouterWithoutRecord.GET("replyList", driftBottleApi.AdminGetReplyList)   // 回复列表
	}
}
