package dbdiff

import (
	"fmt"
	"testing"
)

var s = SchemeSql{}
var dbName = "dbdiff"
var tableName = "student"

func TestSchemeSql_TableSchemeSql(t *testing.T) {
	fmt.Println(s.TableSchemeSql(dbName))
}

func TestSchemeSql_ColumnSchemeSql(t *testing.T) {
	fmt.Println(s.ColumnSchemeSql(dbName, tableName))
}

func TestSchemeSql_IndexSchemeSql(t *testing.T) {
	fmt.Println(s.IndexSchemeSql(tableName))
}
