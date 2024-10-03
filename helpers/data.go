package helpers

type PaginationParams struct {
	Size   int `in:"query=size"`
	Number int `in:"query=number"`
}

type RowCount struct {
	Total int `db:"total"`
}
