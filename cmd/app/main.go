package main

import (
	"context"
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"protravel-finance/internal/config"
	"protravel-finance/internal/cron"
	"protravel-finance/internal/modules/currency"
	"protravel-finance/internal/modules/user"
	"protravel-finance/internal/runner"
	http_server "protravel-finance/internal/server/http"
	"protravel-finance/internal/shared/database/postgres"
	redis2 "protravel-finance/internal/shared/database/redis"
	"protravel-finance/internal/shared/middleware"
	"protravel-finance/internal/shared/response"
	transperr "protravel-finance/internal/shared/transport_error"
	"protravel-finance/pkg/clients/exchangerate"
	"protravel-finance/pkg/logger"
	"syscall"
)

func main() {
	var ctx = context.Background()

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := config.MustConfig(log)

	pool, err := postgres.New(ctx, cfg.Postgres, log)
	if err != nil {
		log.Panic(err)
	}
	defer pool.Close()

	transaction := postgres.NewTransactionsRepos(cfg.Postgres, pool)

	client, err := redis2.NewRedis(ctx, cfg.Redis, log)
	if err != nil {
		log.Panic("failed to connected redis: ", err)
	}
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			log.Error("failed to close redis client: ", err)
		}
	}(client)

	var (
		mid              = middleware.NewMiddleware(log)
		httpResp         = response.NewHTTPResponse(log, true)
		convert          = transperr.NewErrorConverter()
		validationFormat = strfmt.NewFormats()

		router = mux.NewRouter()
	)

	initBusinessLogic(
		router,

		mid,
		httpResp,
		convert,
		validationFormat,

		transaction,
		client,

		log,
		*cfg,
	)
	httpServer := http_server.NewServer(&cfg.Server, router)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server failed: %v", err)
		}
	}()

	log.Infof("server listening on port [%d] | Env %s", cfg.Server.Port, cfg.Env)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Errorf("error shutdown: %s", err)
	}

	log.Info("server shutdown")
}

func initBusinessLogic(
	router *mux.Router,

	mid middleware.Middleware,
	httpResp response.HttpResponse,
	convert transperr.ErrorConverter,
	validationFormat strfmt.Registry,

	transaction postgres.Transaction,
	redisClient *redis.Client,

	log logger.Logger,
	cfg config.Config,
) {
	// Init Repositories
	userRepos := user.NewRepository()
	currencyRepos := currency.NewRepository()

	// Init Services
	userService := user.NewService(userRepos, transaction)
	currencyService := currency.NewService(currencyRepos, transaction, redisClient)

	// Init client api
	exchangeRateApi := exchangerate.NewExchangeRate(cfg.Exchangerate.BaseURL, log)

	runner.InitHandlers(router, mid,
		user.NewRunnerHandlerV1(router, httpResp, convert, log, validationFormat, userService),
	)

	runner.InitCronTasks(log,
		cron.NewUpdateRatesCron(currencyService, exchangeRateApi, log),
	)
}
