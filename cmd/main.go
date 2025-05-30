package main

import (
	"csTrade/config"
	"csTrade/internal/domain/bots"
	"fmt"
	"time"
	// _ "github.com/lib/pq"
)

//		@title			csTrairde Api
//		@version		1.0
//		@description	CSGO trade
//	  @BasePath	/
func main() {
	// LoggerInit()

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	// r := handler.Init()
	cfg := config.LoadEnv()
	client := bots.NewSteamClient(
		cfg.Username,
		cfg.Password,
		cfg.SteamID,
		cfg.SharedSecret,
		cfg.IdentitySecret,
		"",
	)
	err := client.Login()
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
	}
	fmt.Println("Successfully logged in to Steam!")
	fmt.Printf("Access Token: %s\n", client.AccessToken)
	time.Sleep(5 * time.Second)

	// srv := &http.Server{
	// 	Addr:         ":8080",
	// 	Handler:      r,
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	IdleTimeout:  120 * time.Second,
	// }

	// go func() {
	// 	if err := srv.ListenAndServe(); err != nil {
	// 		log.Fatalf("Server error: %v", err)
	// 	}
	// }()

	// slog.Info("Server is running...")
	// <-ctx.Done()
	// slog.Info("Shutting down server...")

	// shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// if err := srv.Shutdown(shutdownCtx); err != nil {
	// 	slog.Error("Server shutdown error", "error", err)
	// } else {
	// 	slog.Info("Server stopped gracefully")
	// }
}
