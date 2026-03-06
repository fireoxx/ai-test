package advertisement

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AdvertisementRouter struct{}

func (r *AdvertisementRouter) InitAdvertisementRouter(PrivateGroup *gin.RouterGroup, PublicGroup *gin.RouterGroup) {
	// 管理后台接口（需要鉴权）
	adRouter := PrivateGroup.Group("advertisement").Use(middleware.OperationRecord())
	adRouterWithoutRecord := PrivateGroup.Group("advertisement")

	{
		adRouter.POST("create", advertisementApi.CreateAdvertisement)   // 创建广告
		adRouter.DELETE("delete", advertisementApi.DeleteAdvertisement) // 删除广告
		adRouter.PUT("update", advertisementApi.UpdateAdvertisement)    // 更新广告
	}

	{
		adRouterWithoutRecord.GET("detail", advertisementApi.GetAdvertisement)     // 获取单个广告
		adRouterWithoutRecord.GET("list", advertisementApi.GetAdvertisementList)   // 获取广告列表
		adRouterWithoutRecord.GET("stats", advertisementApi.GetAdvertisementStats) // 获取广告统计
	}

	// H5公开接口（无需鉴权）
	adPublicRouter := PublicGroup.Group("advertisement")
	{
		adPublicRouter.GET("suitable", advertisementApi.GetSuitableAds)    // 获取适合广告（H5使用）
		adPublicRouter.POST("view", advertisementApi.RecordAdView)         // 记录广告展示（H5使用）
		adPublicRouter.POST("click", advertisementApi.RecordAdClick)       // 记录广告点击（H5使用）
	}
}
