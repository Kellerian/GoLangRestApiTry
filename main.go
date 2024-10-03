package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host     = "localhost"
	port     = 6432
	user     = "postgres"
	password = "postgres"
	dbname   = "wip_1044_1"
	api_url  = "/api/v2/"
)

func handleReqests(DB *pgxpool.Pool, http_port string) {
	h := New(DB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//DataMatrix
	r.Get(api_url+"dm/dm", h.dm_get)
	r.With(httpin.NewInput(DmFilterParams{})).Get(api_url+"dm", h.dm_get_all)
	r.Patch(api_url+"dm/dm", h.dm_patch)
	r.Post(api_url+"dm/dm", h.dm_post)
	r.Delete(api_url+"dm/dm", h.dm_delete)
	//Aggregate
	r.Get(api_url+"aggregate/unit_id", h.aggr_get)
	r.With(httpin.NewInput(AggrFilterParams{})).Get(api_url+"aggregate", h.aggr_get_all)
	r.Patch(api_url+"aggregate/unit_id", h.aggr_patch)
	r.Post(api_url+"aggregate/unit_id", h.aggr_post)
	r.Delete(api_url+"aggregate/unit_id", h.aggr_delete)
	r.Post(api_url+"build_aggregate", h.aggr_build)

	r.Post(api_url+"events_log/", h.plug)
	r.Patch(api_url+"task_count_data/{id}", h.plug)

	log.Fatal(http.ListenAndServe(":"+http_port, r))
}

func main() {
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
