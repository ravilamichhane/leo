package db

import (
	"fmt"

	"github.com/ravilmc/leo/web"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetPostgresClient() *gorm.DB {
	host := web.GetEnv("DB_HOST", "localhost")
	port := web.GetEnv("DB_PORT", "5432")
	user := web.GetEnv("DB_USER", "postgres")
	password := web.GetEnv("DB_PASSWORD", "postgres")
	dbname := web.GetEnv("DB_NAME", "users")

	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	client, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	return client
}
