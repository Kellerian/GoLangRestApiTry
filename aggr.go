package main

import (
	"GoRestApi/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ggicci/httpin"
)

const AGGR_TABLE_NAME = "aggregates"

func (h handler) aggr_get_all(w http.ResponseWriter, r *http.Request) {
	q := r.Context().Value(httpin.Input).(*AggrFilterParams)

	p_size, err := helpers.GetPageSize(r)
	if err != nil {
		log.Println("Ошибка запроса", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p_num, err := helpers.GetPageNumber(r)
	if err != nil {
		log.Println("Ошибка запроса", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterMap := helpers.StructToMap(q)
	sql, _, _ := goqu.Select("*").From(AGGR_TABLE_NAME).Where(goqu.Ex(filterMap)).Order(goqu.C("unit_id").Asc()).Offset(helpers.CalcOffset(p_size, p_num)).Limit(uint(p_size)).ToSQL()
	fmt.Println(sql)
	aggrs := []Aggregate{}
	err = pgxscan.Select(r.Context(), h.DB, &aggrs, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	total_rows := helpers.RowCount{}
	sql_count, _, _ := goqu.Select(goqu.COUNT("*")).From(AGGR_TABLE_NAME).Where(goqu.Ex(filterMap)).ToSQL()
	err = pgxscan.Get(r.Context(), h.DB, &total_rows, sql_count)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	pages := AggrPagedResponse{Next_page: p_num + 1, Page: p_num, Previous_page: p_num - 1, Total_records: total_rows.Total, Response: aggrs}
	response := AggrResponseList{Success: true, Data: &pages}
	json.NewEncoder(w).Encode(response)
}

func (h handler) aggr_get(w http.ResponseWriter, r *http.Request) {
	unit_id := r.URL.Query().Get("unit_id")
	if unit_id == "" {
		http.Error(w, "некорректный запрос. Не передан параметр unit_id", http.StatusBadRequest)
		return
	}
	aggr := Aggregate{}
	sql, _, _ := goqu.Select("*").From(AGGR_TABLE_NAME).Where(goqu.Ex{"unit_id": unit_id}).ToSQL()
	err := pgxscan.Get(r.Context(), h.DB, &aggr, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := AggrResponse{Success: true, Data: aggr}
	json.NewEncoder(w).Encode(response)
}

func (h handler) aggr_patch(w http.ResponseWriter, r *http.Request) {
	unit_id := r.URL.Query().Get("unit_id")
	if unit_id == "" {
		http.Error(w, "некорректный запрос. Не передан параметр unit_id", http.StatusBadRequest)
		return
	}
	var aggr_struct PatchAggregate
	err := json.NewDecoder(r.Body).Decode(&aggr_struct)
	if err != nil {
		log.Println("Ошибка чтения тела запроса!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Запрос изменений кода", unit_id)
	sql, _, _ := goqu.Update(AGGR_TABLE_NAME).Set(aggr_struct).Where(goqu.Ex{"unit_id": unit_id}).Returning("*").ToSQL()
	fmt.Println(sql)
	aggr := Aggregate{}
	err = pgxscan.Get(r.Context(), h.DB, &aggr, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := AggrResponse{Success: true, Data: aggr}
	json.NewEncoder(w).Encode(response)
}

func (h handler) aggr_post(w http.ResponseWriter, r *http.Request) {
	var aggr_struct Aggregate
	err := json.NewDecoder(r.Body).Decode(&aggr_struct)
	if err != nil {
		log.Println("Ошибка чтения тела запроса!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Запрос добавления кода", aggr_struct.Unit_id)
	sql, _, _ := goqu.Insert(AGGR_TABLE_NAME).Rows(aggr_struct).Returning("*").ToSQL()
	fmt.Println(sql)
	aggr := Aggregate{}
	err = pgxscan.Get(r.Context(), h.DB, &aggr, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := AggrResponse{Success: true, Data: aggr}
	json.NewEncoder(w).Encode(response)
}

func (h handler) aggr_delete(w http.ResponseWriter, r *http.Request) {
	unit_id := r.URL.Query().Get("unit_id")
	if unit_id == "" {
		http.Error(w, "некорректный запрос. Не передан параметр unit_id", http.StatusBadRequest)
		return
	}
	log.Println("Запрос удаления кода", unit_id)
	sql, _, _ := goqu.Delete(AGGR_TABLE_NAME).Where(goqu.Ex{"unit_id": unit_id}).ToSQL()
	fmt.Println(sql)
	_, err := h.DB.Query(r.Context(), sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h handler) aggr_build(w http.ResponseWriter, r *http.Request) {
	var aggr_build BuildAggregate
	err := json.NewDecoder(r.Body).Decode(&aggr_build)
	if err != nil {
		log.Println("Ошибка чтения тела запроса!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var tbl_name string
	var parent_field string
	var search_field string
	if aggr_build.Level == 0 {
		tbl_name = "dm"
		search_field = "dm"
		parent_field = "aggregate"
	} else {
		tbl_name = "aggregates"
		search_field = "unit_id"
		parent_field = "parent_id"
	}
	sql, _, _ := goqu.Update(tbl_name).Set(goqu.Ex{parent_field: aggr_build.Parent}).Where(goqu.Ex{search_field: aggr_build.Content}).ToSQL()
	fmt.Println(sql)
	res, err := h.DB.Exec(r.Context(), sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := BuildAggregateResponse{Success: true, Data: AggrLenRecords{LenRecords: res.RowsAffected()}}
	json.NewEncoder(w).Encode(response)
}
