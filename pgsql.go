package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/ini.v1"
)

type SQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func findConfigFile() (string, error) {
	// Проверяем текущий каталог
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	currentDir := filepath.Dir(exePath)
	configPath := filepath.Join(currentDir, "srv_cfg.ini")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	// Проверяем каталог уровнем выше
	parentDir := filepath.Dir(currentDir)
	configPath = filepath.Join(parentDir, "srv_cfg.ini")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	return "", fmt.Errorf("srv_cfg.ini not found in current or parent directory")
}

func loadSQLConfig() (*SQLConfig, error) {
	// Находим файл конфигурации
	configPath, err := findConfigFile()
	if err != nil {
		return nil, err
	}

	// Загружаем INI-файл
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load INI file: %v", err)
	}

	// Читаем секцию SQL
	sqlSection := cfg.Section("SQL")
	if sqlSection == nil {
		return nil, fmt.Errorf("section 'SQL' not found in INI file")
	}

	// Извлекаем параметры
	host := sqlSection.Key("host").String()
	port, err := sqlSection.Key("port").Int()
	if err != nil {
		return nil, fmt.Errorf("invalid port value: %v", err)
	}
	user := sqlSection.Key("username").String()
	password := sqlSection.Key("password").String()
	dbName := sqlSection.Key("dbname").String()

	// Возвращаем конфигурацию
	return &SQLConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}, nil
}

type handler struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) handler {
	return handler{db}
}

func Connect() *pgxpool.Pool {
	config, err := loadSQLConfig()
	if err != nil {
		log.Fatalf("Failed to load SQL config: %v", err)
	}
	connInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.DBName)
	pgxConfig, err := pgxpool.ParseConfig(connInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pgxConfig.MaxConns = 5
	pgxConfig.MaxConnLifetime = 5 * time.Minute
	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.HealthCheckPeriod = time.Minute
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Successfully connected to db %s!\n", connInfo)
	return db
}

func CloseConnection(db *pgxpool.Pool) {
	defer db.Close()
}
