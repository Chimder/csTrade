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

	offerServ := service.NewOfferService(repo,botmanager)
	offerHandler := NewOfferHandler(offerServ)
	userServ := service.NewUserService(repo)
	userHandler := NewUserHandler(userServ)

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
		users.GET("/:id/offers")
	}

	offers := api.Group("/offers")
	{
		offers.POST("create", offerHandler.CreateOffer)
		offers.GET("/:id")
		offers.GET("")
		offers.PATCH("/:id/price")
		offers.DELETE("/:id")
	}

	transaction := api.Group("/transaction").Use(middleware.AuthMiddleware())
	{
		transaction.POST("", func(ctx *gin.Context) {})
		transaction.GET("/:id", func(ctx *gin.Context) {})
		transaction.PATCH("/:id/status", func(ctx *gin.Context) {})
		transaction.GET("/buyer/:id", func(ctx *gin.Context) {})
		transaction.GET("/seller/:id", func(ctx *gin.Context) {})
	}

	return r
}
