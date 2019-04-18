package dbdiff

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DriverType string

func (dt DriverType) String() string {
	return string(dt)
}

const (
	MYSQL DriverType = "mysql"
)

type DBConn struct {
	DriverName DriverType
	Username   string
	Password   string
	Ip         string
	Port       int
	DBName     string
}

func NewDBConn(driverName DriverType, username, password, ip string, port int, dbName string) *DBConn {
	return &DBConn{
		DriverName: driverName,
		Username:   username,
		Password:   password,
		Ip:         ip,
		Port:       port,
		DBName:     dbName,
	}
}

func (dbConn *DBConn) ConnUrl() (string, error) {
	switch dbConn.DriverName {
	case MYSQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConn.Username, dbConn.Password, dbConn.Ip, dbConn.Port, dbConn.DBName), nil
	default:
		return "", &DBNotSupportError{DriverName: dbConn.DriverName.String()}
	}
}

func (dbConn *DBConn) Conn() (*sql.DB, error) {
	connUrl, err := dbConn.ConnUrl()
	if err != nil {
		return nil, err
	}
	return sql.Open(dbConn.DriverName.String(), connUrl)
}
