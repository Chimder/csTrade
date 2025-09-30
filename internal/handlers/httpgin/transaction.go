package httpgin

import (
	"csTrade/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (trh *TransactionHandler) GetByuerTransaction(c *gin.Context) {
	id := c.Param("id")

	data, err := trh.service.GetTransactionByBuyerID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offer not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// func (trh *TransactionHandler) GetSellerTransaction(c *gin.Context) {}
