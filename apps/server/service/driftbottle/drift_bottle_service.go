package driftbottle

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/driftbottle"
	driftbottleReq "github.com/flipped-aurora/gin-vue-admin/server/model/driftbottle/request"
	"gorm.io/gorm"
)

type DriftBottleService struct{}

// ThrowBottle 扔瓶子
func (s *DriftBottleService) ThrowBottle(req driftbottleReq.ThrowBottleRequest) error {
	bottle := driftbottle.DriftBottle{
		Nickname: req.Nickname,
		Content:  req.Content,
		DeviceID: req.DeviceID,
		Status:   1,
	}
	return global.GVA_DB.Create(&bottle).Error
}

// PickRandomBottle 捞一个随机瓶子（不捞自己的，只捞漂流中的）
func (s *DriftBottleService) PickRandomBottle(deviceID string) (driftbottle.DriftBottle, error) {
	var bottle driftbottle.DriftBottle
	err := global.GVA_DB.
		Where("status = ? AND device_id != ?", 1, deviceID).
		Order("RANDOM()").
		First(&bottle).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return bottle, errors.New("暂时没有可捞的瓶子")
	}
	return bottle, err
}

// ReplyBottle 回复瓶子
func (s *DriftBottleService) ReplyBottle(req driftbottleReq.ReplyBottleRequest) error {
	var bottle driftbottle.DriftBottle
	if err := global.GVA_DB.First(&bottle, req.BottleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("瓶子不存在")
		}
		return err
	}
	reply := driftbottle.DriftBottleReply{
		BottleID: req.BottleID,
		Nickname: req.Nickname,
		Content:  req.Content,
		DeviceID: req.DeviceID,
	}
	if err := global.GVA_DB.Create(&reply).Error; err != nil {
		return err
	}
	// 更新瓶子状态为已回复
	return global.GVA_DB.Model(&bottle).Update("status", 3).Error
}

// GetMyBottles 获取我的瓶子列表
func (s *DriftBottleService) GetMyBottles(deviceID string, info request.PageInfo) (list []driftbottle.DriftBottle, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&driftbottle.DriftBottle{}).Where("device_id = ?", deviceID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return
}

// GetBottleReplies 获取瓶子的回复列表
func (s *DriftBottleService) GetBottleReplies(bottleID uint) ([]driftbottle.DriftBottleReply, error) {
	var replies []driftbottle.DriftBottleReply
	err := global.GVA_DB.Where("bottle_id = ?", bottleID).Order("created_at ASC").Find(&replies).Error
	return replies, err
}

// GetBottleDetail 获取瓶子详情含回复
func (s *DriftBottleService) GetBottleDetail(bottleID uint) (driftbottle.DriftBottle, []driftbottle.DriftBottleReply, error) {
	var bottle driftbottle.DriftBottle
	if err := global.GVA_DB.First(&bottle, bottleID).Error; err != nil {
		return bottle, nil, err
	}
	replies, err := s.GetBottleReplies(bottleID)
	return bottle, replies, err
}

// AdminGetBottleList 后台获取瓶子列表
func (s *DriftBottleService) AdminGetBottleList(req driftbottleReq.BottleListRequest) (list []driftbottle.DriftBottle, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.GVA_DB.Model(&driftbottle.DriftBottle{})
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		db = db.Where("content LIKE ? OR nickname LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return
}

// AdminGetReplyList 后台获取回复列表
func (s *DriftBottleService) AdminGetReplyList(req driftbottleReq.ReplyListRequest) (list []driftbottle.DriftBottleReply, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.GVA_DB.Model(&driftbottle.DriftBottleReply{})
	if req.BottleID > 0 {
		db = db.Where("bottle_id = ?", req.BottleID)
	}
	if req.Keyword != "" {
		db = db.Where("content LIKE ? OR nickname LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return
}

// AdminDeleteBottle 后台删除瓶子（软删除）
func (s *DriftBottleService) AdminDeleteBottle(id uint) error {
	return global.GVA_DB.Delete(&driftbottle.DriftBottle{}, id).Error
}

// AdminDeleteReply 后台删除回复（软删除）
func (s *DriftBottleService) AdminDeleteReply(id uint) error {
	return global.GVA_DB.Delete(&driftbottle.DriftBottleReply{}, id).Error
}
