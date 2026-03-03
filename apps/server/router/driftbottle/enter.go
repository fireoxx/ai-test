package driftbottle

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	DriftBottleRouter
}

var driftBottleApi = api.ApiGroupApp.DriftBottleApiGroup.DriftBottleApi
