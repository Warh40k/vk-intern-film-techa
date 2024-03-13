package main

import (
	"context"
	"errors"
	httpserver "github.com/Warh40k/vk-intern-filmotecka/internal/api/handler"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository/postgres"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/service"
	"github.com/Warh40k/vk-intern-filmotecka/internal/app"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("local")
	return viper.ReadInConfig()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func main() {
	log := setupLogger(viper.GetString("env"))

	if err := initConfig(); err != nil {
		log.Error("Ошибка чтения конфигурации: %s", err.Error())
		panic(err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Error("Ошибка чтения переменных окружения: %s", err.Error())
		panic(err.Error())
	}

	pgCfg := postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	db, err := postgres.NewPostgresDB(pgCfg)
	if err != nil {
		log.Error("Ошибка подключения к базе данных: %s", err.Error())
		return
	}

	if err != nil {
		log.Error("Ошибка подключения к кэшу: %s", err.Error())
		return
	}

	repos := repository.NewRepository(db, log)
	services := service.NewService(repos, log)
	handlers := httpserver.NewHandler(services, log)
	serv := new(app.App)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		if err = serv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Ошибка запуска http сервера: %s", err.Error())
			panic(err.Error())
		}
	}()
	log.Info("server started")
	<-quit

	log.Info("trying to gracefull shutdown")
	if err = serv.Shutdown(context.Background()); err != nil {
		log.With(slog.String("err", err.Error())).Error("error occured on server shutting down:")
	}

	if err = db.Close(); err != nil {
		log.Error("error occured on db connection close: %s", err.Error())
		return
	}

	log.Info("gracefully stopped")
}