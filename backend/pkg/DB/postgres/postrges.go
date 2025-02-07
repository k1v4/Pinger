package postgres

//
//import (
//	"context"
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	_ "github.com/lib/pq"
//	"log"
//)
//
//type DBConfig struct {
//	UserName string `env:"POSTGRES_USER" env-default:"root"`
//	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
//	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
//	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
//	DbName   string `env:"POSTGRES_DB" env-default:"user_service"`
//}
//
//type DB struct {
//	Db *sqlx.DB
//}
//
//func New(cfg DBConfig) (*DB, error) {
//	connectString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
//		cfg.UserName, cfg.Password, cfg.DbName, cfg.Host, cfg.Port)
//
//	db, err := sqlx.Connect("postgres", connectString)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	if _, err = db.Conn(context.Background()); err != nil {
//		return nil, fmt.Errorf("failed to connect to database: %w", err)
//	}
//
//	return &DB{Db: db}, nil
//}
