package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	UserName string `env:"POSTGRES_USER" env-default:"root"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
	Host     string `env:"POSTGRES_HOST" env-default:"postgres"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string `env:"POSTGRES_DB" env-default:"order_service"`
}

type DB struct {
	Db *sqlx.DB
}

func New(config Config) (*DB, error) {
	connectString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.UserName, config.Password, config.DbName, config.Host, config.Port)

	fmt.Println(connectString)

	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err = db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{Db: db}, nil
}
