package driftbottle

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// DriftBottleReply 漂流瓶回复
type DriftBottleReply struct {
	global.GVA_MODEL
	BottleID uint        `json:"bottleId" form:"bottleId" gorm:"column:bottle_id;comment:漂流瓶ID"`
	Nickname string      `json:"nickname" form:"nickname" gorm:"column:nickname;comment:回复人昵称"`
	Content  string      `json:"content" form:"content" gorm:"column:content;type:text;comment:回复内容"`
	DeviceID string      `json:"deviceId" form:"deviceId" gorm:"column:device_id;comment:设备ID"`
	Bottle   DriftBottle `json:"bottle" gorm:"foreignKey:BottleID"`
}

func (DriftBottleReply) TableName() string {
	return "drift_bottle_replies"
}
