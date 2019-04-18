package dbdiff

import (
	"database/sql"
	"fmt"
	"testing"
)

/**
CREATE TABLE `student` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
)
*/
type Student struct {
	Id   int    `col:"id"`
	Name string `col:"name"`
	Age  int    `col:"age"`
}

type StudentRowMapper struct {
}

func (this *StudentRowMapper) MapRow(rs *sql.Rows, rowNum int, out interface{}) error {
	s := out.(*Student)
	return rs.Scan(&s.Id, &s.Name, &s.Age)
}

func getDBConn() *DBConn {
	return NewDBConn(MYSQL, "root", "root12345", "localhost", 3306, "sqldiff")
}

func getDB() *sql.DB {
	dbconn := getDBConn()
	db, err := dbconn.Conn()
	if err != nil {
		panic(err)
	}
	return db
}

func TestDBTemplate_QueryList(t *testing.T) {
	db := getDB()
	defer db.Close()

	sql := "select * from student"
	tpl := NewDBTemplate(db)
	out := []Student{}
	err := tpl.QueryList(sql, &out)
	if err != nil {
		panic(err)
	}
	fmt.Println("out : ", out)

	outByMapper := []Student{}
	err = tpl.QueryListByMapper(sql, &StudentRowMapper{}, &outByMapper)
	if err != nil {
		panic(err)
	}
	fmt.Println("out by mapper", outByMapper)
}

func TestDBTemplate_QuerySingle(t *testing.T) {
	db := getDB()
	defer db.Close()

	sql := "select * from student"
	tpl := NewDBTemplate(db)
	out := Student{}
	err := tpl.QuerySingle(sql, &out)
	if err != nil {
		panic(err)
	}
	fmt.Println("out : ", out)

	outByMapper := Student{}
	err = tpl.QuerySingleByMapper(sql, &StudentRowMapper{}, &outByMapper)
	if err != nil {
		panic(err)
	}
	fmt.Println("out by mapper", outByMapper)
}
