package helpers

type PaginationParams struct {
	Size   int `in:"query=size"`
	Number int `in:"query=number"`
}

type RowCount struct {
	Total int `db:"count"`
}

type CodeFilterParams struct {
	Taskid    *int     `in:"query=taskid" db:"taskid" goqu:"omitnil"`
	Gtin      *string  `in:"query=gtin" db:"gtin" goqu:"omitnil, omitempty"`
	Status    []string `in:"query=status[],status" db:"status" goqu:"omitnil, omitempty"`
	Dm        []string `in:"query=dm[],dm" db:"dm" goqu:"omitnil,omitempty"`
	Unit_id   []string `in:"query=unit_id[],unit_id" db:"unit_id" goqu:"omitnil,omitempty"`
	Aggregate *string  `in:"query=aggregate" db:"parent_id" goqu:"omitnil, omitempty"`
	Level     *int     `in:"query=level" db:"level" goqu:"omitnil"`
}
