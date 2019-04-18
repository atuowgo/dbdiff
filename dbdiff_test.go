package dbdiff

import (
	"fmt"
	"log"
	"testing"
)

/**
CREATE TABLE `student` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) DEFAULT NULL,
  `age2` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
)
 */

func getDBConnNew() *DBConn {
	return NewDBConn(MYSQL, "root", "root12345", "localhost", 3306, "sqldiff2")
}

func TestSqlDiff_ParseDiff(t *testing.T) {
	connOld := getDBConn()
	connNew := getDBConnNew()
	sqlDiff := NewSqlDiff()
	diffDataBase, err := sqlDiff.ParseDiff(connOld,connNew)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(diffDataBase)

}
