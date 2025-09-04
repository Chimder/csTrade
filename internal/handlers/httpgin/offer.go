package httpgin

import (
	// offer "csTrade/internal/app"

	"csTrade/internal/domain/offer"
	"csTrade/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OfferHandler struct {
	service *service.OfferService
}

func NewOfferHandler(service *service.OfferService) *OfferHandler {
	return &OfferHandler{service: service}
}

func (ofh *OfferHandler) ListSkin(c *gin.Context) {

	var req offer.OfferCreateReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := ofh.service.ReceiveFromUserOffer(c.Request.Context(), &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (ofh *OfferHandler) Purchase(c *gin.Context) {

	var req offer.OfferCreateReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err := ofh.service.SendToBuyerOffer(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (ofh *OfferHandler) GetOfferByID(c *gin.Context) {
	offerID := c.Param("id")
	data, err := ofh.service.GetByID(c.Request.Context(), offerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offer not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (ofh *OfferHandler) GetAllOffers(c *gin.Context) {
	data, err := ofh.service.GetAllOffers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
func (ofh *OfferHandler) GetTradeStatus(c *gin.Context) {
	steamOfferID := c.Query("steam_id")

	log.Info().Msg("start")
	status, err := ofh.service.GetTradeStatus(c.Request.Context(), steamOfferID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": status})
}

func (ofh *OfferHandler) CancelTrade(c *gin.Context) {
	steamOfferID := c.Query("steam_id")

	log.Info().Msg("start")
	err := ofh.service.CancelTrade(c.Request.Context(), steamOfferID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (ofh *OfferHandler) ChangePrice(c *gin.Context) {
	id := c.Param("id")
	priceStr := c.Param("price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	err = ofh.service.ChangePriceByID(c.Request.Context(), id, price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (ofh *OfferHandler) UserOffers(c *gin.Context) {
	id := c.Param("id")

	data, err := ofh.service.GetUserOffers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}

func (ofh *OfferHandler) DeleteByID(c *gin.Context) {
	offerId := c.Param("id")

	err := ofh.service.ChangeStatusByID(c.Request.Context(), offer.OfferCanceled, offerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "offer deleted"})
}

// func (ofh *OfferHandler) DeleteByID(c *gin.Context) {
// 	id := c.Param("id")

// 	err := ofh.service.DeleteByID(c.Request.Context(), id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(200, gin.H{"message": "ok"})
// }
