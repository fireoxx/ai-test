package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/driftbottle"

// BottleDetailResponse 瓶子详情（含回复列表）
type BottleDetailResponse struct {
	Bottle  driftbottle.DriftBottle        `json:"bottle"`
	Replies []driftbottle.DriftBottleReply `json:"replies"`
}
