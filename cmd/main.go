package main

import (
	"context"
	"csTrade/db"
	router "csTrade/internal/api/http"
	"csTrade/internal/app/bots"
	"csTrade/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

//		@title			csTrairde Api
//		@version		1.0
//		@description	CSGO trade
//	  @BasePath	/
func main() {
	SetupLogger()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbconn, err := db.DBConn(ctx)
	if err != nil {
		log.Panic().Msg("Err conn to db")
		return
	}
	repo := repository.NewRepository(dbconn)

	///////////////////
	botmanager := bots.NewBotManager(repo)
	go botmanager.Start(ctx)
	//////////////////////

	r := router.Init()
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()

	log.Info().Msg("Server is running...")
	<-ctx.Done()
	log.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	} else {
		log.Info().Msg("Server stopped gracefully")
	}
}
