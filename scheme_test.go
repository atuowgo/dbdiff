package dbdiff

import (
	"fmt"
	"log"
	"testing"
)

func TestSchemeSqlResult(t *testing.T) {
	var s = SchemeSql{}
	var dbName = "dbdiff"
	var tableName = "student"

	db := getDB()
	defer db.Close()

	tpl := NewDBTemplate(db)
	table := []TableScheme{}
	err := tpl.QueryList(s.TableSchemeSql(dbName), &table)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(table)

	columns := []ColumnScheme{}
	err = tpl.QueryList(s.ColumnSchemeSql(dbName, tableName), &columns)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(columns)

	indexs := []IndexScheme{}
	err = tpl.QueryList(s.IndexSchemeSql(tableName), &indexs)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(indexs)
}

func TestScheme_Parse(t *testing.T) {
	var (
		dbConn = getDBConn()
		db     = getDB()
		scheme = NewScheme(dbConn, db)
	)

	database, err := scheme.Parse()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(database)
	cols := database.Tables[0].ColumnList
	fmt.Println(cols[0])
	fmt.Println(cols[1])
}
