package httpgin

import (
	_ "csTrade/docs"
	"csTrade/internal/handlers/middleware"
	"csTrade/internal/repository"
	"csTrade/internal/service"
	"csTrade/internal/service/bots"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(repo *repository.Repository, botmanager *bots.BotManager) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	offerServ := service.NewOfferService(repo, botmanager)
	offerHandler := NewOfferHandler(offerServ)

	userServ := service.NewUserService(repo)
	userHandler := NewUserHandler(userServ)

	transactionServ := service.NewTransactionService(repo)
	transactionHandler := NewTransactionHandler(transactionServ)

	{
		r.GET("/swagger", ginSwagger.WrapHandler(swaggerfiles.Handler))
		r.GET("/healthz", func(c *gin.Context) {
			c.String(200, "ok")
		})
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	api := r.Group("/api/v1")

	api.POST("/users/create", userHandler.CreateUser)

	users := api.Group("/users").Use(middleware.AuthMiddleware())
	{
		users.GET("/:id")
		users.GET("/:id/cash")
		users.PATCH("/:id/cash")
	}

	listings := api.Group("/market/listings")
	{
		listings.GET("", offerHandler.GetAllOffers)
		listings.POST("", offerHandler.ListSkin)              // sell
		listings.POST("/:id/purchase", offerHandler.Purchase) // buy
		listings.GET("/:id", offerHandler.GetOfferByID)
		listings.GET("/user/:id", offerHandler.UserOffers)
		listings.POST("/cancel", offerHandler.CancelTrade)
		listings.PATCH("/:id/price", offerHandler.ChangePrice)
		listings.DELETE("/:id", offerHandler.DeleteByID)
	}

	transaction := api.Group("/transaction").Use(middleware.AuthMiddleware())
	{
		transaction.GET("/:id")
		transaction.PATCH("/:id/status")
		transaction.GET("/user/:id", transactionHandler.GetByuerTransaction)
	}

	return r
}
