package main

import (
	"log"
	"net/http"
	"os"

	_ "GoRestApi/docs"
	"GoRestApi/helpers"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

const api_url = "/api/v2/"

// var logger *log.Logger
// func getLogFilePath() (string, error) {
// 	// Получаем путь к исполняемому файлу
// 	exePath, err := os.Executable()
// 	if err != nil {
// 		return "", err
// 	}
// 	// Получаем директорию, в которой находится исполняемый файл
// 	exeDir := filepath.Dir(exePath)
// 	// Формируем путь к лог-файлу
// 	logFilePath := filepath.Join(exeDir, "logs", "GoRestApi.log")
// 	// Создаем директорию для логов, если она не существует
// 	logDir := filepath.Dir(logFilePath)
// 	if err := os.MkdirAll(logDir, 0755); err != nil {
// 		return "", err
// 	}
// 	return logFilePath, nil
// }
// func initLogger() {
// 	logFilePath, err := getLogFilePath()
// 	if err != nil {
// 		log.Fatalf("Failed to get log file path: %v", err)
// 	}
// 	// Открываем файл для записи логов
// 	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal("Failed to open log file:", err)
// 	}

// 	// Настраиваем логгер для записи в файл
// 	logger = log.New(file, "API: ", log.Ldate|log.Ltime|log.Lshortfile)
// }

func (h handler) group_check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleReqests(DB *pgxpool.Pool, http_port string) {
	h := New(DB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//swagger
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+http_port+"/docs/doc.json")),
	)
	//DataMatrix
	r.Get(api_url+"dm/dm", h.dm_get)
	r.With(httpin.NewInput(helpers.CodeFilterParams{})).Get(api_url+"dm/", h.dm_get_all)
	r.Patch(api_url+"dm/dm", h.dm_patch)
	r.Post(api_url+"dm/", h.dm_post)
	r.Delete(api_url+"dm/dm", h.dm_delete)
	//Aggregate
	r.Get(api_url+"aggregate/unit_id", h.aggr_get)
	r.With(httpin.NewInput(helpers.CodeFilterParams{})).Get(api_url+"aggregate/", h.aggr_get_all)
	r.Patch(api_url+"aggregate/unit_id", h.aggr_patch)
	r.Post(api_url+"aggregate/", h.aggr_post)
	r.Delete(api_url+"aggregate/unit_id", h.aggr_delete)
	r.Post(api_url+"build_aggregate/", h.aggr_build)
	//Can group check
	r.Get(api_url+"group_check", h.group_check)
	// r.Post(api_url+"events_log/", h.plug)
	// r.Patch(api_url+"task_count_data/{id}", h.plug)

	log.Fatal(http.ListenAndServe(":"+http_port, r))
}

// @title DMC GO Rest API
// @version 1.0.4.5
// @description Это Rest API для DMC на базе Go
// @BasePath /api/v2
func main() {
	// initLogger()
	http_port := "8555"
	if len(os.Args) > 1 {
		argsWithoutProg := os.Args[1]
		if argsWithoutProg != "" {
			http_port = argsWithoutProg
		}
	}
	log.Println("Using port " + http_port)
	DB := Connect()
	handleReqests(DB, http_port)
	defer CloseConnection(DB)
}
