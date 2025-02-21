package main

import (
	"GoRestApi/helpers"

	"github.com/guregu/null/v5"
)

type Aggregate struct {
	Status                 int                `db:"status" json:"status"`
	Aggr_gtin              null.String        `db:"aggr_gtin" json:"aggr_gtin"`
	Datetime               helpers.CustomTime `db:"datetime" json:"datetime" goqu:"omitnil"`
	Emission_date          helpers.CustomTime `db:"emission_date" json:"emission_date" goqu:"omitnil"`
	Emission_document      null.String        `db:"emission_document" json:"emission_document" goqu:"omitnil"`
	Crpt_state             int                `db:"crpt_state" json:"crpt_state" goqu:"omitnil"`
	Taskid                 null.Int64         `db:"taskid" json:"taskid" goqu:"omitnil"`
	Serial                 string             `db:"serial" json:"serial" goqu:"omitnil"`
	Parent_id              null.String        `db:"parent_id" json:"parent_id" goqu:"omitnil"`
	Volume                 int                `db:"volume" json:"volume" goqu:"omitnil"`
	Weight                 int                `db:"weight" json:"weight" goqu:"omitnil"`
	Dm_91                  string             `db:"dm_91" json:"dm_91" goqu:"omitnil"`
	Dm_92                  string             `db:"dm_92" json:"dm_92" goqu:"omitnil"`
	Dm_93                  string             `db:"dm_93" json:"dm_93" goqu:"omitnil"`
	Mrp                    string             `db:"mrp" json:"mrp" goqu:"omitnil"`
	Processing_in_document null.String        `db:"processing_in_document" json:"processing_in_document" goqu:"omitnil"`
	Crpt_template_idx      int                `db:"crpt_template_idx" json:"crpt_template_idx" goqu:"omitnil"`
	Unit_id                string             `db:"unit_id" json:"unit_id"`
	Level                  int                `db:"level" json:"level"`
	Aggregated             bool               `db:"aggregated" json:"aggregated"`
	Payment_on_emission    bool               `db:"payment_on_emission" json:"payment_on_emission"`
	Number_in_task         null.Int64         `db:"number_in_task" json:"number_in_task"`
	Place_in_warehouse     null.String        `db:"place_in_warehouse" json:"place_in_warehouse" goqu:"omitnil"`
}

type PatchAggregate struct {
	Status             *int                `db:"status" json:"status" goqu:"omitnil"`
	Datetime           *helpers.CustomTime `db:"datetime" json:"datetime" goqu:"omitnil"`
	Crpt_state         *int                `db:"crpt_state" json:"crpt_state" goqu:"omitnil"`
	Taskid             *null.Int64         `db:"taskid" json:"taskid" goqu:"omitnil"`
	Serial             *string             `db:"serial" json:"serial" goqu:"omitnil"`
	Parent_id          *null.String        `db:"parent_id" json:"parent_id" goqu:"omitnil"`
	Volume             *int                `db:"volume" json:"volume" goqu:"omitnil"`
	Weight             *int                `db:"weight" json:"weight" goqu:"omitnil"`
	Dm_91              *string             `db:"dm_91" json:"dm_91" goqu:"omitnil"`
	Dm_92              *string             `db:"dm_92" json:"dm_92" goqu:"omitnil"`
	Dm_93              *string             `db:"dm_93" json:"dm_93" goqu:"omitnil"`
	Level              *int                `db:"level" json:"level" goqu:"omitnil"`
	Number_in_task     *int                `db:"number_in_task" json:"number_in_task" goqu:"omitnil"`
	Place_in_warehouse *null.String        `db:"place_in_warehouse" json:"place_in_warehouse" goqu:"omitnil"`
}

type AggrResponse struct {
	Success bool      `json:"success"`
	Data    Aggregate `json:"data"`
}

type AggrPagedResponse struct {
	Page          int         `json:"page"`
	Next_page     int         `json:"next_page"`
	Previous_page int         `json:"previous_page"`
	Total_records int         `json:"total_records"`
	Size          int         `json:"size"`
	Response      []Aggregate `json:"response"`
}

type AggrResponseList struct {
	Success bool               `json:"success"`
	Data    *AggrPagedResponse `json:"data"`
}

type BuildAggregate struct {
	Level   int      `json:"level"`
	Parent  string   `json:"parent"`
	Content []string `json:"content"`
}

type AggrLenRecords struct {
	LenRecords int64 `json:"len_records"`
}

type BuildAggregateResponse struct {
	Success bool           `json:"success"`
	Data    AggrLenRecords `json:"data"`
}
