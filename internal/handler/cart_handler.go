package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexRijikov/go-petshop-api/internal/service"
	"github.com/gin-gonic/gin"
)

// CartHandler обробляє HTTP-запити, пов'язані з кошиком користувача

type CartHandler struct {
	svc services.CartService
}

// NewCartHandler створює новий CartHandler з наданим сервісом

func NewCartHandler(s services.CartService) *CartHandler {
	return &CartHandler{svc: s}
}

// RegisterRoutes реєструє маршрути кошика у вказаній групі маршрутизатора (rg *gin.RouterGroup)

func (h *CartHandler) RegisterRoutes(rg *gin.RouterGroup) {
	grp := rg.Group("/cart")        
	grp.GET("", h.List)             
	grp.POST("/add", h.Add)         
	grp.PUT("/update", h.Update)    
	grp.DELETE("/remove", h.Remove) 

}

// List отримує всі товари в кошику користувача

func (h *CartHandler) List(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"}) //
		return
	}
	items, err := h.svc.List(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// modifyCartRequest - ( ) 

type modifyCartRequest struct {
	UserID    uint `json:"user_id" binding:"required"`  
	ProductID uint `json:"product_id" binding:"required"` 
	Quantity  int  `json:"quantity" binding:"required,gt=0"` 
}

// Add - 

func (h *CartHandler) Add(c *gin.Context) {
	var req modifyCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddItem(c.Request.Context(), req.UserID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "item added to cart"})
}

// Update -

func (h *CartHandler) Update(c *gin.Context) {
	var req modifyCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateQuantity(c.Request.Context(), req.UserID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "quantity updated"})
}

// Remove -

func (h *CartHandler) Remove(c *gin.Context) {
	userIDStr := c.Query("user_id")
	productIDStr := c.Query("product_id")

	userID, err1 := strconv.ParseUint(userIDStr, 10, 64)
	productID, err2 := strconv.ParseUint(productIDStr, 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id or product_id"})
		return
	}
	if err := h.svc.RemoveItem(c.Request.Context(), uint(userID), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
}
