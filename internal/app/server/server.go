package server

import (
	"context"
	_ "cryptoObserver/docs"
	"cryptoObserver/internal/app/handlers"
	"cryptoObserver/internal/app/store/sqlstore"
	worker "cryptoObserver/internal/app/workers"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type App struct {
	router *chi.Mux
	Server *http.Server
	store  sqlstore.StoreInterface
	logger *logrus.Logger
	ctx    context.Context
	config Config
	pool   *worker.WorkerPool
}

func newApp(ctx context.Context, store sqlstore.StoreInterface, config Config, logger *logrus.Logger, pool *worker.WorkerPool) *App {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", config.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,  // Защита от Slowloris
		ReadTimeout:       10 * time.Second, // Максимальное время чтения всего запроса
		WriteTimeout:      30 * time.Second, // Максимальное время записи ответа
		IdleTimeout:       60 * time.Second, // Таймаут для keep-alive соединений
	}

	logger.Out = os.Stdout
	log.SetOutput(os.Stdout)
	a := &App{
		router: router,
		Server: server,
		store:  store,
		logger: logger,
		ctx:    ctx,
		config: config,
		pool:   pool,
	}
	a.configureRouter()
	return a
}

func (a *App) configureRouter() {
	a.router.Use(middleware.Recoverer)
	a.router.Use(a.logRequest)
	a.router.Route("/currency", func(r chi.Router) {
		r.Post("/add", handlers.NewAddCurrencyHandler(a.logger, a.store.Currency(), a.pool))
		r.Delete("/remove", handlers.NewRemoveCurrencyHandler(a.logger, a.store.Currency(), a.pool))
		r.Post("/price", handlers.NewGetPriceHandler(a.logger, a.store.Currency()))
	})
	a.router.Get("/api/doc/*", httpSwagger.WrapHandler)
}

func (a *App) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := a.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"method":      r.Method,
			"real_ip":     getClientIP(r),
			"build_type":  r.Header.Get("X-App-Build-Type"),
			"version":     r.Header.Get("X-App-Version"),
		})

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		if strings.HasPrefix(r.Method, "/metrics") {
			return
		}
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func getClientIP(r *http.Request) string {
	// Проверяем X-Forwarded-For
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ips := strings.Split(fwd, ",")
		// Возвращаем первый IP в цепочке
		return strings.TrimSpace(ips[0])
	}

	// Проверяем X-Real-IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Если ничего не найдено, используем RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *App) Close() error {
	err := a.Server.Close()
	if err != nil {
		return err
	}
	return a.Server.Close()
}
