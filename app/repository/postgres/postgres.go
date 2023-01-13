package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	ConnMaxLifeTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}

const (
	userTable       = "users"
	forumTable      = "forums"
	threadTable     = "threads"
	forumUsersTable = "nickname_forum"
	postTable       = "posts"
	voteTable       = "votes"
)

func NewPostgresDB(conf Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		conf.Host, conf.Port, conf.Username, conf.DBName, conf.Password, conf.SSLMode))

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(conf.ConnMaxLifeTime)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
