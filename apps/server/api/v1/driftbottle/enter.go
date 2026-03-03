package driftbottle

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	DriftBottleApi
}

var driftBottleService = service.ServiceGroupApp.DriftBottleServiceGroup.DriftBottleService
