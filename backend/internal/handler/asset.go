package handler

import (
	"mini-assets/internal/model"
	"mini-assets/internal/service"
	"mini-assets/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssetHandler 处理资产相关的HTTP请求
type AssetHandler struct {
	AssetService service.AssetService
}

// NewAssetHandler 创建AssetHandler实例
func NewAssetHandler(assetService service.AssetService) *AssetHandler {
	return &AssetHandler{AssetService: assetService}
}

// CreateAsset 创建资产
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var asset model.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 假设从上下文中获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	asset.CreatorID = userID.(uuid.UUID)
	asset.ID = uuid.New()

	if err := h.AssetService.CreateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

// UpdateAsset 更新资产
func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	var asset model.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户是否为资产的创建者
	userID, exists := c.Get("userID")
	if !exists || asset.CreatorID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to update this asset"})
		return
	}

	if err := h.AssetService.UpdateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// DiscardAsset 废弃资产
func (h *AssetHandler) DiscardAsset(c *gin.Context) {
	assetID, exists := c.Get("assetID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "asset not found"})
		return
	}

	// 验证用户是否为资产的创建者
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if err := h.AssetService.DiscardAsset(assetID.(uuid.UUID), userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "asset discarded"})
}

// SetupAssetRoutes 设置资产相关的路由
func SetupAssetRoutes(r *gin.Engine, assetService service.AssetService) {
	h := NewAssetHandler(assetService)
	r.POST("/assets", util.AuthMiddleware(), h.CreateAsset)
	r.PUT("/assets/:id", util.AuthMiddleware(), h.UpdateAsset)
	r.DELETE("/assets/:id", util.AuthMiddleware(), h.DiscardAsset)
}
