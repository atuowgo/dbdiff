package dbdiff

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Scheme struct {
	DbConn    *DBConn
	Db        *sql.DB
	schemeSql *SchemeSql
	tpl       *DBTemplate
}

func NewScheme(dbConn *DBConn, db *sql.DB) *Scheme {
	var (
		tpl = NewDBTemplate(db)
	)
	return &Scheme{
		DbConn:    dbConn,
		Db:        db,
		schemeSql: &SchemeSql{},
		tpl:       tpl,
	}
}

func (scheme *Scheme) Parse() (*DataBase, error) {
	return scheme.parseDataBase(scheme.DbConn.DBName)
}

func (scheme *Scheme) parseDataBase(dbName string) (*DataBase, error) {
	dataBase := &DataBase{}
	options, err := scheme.parseOptions()
	if err != nil {
		return nil, err
	}

	dataBase.Options = options
	tables, err := scheme.parseTables()
	if err != nil {
		return nil, err
	}
	dataBase.Tables = tables

	return dataBase, nil
}

func (scheme *Scheme) parseOptions() ([]*Variable, error) {
	variableSchemes := []VariableScheme{}
	err := scheme.tpl.QueryList(scheme.schemeSql.VariablesSchemeSql(Session), &variableSchemes)
	if err != nil {
		return nil, err
	}
	options := make([]*Variable, len(variableSchemes))
	for i, variable := range variableSchemes {
		variable := &Variable{VariableScheme: variable}
		options[i] = variable
	}
	return options, nil
}

func (scheme *Scheme) parseTables() ([]*Table, error) {
	tableSchemes := []TableScheme{}
	err := scheme.tpl.QueryList(scheme.schemeSql.TableSchemeSql(scheme.DbConn.DBName), &tableSchemes)
	if err != nil {
		return nil, err
	}

	tables := make([]*Table, len(tableSchemes))
	for i, tableScheme := range tableSchemes {
		var (
			table     = &Table{TableScheme: tableScheme}
			tableName = tableScheme.TableName
		)

		columns, err := scheme.parseColumns(tableName)
		if err != nil {
			return nil, err
		}
		table.ColumnList = columns

		indexes, err := scheme.parseIndexes(tableName)
		if err != nil {
			return nil, err
		}
		table.IndexList = indexes

		createTableScheme := CreateTableScheme{}
		err = scheme.tpl.QuerySingle(scheme.schemeSql.ShowCreateTableSql(tableName), &createTableScheme)
		if err != nil {
			return nil, err
		}
		table.CreateTableSql = createTableScheme.CreateTable
		table.DropTableSql = scheme.schemeSql.DropTableSql(tableName)

		tables[i] = table
	}

	return tables, nil
}

func (scheme *Scheme) parseColumns(tableName string) ([]*Column, error) {
	columnSchemes := []ColumnScheme{}
	err := scheme.tpl.QueryList(scheme.schemeSql.ColumnSchemeSql(scheme.DbConn.DBName, tableName), &columnSchemes)
	if err != nil {
		return nil, err
	}

	columns := make([]*Column, len(columnSchemes))
	for i, columnScheme := range columnSchemes {
		column := &Column{ColumnScheme: columnScheme}
		column.fillAddColumnSql()
		column.fillModifyColumnSql()
		column.fillDropColumnSql()

		columns[i] = column
	}

	return columns, nil
}

func (scheme *Scheme) parseIndexes(tableName string) ([]*Index, error) {
	indexSchemes := []IndexScheme{}
	err := scheme.tpl.QueryList(scheme.schemeSql.IndexSchemeSql(tableName), &indexSchemes)
	if err != nil {
		return nil, err
	}
	indexMap := make(map[string]*Index)
	for i, _ := range indexSchemes {
		var (
			indexScheme = indexSchemes[i]
			keyName     = indexScheme.KeyName
			index       *Index
		)
		if _, ok := indexMap[keyName]; !ok {
			inner := &Index{
				TableName: tableName,
				KeyName:   keyName,
			}
			indexMap[keyName] = inner
			index = inner
		} else {
			index = indexMap[keyName]
		}
		index.Columns = append(index.Columns, indexScheme.ColumnName)
		index.ColumnIndex = append(index.ColumnIndex, &indexScheme)
	}

	indexes := make([]*Index, len(indexMap))
	var i = 0
	for _, index := range indexMap {
		index.fillAddIndexSql()
		index.fillDropIndexSql()
		indexes[i] = index
		i++
	}
	return indexes, nil
}

const SCHEME_KEY_COMPARATOR_TAG_NAME = "comp"

func SchemeKeyComparator(left, right interface{}) int {
	if reflect.ValueOf(left).Type().Elem().Name() !=
		reflect.ValueOf(right).Type().Elem().Name() {
		return -1
	}
	originKeys := make(map[string]interface{})
	keyValues(left, originKeys)

	destKeys := make(map[string]interface{})
	keyValues(right, destKeys)

	for k, sv := range originKeys {
		dv := destKeys[k]
		res := strings.Compare(sv.(string), dv.(string))
		if res != 0 {
			return res
		}
	}
	return 0
}

func keyValues(item interface{}, keys map[string]interface{}) {
	var (
		v = reflect.ValueOf(item)
		t = v.Type().Elem()
	)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		keyTag := field.Tag.Get(SCHEME_KEY_COMPARATOR_TAG_NAME)
		if "" != keyTag {
			keyVal := v.Elem().Field(i).Interface()
			keys[field.Name] = keyVal
		} else {
			if v.Elem().Field(i).Type().Kind() == reflect.Struct {
				keyValues(v.Elem().Field(i).Addr().Interface(), keys)
			}
		}
	}
}

type DataBase struct {
	Tables  []*Table
	Options []*Variable
}

type VariableScheme struct {
	VariableName string `col:"Variable_name" comp:"_"`
	Value        string `col:"Value"`
}

type Variable struct {
	VariableScheme
}

type TableScheme struct {
	TableName      string `col:"TABLE_NAME" comp:"_"`
	Engine         string `col:"ENGINE"`
	RowFormat      string `col:"ROW_FORMAT"`
	AutoIncrement  string `col:"AUTO_INCREMENT"`
	CreateOptions  string `col:"CREATE_OPTIONS"`
	TableCollation string `col:"TABLE_COLLATION"`
	TableComment   string `col:"TABLE_COMMENT"`
}

type CreateTableScheme struct {
	TableName   string `col:"Table"`
	CreateTable string `col:"Create Table"`
}

type Table struct {
	TableScheme

	CreateTableSql string
	DropTableSql   string

	ColumnList []*Column
	IndexList  []*Index
}

type ColumnScheme struct {
	TableName              string `col:"TABLE_NAME" comp:"_"`
	ColumnName             string `col:"COLUMN_NAME" comp:"_"`
	OrdinalPosition        int    `col:"ORDINAL_POSITION"`
	ColumnDefault          string `col:"COLUMN_DEFAULT"`
	NullAble               string `col:"IS_NULLABLE"`
	ColumnType             string `col:"COLUMN_TYPE"`
	ColumnKey              string `col:"COLUMN_KEY"`
	CharacterMaximumLength int    `col:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterSetName       string `col:"CHARACTER_SET_NAME"`
	CollationName          string `col:"COLLATION_NAME"`
	Extra                  string `col:"EXTRA"`
	ColumnComment          string `col:"COLUMN_COMMENT"`
}

type Column struct {
	ColumnScheme

	AddColumnSql    string
	DropColumnSql   string
	ModifyColumnSql string
}

func (column *Column) fillAddColumnSql() {
	var buff bytes.Buffer
	buff.WriteString("ALTER TABLE ")
	buff.WriteString(column.TableName)
	buff.WriteString(" ADD COLUMN ")
	buff.WriteString(column.ColumnName)
	buff.WriteString(" ")
	buff.WriteString(column.ColumnType)
	if "NO" == column.NullAble {
		buff.WriteString(" NOT NULL ")
	}
	if !AssertStrEmpty(column.ColumnDefault) {
		buff.WriteString(" DEFAULT ")
		if "CURRENT_TIMESTAMP" != strings.ToUpper(column.ColumnDefault) && !AssertStrEmpty(column.CollationName) {
			buff.WriteString(fmt.Sprintf("'%s'", column.ColumnDefault))
		} else {
			buff.WriteString(column.ColumnDefault)
		}
	}
	if !AssertStrEmpty(column.Extra) {
		buff.WriteString("")
		buff.WriteString(column.Extra)
	}

	if !AssertStrBlank(column.ColumnComment) {
		buff.WriteString(" COMMENT \"")
		buff.WriteString(column.ColumnComment)
		buff.WriteString("\"")
	}

	column.AddColumnSql = buff.String()
}

func (column *Column) fillModifyColumnSql() {
	var buff bytes.Buffer
	buff.WriteString("ALTER TABLE ")
	buff.WriteString(column.TableName)
	buff.WriteString(" MODIFY COLUMN ")
	buff.WriteString(column.ColumnName)
	buff.WriteString(" ")
	buff.WriteString(column.ColumnType)
	if "NO" == column.NullAble {
		buff.WriteString(" NOT NULL ")
	}
	if !AssertStrEmpty(column.ColumnDefault) {
		buff.WriteString(" DEFAULT ")
		if "CURRENT_TIMESTAMP" != strings.ToUpper(column.ColumnDefault) && !AssertStrEmpty(column.CollationName) {
			buff.WriteString(fmt.Sprintf("'%s'", column.ColumnDefault))
		} else {
			buff.WriteString(column.ColumnDefault)
		}
	}
	if !AssertStrEmpty(column.Extra) {
		buff.WriteString("")
		buff.WriteString(column.Extra)
	}

	if !AssertStrBlank(column.ColumnComment) {
		buff.WriteString(" COMMENT \"")
		buff.WriteString(column.ColumnComment)
		buff.WriteString("\"")
	}

	column.AddColumnSql = buff.String()
}

func (column *Column) fillDropColumnSql() {
	column.DropColumnSql = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", column.TableName, column.ColumnName)
}

type IndexScheme struct {
	TableName    string `col:"Table" comp:"_"`
	NonUnique    int    `col:"Non_unique"`
	KeyName      string `col:"Key_name" comp:"_"`
	SeqInIndex   int    `col:"Seq_in_index"`
	ColumnName   string `col:"Column_name"`
	Collation    string `col:"Collation"`
	Cardinality  string `col:"Cardinality"`
	SubPart      string `col:"Sub_part"`
	Packed       string `col:"Packed"`
	Null         string `col:"Null"`
	IndexType    string `col:"Index_type"`
	Comment      string `col:"Comment"`
	IndexComment string `col:"Index_comment"`
}

type Index struct {
	TableName   string
	KeyName     string
	Columns     []string
	ColumnIndex []*IndexScheme

	AddIndexSql    string
	DropIndexSql   string
	ModifyIndexSql string
}

func (index *Index) fillAddIndexSql() {
	var buff bytes.Buffer
	buff.WriteString("ALTER TABLE ")
	buff.WriteString(index.TableName)

	if "PRIMARY" == strings.ToUpper(index.KeyName) {
		buff.WriteString(" ADD PRIMARY KEY ")
	} else if len(index.ColumnIndex) != 0 && index.ColumnIndex[0].NonUnique == 0 {
		buff.WriteString(fmt.Sprintf(" ADD UNIQUE INDEX `%s`", index.KeyName))
	} else {
		buff.WriteString(fmt.Sprintf(" ADD INDEX `%s`", index.KeyName))
	}
	buff.WriteString(" (")
	for i := 0; i < len(index.ColumnIndex); i++ {
		if i != 0 {
			buff.WriteString(" , ")
		}
		buff.WriteString(fmt.Sprintf(`%s`, index.ColumnIndex[i].ColumnName))
	}
	buff.WriteString(")")

	index.AddIndexSql = buff.String()
}

func (index *Index) fillDropIndexSql() {
	var buff bytes.Buffer
	buff.WriteString("ALTER TABLE ")
	buff.WriteString(index.TableName)
	if "PRIMARY" == strings.ToUpper(index.KeyName) {
		buff.WriteString(" DROP PRIMARY KEY")
	} else {
		buff.WriteString(fmt.Sprintf(" DROP INDEX `%s`", index.KeyName))
	}

	index.DropIndexSql = buff.String()
}

func (index *Index) fillModifyIndexSql() {
}
