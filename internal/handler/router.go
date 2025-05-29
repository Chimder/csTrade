package handler

import (
	_ "csTrade/docs"
	"csTrade/internal/handler/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() *gin.Engine {
	// cfg := config.LoadEnv()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// client := bots.NewSteamClient(
	// 	cfg.Username,
	// 	cfg.Password,
	// 	cfg.SteamID,
	// 	cfg.SharedSecret,
	// 	cfg.IdentitySecret,
	// 	"",
	// )
	// err := client.Login()
	// if err != nil {
	// 	fmt.Printf("Login failed: %v\n", err)
	// }
	// fmt.Println("Successfully logged in to Steam!")
	// fmt.Printf("Access Token: %s\n", client.AccessToken)

	{
		r.GET("/swagger", ginSwagger.WrapHandler(swaggerfiles.Handler))
		r.GET("/healthz", func(c *gin.Context) {
			c.String(200, "ok")
		})
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.GET("/steam", func(ctx *gin.Context) {})
		auth.GET("/steam/callback", func(ctx *gin.Context) {})
	}

	user := api.Group("/user").Use(middleware.AuthMiddleware())
	{
		user.GET("/inventory", func(ctx *gin.Context) {})
		user.POST("/inventory/sync", func(ctx *gin.Context) {})
	}

	market := api.Group("/market")
	{
		market.GET("/items", func(ctx *gin.Context) {})
		market.GET("/item/:id", func(ctx *gin.Context) {})

		marketPrivate := market.Use(middleware.AuthMiddleware())
		{
			marketPrivate.POST("/item/sell", func(ctx *gin.Context) {})
			marketPrivate.PUT("/item/:id", func(ctx *gin.Context) {}) //change price
			marketPrivate.DELETE("/item/:id", func(ctx *gin.Context) {})
			marketPrivate.POST("/listings/:id/buy", func(ctx *gin.Context) {}) //buy
		}
	}

	return r
}
