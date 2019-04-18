package dbdiff

import "fmt"

const (
	tableSchemeTpl = "SELECT TABLE_NAME, ENGINE, ROW_FORMAT, AUTO_INCREMENT,CREATE_OPTIONS, TABLE_COLLATION, " +
		"TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA='%s' " +
		"AND not isnull(ENGINE) order by TABLE_NAME"

	columnSchemeTpl = "SELECT TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE," +
		" COLUMN_TYPE, COLUMN_KEY, CHARACTER_MAXIMUM_LENGTH,CHARACTER_SET_NAME, COLLATION_NAME, EXTRA, " +
		"COLUMN_COMMENT  FROM information_schema.COLUMNS  WHERE TABLE_SCHEMA='%s'  AND TABLE_NAME='%s'  " +
		"ORDER BY ORDINAL_POSITION"

	indexSchemeTpl = "SHOW INDEX FROM %s"

	sessionVariablesSchemeTpl = "show variables"

	globalVariablesSchemeTpl = "show GLOBAL variables"

	showCreateTableTpl = "show CREATE TABLE %s"

	dropTableTpl = "drop TABLE %s IF EXISTS"
)

type VariableScope int

const (
	_ VariableScope = iota
	Session
	Global
)

type SchemeSql struct {
}

func (this *SchemeSql) TableSchemeSql(dbName string) string {
	return fmt.Sprintf(tableSchemeTpl, dbName)
}

func (this *SchemeSql) ColumnSchemeSql(dbName, tableName string) string {
	return fmt.Sprintf(columnSchemeTpl, dbName, tableName)
}

func (this *SchemeSql) IndexSchemeSql(tableName string) string {
	return fmt.Sprintf(indexSchemeTpl, tableName)
}

func (this *SchemeSql) VariablesSchemeSql(scope VariableScope) string {
	switch scope {
	case Session:
		return sessionVariablesSchemeTpl
	case Global:
		return globalVariablesSchemeTpl
	}
	return ""
}

func (this *SchemeSql) ShowCreateTableSql(tableName string) string {
	return fmt.Sprintf(showCreateTableTpl, tableName)
}

func (this *SchemeSql) DropTableSql(tableName string) string {
	return fmt.Sprintf(dropTableTpl, tableName)
}
