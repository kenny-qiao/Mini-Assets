package service

import (
	"fmt"
	"mini-assets/internal/model"
	"mini-assets/internal/repository"

	"github.com/google/uuid"
)

// AssetService 定义了资产相关的业务逻辑接口
type AssetService interface {
	CreateAsset(asset *model.Asset) error
	GetAssetByID(id uuid.UUID) (*model.Asset, error)
	UpdateAsset(asset *model.Asset) error
	DeleteAsset(id uuid.UUID) error
	ListAssets() ([]*model.Asset, error)
	DiscardAsset(assetID uuid.UUID, userID uuid.UUID) error
}

// assetService 是AssetService接口的实现
type assetService struct {
	assetRepo repository.AssetRepository
}

// NewAssetService 创建并返回一个AssetService实例
func NewAssetService(assetRepo repository.AssetRepository) AssetService {
	return &assetService{assetRepo: assetRepo}
}

// CreateAsset 创建一个新的资产记录
func (s *assetService) CreateAsset(asset *model.Asset) error {
	return s.assetRepo.Create(asset)
}

// GetAssetByID 根据资产ID获取资产信息
func (s *assetService) GetAssetByID(id uuid.UUID) (*model.Asset, error) {
	return s.assetRepo.GetByID(id)
}

// UpdateAsset 更新资产信息
func (s *assetService) UpdateAsset(asset *model.Asset) error {
	return s.assetRepo.Update(asset)
}

// DeleteAsset 删除资产记录
func (s *assetService) DeleteAsset(id uuid.UUID) error {
	return s.assetRepo.Delete(id)
}

// ListAssets 获取所有资产记录
func (s *assetService) ListAssets() ([]*model.Asset, error) {
	return s.assetRepo.List()
}

// DiscardAsset 根据资产ID和用户ID删除资产
func (s *assetService) DiscardAsset(assetID uuid.UUID, userID uuid.UUID) error {
	// 检查资产是否属于该用户
	asset, err := s.assetRepo.GetByID(assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %w", err)
	}

	if asset.CreatorID != userID {
		return fmt.Errorf("asset does not belong to user")
	}

	// 删除资产
	err = s.assetRepo.Delete(assetID)
	if err != nil {
		return fmt.Errorf("failed to discard asset: %w", err)
	}

	return nil
}
