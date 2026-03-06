package advertisement

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// Advertisement 广告表
type Advertisement struct {
	global.GVA_MODEL
	Title       string `json:"title" form:"title" gorm:"column:title;comment:广告标题"`
	Description string `json:"description" form:"description" gorm:"column:description;type:text;comment:广告描述"`
	ImageURL    string `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:图片URL"`
	Link        string `json:"link" form:"link" gorm:"column:link;comment:跳转链接"`
	Position    string `json:"position" form:"position" gorm:"column:position;comment:广告位置(bottom:底部,top:顶部,popup:弹窗)"`
	Status      int    `json:"status" form:"status" gorm:"column:status;default:1;comment:状态 1启用 2禁用"`
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序"`
	StartTime   string `json:"startTime" form:"startTime" gorm:"column:start_time;comment:开始时间"`
	EndTime     string `json:"endTime" form:"endTime" gorm:"column:end_time;comment:结束时间"`
	ClickCount  int    `json:"clickCount" form:"clickCount" gorm:"column:click_count;default:0;comment:点击次数"`
	ViewCount   int    `json:"viewCount" form:"viewCount" gorm:"column:view_count;default:0;comment:展示次数"`
	DeviceType  string `json:"deviceType" form:"deviceType" gorm:"column:device_type;default:all;comment:设备类型(all:全部,mobile:手机,pc:电脑)"`
	ForcePopup  bool   `json:"forcePopup" form:"forcePopup" gorm:"column:force_popup;default:false;comment:是否强制弹窗显示"`
}

func (Advertisement) TableName() string {
	return "advertisements"
}