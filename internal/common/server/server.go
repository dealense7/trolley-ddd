package server

import (
	"net/http"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	//r.Use(middleware.LanguageMiddleware)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))

	return r
}

func Start(cfg *cfg.Config, logFactory logger.Factory, r *chi.Mux) {
	go func() {
		log := logFactory.For(logger.General)

		if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
			log.Error("Could not start server", zap.Error(err))
			panic(err)
		}

		log.Info("Server started", zap.String("time", time.Now().Format(time.RFC3339)))
	}()
}
