package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/advertisement"
	"github.com/flipped-aurora/gin-vue-admin/server/service/driftbottle"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup       system.ServiceGroup
	ExampleServiceGroup      example.ServiceGroup
	DriftBottleServiceGroup  driftbottle.ServiceGroup
	AdvertisementServiceGroup advertisement.ServiceGroup
}
