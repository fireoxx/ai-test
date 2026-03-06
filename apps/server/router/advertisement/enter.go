package advertisement

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	AdvertisementRouter
}

var advertisementApi = api.ApiGroupApp.AdvertisementApiGroup.AdvertisementApi
