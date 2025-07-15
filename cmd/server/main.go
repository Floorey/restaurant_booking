package main

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"restaurant_booking/internal/handler"
	"restaurant_booking/internal/job"
	"time"
)

func main() {
	// Logger
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// DB connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal().Msg("DATABASE_URL not set")
	}
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("DB connection failed")
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	// start backround-jobs
	job.StartBackroundJobs(db)

	// Gin-router + sessions

	r := gin.Default()
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("sess", store))

	r.SetHTMLTemplate(template.Must(
		template.ParseGlob("web/templates/*.tmpl")))

	handler.RegisterWeb(r, db)
	handler.RegisterAdmin(r.Group("/admin"), db)

	srv := &http.Server{
		Addr:        ":8080",
		Handler:     r,
		ReadTimeout: 5 * time.Second,
	}

	// start async
	go func() {
		log.Info().Msg("⇢ server running at: http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("ListenAndServe")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server Shutdown")
	}
}
