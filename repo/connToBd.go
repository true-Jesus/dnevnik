package repo

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectToDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=admin password=1233 dbname=bobr sslmode=disable" // Замените на ваши настройки подключения
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки соединения: %w", err)
	}

	_, err = db.Exec(`
		 CREATE TABLE IF NOT EXISTS users (
      username VARCHAR(255) PRIMARY KEY,
      password_hash VARCHAR(255) NOT NULL
    )
  `)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
