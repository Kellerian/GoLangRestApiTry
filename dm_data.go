package main

import (
	"time"

	"github.com/guregu/null/v5"
)

type DmFilterParams struct {
	Taskid    *int     `in:"query=taskid" db:"taskid" goqu:"omitnil"`
	Gtin      *string  `in:"query=gtin" db:"gtin" goqu:"omitnil, omitempty"`
	Status    []string `in:"query=status[],status" db:"status" goqu:"omitnil, omitempty"`
	Aggregate *string  `in:"query=aggregate" db:"aggregate" goqu:"omitnil, omitempty"`
}

type DataMatrixCode struct {
	Status                 int         `db:"status" json:"status"`
	Gtin                   string      `db:"gtin" json:"gtin"`
	ScanDate               time.Time   `db:"scandate" json:"scandate" goqu:"omitnil"`
	Emission_date          time.Time   `db:"emission_date" json:"emission_date" goqu:"omitnil"`
	Emission_document      null.String `db:"emission_document" json:"emission_document" goqu:"omitnil"`
	Insert_by              string      `db:"insert_by" json:"insert_by" goqu:"omitnil"`
	Crpt_state             int         `db:"crpt_state" json:"crpt_state" goqu:"omitnil"`
	Taskid                 null.Int64  `db:"taskid" json:"taskid" goqu:"omitnil"`
	Dm_tail                string      `db:"dm_tail" json:"dm_tail" goqu:"omitnil"`
	Verification_device    null.String `db:"verification_device" json:"verification_device" goqu:"omitnil"`
	Aggregate              null.String `db:"aggregate" json:"aggregate" goqu:"omitnil"`
	Volume                 int         `db:"volume" json:"volume" goqu:"omitnil"`
	Weight                 int         `db:"weight" json:"weight" goqu:"omitnil"`
	Dm_91                  string      `db:"dm_91" json:"dm_91" goqu:"omitnil"`
	Dm_92                  string      `db:"dm_92" json:"dm_92" goqu:"omitnil"`
	Dm_93                  string      `db:"dm_93" json:"dm_93" goqu:"omitnil"`
	Dm_8005                string      `db:"dm_8005" json:"dm_8005" goqu:"omitnil"`
	Processing_in_document null.String `db:"processing_in_document" json:"processing_in_document" goqu:"omitnil"`
	Crpt_template_idx      int         `db:"crpt_template_idx" json:"crpt_template_idx" goqu:"omitnil"`
	Dm                     string      `db:"dm" json:"dm"`
}

type PatchDataMatrixCode struct {
	Status              *int       `json:"status,omitempty" db:"status" goqu:"omitnil"`
	ScanDate            *time.Time `json:"scandate,omitempty" db:"scandate" goqu:"omitnil"`
	Insert_by           *string    `json:"insert_by,omitempty" db:"insert_by" goqu:"omitnil"`
	Crpt_state          *int       `json:"crpt_state,omitempty" db:"crpt_state" goqu:"omitnil"`
	Taskid              *int       `json:"taskid,omitempty" db:"taskid" goqu:"omitnil"`
	Dm_tail             *string    `json:"dm_tail,omitempty" db:"dm_tail" goqu:"omitnil"`
	Verification_device *string    `json:"verification_device,omitempty" db:"verification_device" goqu:"omitnil"`
	Aggregate           *string    `json:"aggregate,omitempty" db:"aggregate" goqu:"omitnil"`
	Volume              *int       `json:"volume,omitempty" db:"volume" goqu:"omitnil"`
	Weight              *int       `json:"weight,omitempty" db:"weight" goqu:"omitnil"`
	Dm_91               *string    `json:"dm_91,omitempty" db:"dm_91" goqu:"omitnil"`
	Dm_92               *string    `json:"dm_92,omitempty" db:"dm_92" goqu:"omitnil"`
	Dm_93               *string    `json:"dm_93,omitempty" db:"dm_93" goqu:"omitnil"`
	Dm_8005             *string    `json:"dm_8005,omitempty" db:"dm_8005" goqu:"omitnil"`
}

type DmResponse struct {
	Success bool           `json:"success"`
	Data    DataMatrixCode `json:"data"`
}

type DmPagedResponse struct {
	Page          int              `json:"page"`
	Next_page     int              `json:"next_page"`
	Previous_page int              `json:"previous_page"`
	Total_records int              `json:"total_records"`
	Response      []DataMatrixCode `json:"response"`
}

type DmResponseList struct {
	Success bool             `json:"success"`
	Data    *DmPagedResponse `json:"data"`
}
