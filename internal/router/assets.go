package router

import (
	"net/http"
	"strconv"

	"github.com/B022MC/soraka-backend/internal/service/assets"
	"github.com/gin-gonic/gin"
)

type AssetsRouter struct {
	assetsService *assets.AssetsService
}

func NewAssetsRouter() *AssetsRouter {
	return &AssetsRouter{
		assetsService: assets.NewAssetsService(),
	}
}

func (r *AssetsRouter) InitRouter(group *gin.RouterGroup) {
	assetsGroup := group.Group("/assets")
	{
		assetsGroup.GET("/champion/:id", r.getChampionIcon)
		assetsGroup.GET("/profile/:id", r.getProfileIcon)
		assetsGroup.GET("/item/:id", r.getItemIcon)
		assetsGroup.GET("/spell/:name", r.getSpellIcon)
		assetsGroup.GET("/perk/:style", r.getPerkIcon)
	}
}

func (r *AssetsRouter) getChampionIcon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, contentType, err := r.assetsService.GetChampionIcon(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000")
	c.Data(http.StatusOK, contentType, data)
}

func (r *AssetsRouter) getProfileIcon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, contentType, err := r.assetsService.GetProfileIcon(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000")
	c.Data(http.StatusOK, contentType, data)
}

func (r *AssetsRouter) getItemIcon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, contentType, err := r.assetsService.GetItemIcon(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000")
	c.Data(http.StatusOK, contentType, data)
}

func (r *AssetsRouter) getSpellIcon(c *gin.Context) {
	name := c.Param("name")

	data, contentType, err := r.assetsService.GetSpellIcon(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000")
	c.Data(http.StatusOK, contentType, data)
}

func (r *AssetsRouter) getPerkIcon(c *gin.Context) {
	style := c.Param("style")

	data, contentType, err := r.assetsService.GetPerkIcon(style)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000")
	c.Data(http.StatusOK, contentType, data)
}
