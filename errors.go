package dbdiff

import "fmt"

type DBNotSupportError struct {
	DriverName string
}

func (err *DBNotSupportError) Error() string {
	return fmt.Sprintf("%s not support", err.DriverName)
}

type DataAccessError struct {
	Message string
	Err     error
}

func (dae *DataAccessError) Error() string {
	return fmt.Sprintf("access data error:%s with %s", dae.Message, dae.Err.Error())
}
