package dbdiff

import (
	"fmt"
	"testing"
)

func TestDBConn_ConnUrl(t *testing.T) {
	dbConn := DBConn{
		DriverName: MYSQL,
		Username:   "user",
		Password:   "pass",
		Ip:         "localhost",
		Port:       3306,
		DBName:     "test",
	}
	fmt.Println(dbConn.ConnUrl())
}
