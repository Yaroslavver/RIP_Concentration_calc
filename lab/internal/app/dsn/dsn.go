package dsn

import (
	"fmt"
	"os"
)

// FromEnv возвращает строку DSN для подключения к PostgreSQL из переменных окружения
func FromEnv() string {
	// Читаем переменные окружения
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// Проверяем, что все необходимые переменные заданы
	if host == "" || port == "" || user == "" || pass == "" || dbname == "" {
		// Можно вернуть пустую строку или значения по умолчанию
		// Для отладки можно вывести предупреждение
		fmt.Println("Warning: not all database environment variables are set")
	}

	// Формируем DSN строку для PostgreSQL
	// sslmode=disable - отключаем SSL для локальной разработки
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)
}

// FromEnvWithSSL - альтернативная версия с поддержкой SSL
func FromEnvWithSSL() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	
	if sslmode == "" {
		sslmode = "disable" // по умолчанию отключаем SSL
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode)
}