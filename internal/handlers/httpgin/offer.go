package httpgin

import (
	// offer "csTrade/internal/app"

	"csTrade/internal/domain/offer"
	"csTrade/internal/service"

	"github.com/gin-gonic/gin"
)

type OfferHandler struct {
	service *service.OfferService
}

func NewOfferHandler(service *service.OfferService) *OfferHandler {
	return &OfferHandler{service: service}
}

func (ofh *OfferHandler) CreateOffer(c *gin.Context) {

	var req *offer.OfferCreateReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := ofh.service.CreateOffer(c.Request.Context(), req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}
