package source

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/advertisement"
	"gorm.io/gorm"
)

var Advertisement = new(advertisementSource)

type advertisementSource struct{}

func (a *advertisementSource) Initialize() error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 检查表是否存在
		if tx.Migrator().HasTable(&advertisement.Advertisement{}) {
			// 清空表数据
			if err := tx.Where("1 = 1").Delete(&advertisement.Advertisement{}).Error; err != nil {
				return err
			}
		} else {
			// 创建表
			if err := tx.AutoMigrate(&advertisement.Advertisement{}); err != nil {
				return err
			}
		}
		
		// 插入初始广告数据
		ads := []advertisement.Advertisement{
			{
				Title:       "探索更多有趣应用",
				Description: "发现更多创意H5应用，体验不一样的互动乐趣",
				ImageURL:    "",
				Link:        "https://example.com/more-apps",
				Position:    "bottom",
				Status:      1,
				Sort:        100,
				StartTime:   "",
				EndTime:     "",
				DeviceType:  "all",
			},
			{
				Title:       "AI创意工具推荐",
				Description: "使用AI工具提升创作效率，激发无限创意",
				ImageURL:    "",
				Link:        "https://example.com/ai-tools",
				Position:    "bottom",
				Status:      1,
				Sort:        90,
				StartTime:   "",
				EndTime:     "",
				DeviceType:  "all",
			},
			{
				Title:       "开发者学习资源",
				Description: "免费学习前端开发，掌握最新技术栈",
				ImageURL:    "",
				Link:        "https://example.com/learn-dev",
				Position:    "bottom",
				Status:      1,
				Sort:        80,
				StartTime:   "",
				EndTime:     "",
				DeviceType:  "all",
			},
			{
				Title:       "漂流瓶技巧分享",
				Description: "学习如何写出吸引人的漂流瓶内容",
				ImageURL:    "",
				Link:        "https://example.com/drift-tips",
				Position:    "bottom",
				Status:      1,
				Sort:        70,
				StartTime:   time.Now().Format("2006-01-02"),
				EndTime:     time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
				DeviceType:  "mobile",
			},
		}
		
		for i := range ads {
			if err := tx.Create(&ads[i]).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

func (a *advertisementSource) TableName() string {
	return "advertisements"
}

func (a *advertisementSource) Data() []advertisement.Advertisement {
	return []advertisement.Advertisement{}
}