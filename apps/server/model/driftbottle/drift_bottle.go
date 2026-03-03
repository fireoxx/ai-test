package driftbottle

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// DriftBottle 漂流瓶
type DriftBottle struct {
	global.GVA_MODEL
	Nickname string `json:"nickname" form:"nickname" gorm:"column:nickname;comment:投瓶人昵称"`
	Content  string `json:"content" form:"content" gorm:"column:content;type:text;comment:瓶子内容"`
	Status   int    `json:"status" form:"status" gorm:"column:status;default:1;comment:状态 1漂流中 2已被捞起 3已回复"`
	DeviceID string `json:"deviceId" form:"deviceId" gorm:"column:device_id;comment:设备ID"`
}

func (DriftBottle) TableName() string {
	return "drift_bottles"
}
