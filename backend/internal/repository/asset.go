package repository

import (
	"database/sql"
	"fmt"
	"mini-assets/internal/model"

	"github.com/google/uuid"
)

// AssetRepository 定义了资产相关的数据库操作接口
type AssetRepository interface {
	Create(asset *model.Asset) error
	GetByID(id uuid.UUID) (*model.Asset, error)
	Update(asset *model.Asset) error
	Delete(id uuid.UUID) error
	List() ([]*model.Asset, error)
}

// assetRepository 是AssetRepository接口的实现
type assetRepository struct {
	db *sql.DB
}

// NewAssetRepository 创建并返回一个AssetRepository实例
func NewAssetRepository(db *sql.DB) AssetRepository {
	return &assetRepository{db: db}
}

// Create 插入一个新的资产记录到数据库
func (r *assetRepository) Create(asset *model.Asset) error {
	query := `INSERT INTO assets (id, category, name, currency, amount, creator_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, uuid.New().String(), asset.Category, asset.Name, asset.Currency, asset.Amount, asset.CreatorID, asset.CreatedAt, asset.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}
	return nil
}

// GetByID 根据资产ID从数据库中获取资产信息
func (r *assetRepository) GetByID(id uuid.UUID) (*model.Asset, error) {
	query := `SELECT id, category, name, currency, amount, creator_id, created_at, updated_at FROM assets WHERE id = ?`
	row := r.db.QueryRow(query, id.String())

	var asset model.Asset
	if err := row.Scan(&asset.ID, &asset.Category, &asset.Name, &asset.Currency, &asset.Amount, &asset.CreatorID, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, fmt.Errorf("failed to get asset by ID: %w", err)
	}
	return &asset, nil
}

// Update 更新数据库中的资产信息
func (r *assetRepository) Update(asset *model.Asset) error {
	query := `UPDATE assets SET category = ?, name = ?, currency = ?, amount = ?, creator_id = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, asset.Category, asset.Name, asset.Currency, asset.Amount, asset.CreatorID, asset.UpdatedAt, asset.ID.String())
	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}
	return nil
}

// Delete 根据资产ID删除数据库中的资产记录
func (r *assetRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM assets WHERE id = ?`
	_, err := r.db.Exec(query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	return nil
}

// List 获取数据库中的所有资产记录
func (r *assetRepository) List() ([]*model.Asset, error) {
	query := `SELECT id, category, name, currency, amount, creator_id, created_at, updated_at FROM assets`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		var asset model.Asset
		if err := rows.Scan(&asset.ID, &asset.Category, &asset.Name, &asset.Currency, &asset.Amount, &asset.CreatorID, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over assets: %w", err)
	}

	return assets, nil
}
