package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"github.com/o3labs/openpoint/platform/config"
)

var (
	DB_URL = os.Getenv("DB_URL")
)
var db *sql.DB
var once sync.Once

func Connect() *sql.DB {
	once.Do(func() {
		connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			config.Env.Database.Host,
			config.Env.Database.User,
			config.Env.Database.Name,
			config.Env.Database.Password)

		var err error
		if DB_URL == "" {
			DB_URL = connectionString
		}
		db, err = sql.Open("postgres", DB_URL)
		if err != nil {
			log.Printf("cannot open database %v", err)
			db = nil
		}
		db.SetMaxOpenConns(256)
		// db.SetConnMaxLifetime(time.Hour * 24)
		db.SetMaxIdleConns(256)

	})
	return db
}
