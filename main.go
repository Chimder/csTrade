package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	LoggerInit()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Bit()
	Parse()
	// r := gin.Default()
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"https://*", "http://*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposeHeaders:    []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300,
	// }))

	// r.GET("/healthz", func(c *gin.Context) {
	// 	c.String(200, "ok")
	// })

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
	slog.Info("Server is running...")

	<-ctx.Done()
	slog.Info("Shutting down server...")

	// shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// if err := srv.Shutdown(shutdownCtx); err != nil {
	// 	slog.Error("Server shutdown error", "error", err)
	// } else {
	// 	slog.Info("Server stopped gracefully")
	// }
}
