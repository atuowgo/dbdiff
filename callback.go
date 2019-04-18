package dbdiff

import "database/sql"

type ResultSetExtractor interface {
	ExtractData(rs *sql.Rows, out interface{}) error
}

type RowMapper interface {
	MapRow(rs *sql.Rows, rowNum int, out interface{}) error
}
