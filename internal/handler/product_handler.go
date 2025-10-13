package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/service"
	"github.com/gin-gonic/gin"
)

// Відповідає за обробку HTTP-запитів

// Використовуємо binding для базової валідації.

// List підтримує limit і offset (пагінація).

// Помилки переводимо в HTTP-статуси.

// ProductHandler обробляє HTTP-запити, пов'язані з продуктами
type ProductHandler struct {
	svc services.ProductService
}

// NewProductHandler створює новий ProductHandler з наданим сервісом
func NewProductHandler(s services.ProductService) *ProductHandler {
	return &ProductHandler{svc: s}
}

// RegisterRoutes реєструє маршрути продуктів у вказаній групі маршрутизатора (rg *gin.RouterGroup)
func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
	grp := rg.Group("/products")
	grp.GET("", h.List)
	grp.POST("", h.Create)
	grp.GET("/:id", h.GetByID)
	grp.PUT("/:id", h.Update)
	grp.DELETE("/:id", h.Delete)
}

// createProductRequest використовується для прив'язки та валідації вхідних даних при створенні або оновленні продукту

type createProductRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	PriceCents  int64  `json:"price_cents" binding:"required,gt=0"`
	Stock       int    `json:"stock" binding:"gte=0"`
	SKU         string `json:"sku" binding:"omitempty,max=100"`
}

// Create (Створення нового продукту)

func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}
	created, err := h.svc.CreateProduct(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// List (Список продуктів з пагінацією)

func (h *ProductHandler) List(c *gin.Context) {
	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	items, total, err := h.svc.ListProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "limit": limit, "offset": offset})
}

// GetByID (Отримання продукту за ID)

func (h *ProductHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Update (Оновлення продукту)

func (h *ProductHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &models.Product{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}
	updated, err := h.svc.UpdateProduct(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete (Видалення продукту)

func (h *ProductHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
