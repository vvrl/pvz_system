package app

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Server struct {
		Port string `envconfig:"PORT"`
	}

	Database struct {
		Dsn    string `envconfig:"DB_DSN"`
		Driver string `envconfig:"DB_DATABASE" default:"postgres"`
	}

	Migrations struct {
		Dir string `envconfig:"MIGRATIONS_DIR"`
	}
}

type App struct {
	config Config
}

func NewApp() (*App, error) {

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования в структуру config %w", err)
	}

	return &App{cfg}, nil

}

func (a *App) Run() {

	config, err := pgx.ParseConfig(a.config.Database.Dsn)
	if err != nil {
		log.Fatalf("Ошибка парсинга pgx конфига: %v", err)
	}

	db := stdlib.OpenDB(*config)

	if err := goose.Up(db, a.config.Migrations.Dir); err != nil {
		log.Fatalf("Ошибка миграций: %v", err)
	}

}
