package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/advertisement"

// AdvertisementResponse 广告响应
type AdvertisementResponse struct {
	advertisement.Advertisement
}

// AdvertisementListResponse 广告列表响应
type AdvertisementListResponse struct {
	List  []advertisement.Advertisement `json:"list"`
	Total int64                         `json:"total"`
}

// SuitableAdResponse 适合广告响应
type SuitableAdResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Link        string `json:"link"`
	Position    string `json:"position"`
	ForcePopup  bool   `json:"forcePopup"`
}

// AdvertisementStatsResponse 广告统计响应
type AdvertisementStatsResponse struct {
	TotalAds     int64 `json:"totalAds"`
	ActiveAds    int64 `json:"activeAds"`
	TotalClicks  int64 `json:"totalClicks"`
	TotalViews   int64 `json:"totalViews"`
	TodayClicks  int64 `json:"todayClicks"`
	TodayViews   int64 `json:"todayViews"`
}