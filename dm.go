package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"GoRestApi/helpers"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ggicci/httpin"
)

const DM_TABLE_NAME = "dm"

func (h handler) dm_get_all(w http.ResponseWriter, r *http.Request) {
	q := r.Context().Value(httpin.Input).(*DmFilterParams)

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
	sql, _, _ := goqu.Select("*").From(DM_TABLE_NAME).Where(goqu.Ex(filterMap)).Order(goqu.C("dm").Asc()).Offset(helpers.CalcOffset(p_size, p_num)).Limit(uint(p_size)).ToSQL()
	fmt.Println(sql)
	dms := []DataMatrixCode{}
	err = pgxscan.Select(r.Context(), h.DB, &dms, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	total_rows := helpers.RowCount{}
	sql_count, _, _ := goqu.Select(goqu.COUNT("*")).From(DM_TABLE_NAME).Where(goqu.Ex(filterMap)).ToSQL()
	err = pgxscan.Get(r.Context(), h.DB, &total_rows, sql_count)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	pages := DmPagedResponse{Next_page: p_num + 1, Page: p_num, Previous_page: p_num - 1, Total_records: total_rows.Total, Response: dms}
	response := DmResponseList{Success: true, Data: &pages}
	json.NewEncoder(w).Encode(response)
}

func (h handler) dm_get(w http.ResponseWriter, r *http.Request) {
	dm_id := r.URL.Query().Get("dm_id")
	if dm_id == "" {
		http.Error(w, "некорректный запрос. Не передан параметр dm_id", http.StatusBadRequest)
		return
	}
	dm := DataMatrixCode{}
	sql, _, _ := goqu.Select("*").From(DM_TABLE_NAME).Where(goqu.Ex{"dm": dm_id}).ToSQL()
	err := pgxscan.Get(r.Context(), h.DB, &dm, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := DmResponse{Success: true, Data: dm}
	json.NewEncoder(w).Encode(response)
}

func (h handler) dm_patch(w http.ResponseWriter, r *http.Request) {
	dm_id := r.URL.Query().Get("dm_id")
	if dm_id == "" {
		http.Error(w, "некорректный запрос. Не передан параметр dm_id", http.StatusBadRequest)
		return
	}
	var dm_struct PatchDataMatrixCode
	err := json.NewDecoder(r.Body).Decode(&dm_struct)
	if err != nil {
		log.Println("Ошибка чтения тела запроса!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Запрос изменений кода", dm_id)
	sql, _, _ := goqu.Update(DM_TABLE_NAME).Set(dm_struct).Where(goqu.Ex{"dm": dm_id}).Returning("*").ToSQL()
	fmt.Println(sql)
	dm := DataMatrixCode{}
	err = pgxscan.Get(r.Context(), h.DB, &dm, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := DmResponse{Success: true, Data: dm}
	json.NewEncoder(w).Encode(response)
}

func (h handler) dm_post(w http.ResponseWriter, r *http.Request) {
	var dm_struct DataMatrixCode
	err := json.NewDecoder(r.Body).Decode(&dm_struct)
	if err != nil {
		log.Println("Ошибка чтения тела запроса!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Запрос добавления кода", dm_struct.Dm)
	sql, _, _ := goqu.Insert(DM_TABLE_NAME).Rows(dm_struct).Returning("*").ToSQL()
	fmt.Println(sql)
	dm := DataMatrixCode{}
	err = pgxscan.Get(r.Context(), h.DB, &dm, sql)
	if err != nil {
		log.Println("Ошибка обмена с БД!", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := DmResponse{Success: true, Data: dm}
	json.NewEncoder(w).Encode(response)
}

func (h handler) dm_delete(w http.ResponseWriter, r *http.Request) {
	dm_id := r.URL.Query().Get("dm_id")
	if dm_id == "" {
		http.Error(w, "некорректный запрос. Не передан параметр dm_id", http.StatusBadRequest)
		return
	}
	log.Println("Запрос удаления кода", dm_id)
	sql, _, _ := goqu.Delete(DM_TABLE_NAME).Where(goqu.Ex{"dm": dm_id}).ToSQL()
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
