package dbdiff

import (
	"fmt"
	"reflect"
	"strings"
)

type DBDiff struct {
}

func NewSqlDiff() *DBDiff {
	return &DBDiff{}
}

func (diff *DBDiff) ParseDiff(connOld, connNew *DBConn) (*DiffDataBase, error) {
	return diff.parseDiff(connOld, connNew)
}

func (diff *DBDiff) parseDiff(connOld, connNew *DBConn) (*DiffDataBase, error) {
	dataBaseOld, err := diff.newDatabase(connOld)
	if err != nil {
		return nil, err
	}
	dataBaseNew, err := diff.newDatabase(connNew)
	if err != nil {
		return nil, err
	}

	return diff.parseDatabaseDiff(dataBaseOld, dataBaseNew)
}

func (diff *DBDiff) newDatabase(conn *DBConn) (*DataBase, error) {
	if conn == nil {
		return nil, nil
	}

	db, err := conn.Conn()
	if err != nil {
		return nil, err
	}
	defer func() {
		db.Close()
	}()

	scheme := NewScheme(conn, db)
	return scheme.Parse()
}

func (diff *DBDiff) parseDatabaseDiff(databaseOld, dataBaseNew *DataBase) (*DiffDataBase, error) {
	if databaseOld == nil {
		return diff.copyDatabaseDiff(dataBaseNew, false), nil
	}

	if dataBaseNew == nil {
		return diff.copyDatabaseDiff(databaseOld, true), nil
	}

	diffDataBase := &DiffDataBase{}
	//diff tables
	diffTables := []*DiffTable{}
	tablesComp := KeySlice{
		keyCompareAction: &compDiffTables{
			items: &diffTables,
		},
		keyComparator: SchemeKeyComparator,
	}
	tablesComp.Compare(&databaseOld.Tables, &dataBaseNew.Tables)
	diffDataBase.DiffTables = diffTables

	//diff options
	diffOptions := []*DiffOption{}
	optionsComp := KeySlice{
		keyCompareAction: &diffItems{
			items: &diffOptions,
		},
		keyComparator: SchemeKeyComparator,
	}
	optionsComp.Compare(&databaseOld.Options, &dataBaseNew.Options)
	diffDataBase.DiffOptions = diffOptions

	return diffDataBase, nil
}

func (diff *DBDiff) copyDatabaseDiff(database *DataBase, isOld bool) *DiffDataBase {
	dataBase := &DiffDataBase{}
	dataBase.Copy(database, isOld)
	return dataBase
}

const (
	defaultItemOldName = "ItemOld"
	defaultItemNewName = "ItemNew"
)

type diffItems struct {
	items interface{}
}

func (this *diffItems) ActionBothExists(itemLeft, itemRight interface{}) {
	switch itemLeft.(type) {
	case *Column:
		var (
			left  = itemLeft.(*Column)
			right = itemRight.(*Column)
		)
		if left.ColumnScheme != right.ColumnScheme {
			diffColumn := &DiffColumn{
				ItemOld: left,
				ItemNew: right,
			}
			this.appendItem(diffColumn)
			fmt.Println()
		}
	case *Index:
		var (
			left  = itemLeft.(*Index)
			right = itemRight.(*Index)
		)
		srcCols := strings.Join(left.Columns, ",")
		destCols := strings.Join(right.Columns, ",")
		if srcCols != destCols {
			diffIndex := &DiffIndex{
				ItemOld: left,
				ItemNew: right,
			}
			this.appendItem(diffIndex)
		}
	case *Variable:
		var (
			left  = itemLeft.(*Variable)
			right = itemRight.(*Variable)
		)
		if left.VariableName != right.VariableName {
			diffOption := &DiffOption{
				ItemOld: left,
				ItemNew: right,
			}
			this.appendItem(diffOption)
		}
	}
}

func (this *diffItems) appendItem(item interface{}) {
	sliceValue := reflect.ValueOf(this.items).Elem()
	sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(item)))
}

func (this *diffItems) ActionLeftExists(itemLeft interface{}) {
	this.actionSingleExists(itemLeft, true)
}

func (this *diffItems) ActionRightExists(itemRight interface{}) {
	this.actionSingleExists(itemRight, false)
}

func (this *diffItems) actionSingleExists(item interface{}, isOld bool) {
	sliceValue := reflect.ValueOf(this.items).Elem()
	diffItemType := sliceValue.Type().Elem().Elem()
	diffItemValue := reflect.New(diffItemType)
	var itemName string
	if isOld {
		itemName = defaultItemOldName
	} else {
		itemName = defaultItemNewName
	}
	diffItemValue.Elem().FieldByName(itemName).Set(reflect.ValueOf(item))
	sliceValue.Set(reflect.Append(sliceValue, diffItemValue))
}

type DiffDataBase struct {
	DiffTables  []*DiffTable
	DiffOptions []*DiffOption
}

func (diff *DiffDataBase) Copy(database *DataBase, isOld bool) {
	tables := make([]*DiffTable, len(database.Tables))
	for i, _ := range database.Tables {
		diffTable := new(DiffTable)
		diffTable.Copy(database.Tables[i], isOld)
		tables[i] = diffTable
	}
	diff.DiffTables = tables

	options := make([]*DiffOption, len(database.Options))
	for i, _ := range database.Options {
		diffOption := new(DiffOption)
		diffOption.Copy(database.Options[i], isOld)
		options[i] = diffOption
	}
	diff.DiffOptions = options
}

type DiffTable struct {
	TableName   string
	TableOld    *Table
	TableNew    *Table
	DiffColumns []*DiffColumn
	DiffIndex   []*DiffIndex
}

type compDiffTables struct {
	items *[]*DiffTable
}

func (this *compDiffTables) ActionBothExists(itemLeft, itemRight interface{}) {
	var (
		left      = itemLeft.(*Table)
		right     = itemRight.(*Table)
		diffTable = &DiffTable{}
	)
	diffTable.TableName = left.TableName
	diffTable.TableOld = left
	diffTable.TableNew = right

	diffColumns := []*DiffColumn{}
	columnComp := KeySlice{
		keyCompareAction: &diffItems{
			items: &diffColumns,
		},
		keyComparator: SchemeKeyComparator,
	}
	columnComp.Compare(&left.ColumnList, &right.ColumnList)
	diffTable.DiffColumns = diffColumns

	diffIndex := []*DiffIndex{}
	indexComp := KeySlice{
		keyCompareAction: &diffItems{
			items: &diffIndex,
		},
		keyComparator: SchemeKeyComparator,
	}
	indexComp.Compare(&left.IndexList, &right.IndexList)
	diffTable.DiffIndex = diffIndex
	*this.items = append(*this.items, diffTable)
}
func (this *compDiffTables) ActionLeftExists(itemLeft interface{}) {
	var (
		table     = itemLeft.(*Table)
		diffTable = &DiffTable{}
	)
	diffTable.Copy(table, true)
	diffTable.TableName = table.TableName
	diffTable.TableOld = table
	*this.items = append(*this.items, diffTable)
}
func (this *compDiffTables) ActionRightExists(itemRight interface{}) {
	var (
		table     = itemRight.(*Table)
		diffTable = &DiffTable{}
	)
	diffTable.Copy(table, false)
	diffTable.TableName = table.TableName
	diffTable.TableOld = table
	*this.items = append(*this.items, diffTable)
}

func (diff *DiffTable) Copy(table *Table, isOld bool) {
	diff.TableName = table.TableName
	if isOld {
		diff.TableOld = table
	} else {
		diff.TableNew = table
	}
	columns := make([]*DiffColumn, len(table.ColumnList))
	for i, _ := range table.ColumnList {
		diffColumn := new(DiffColumn)
		diffColumn.Copy(table.ColumnList[i], isOld)
		columns[i] = diffColumn
	}
	diff.DiffColumns = columns

	indexes := make([]*DiffIndex, len(table.IndexList))
	for i, _ := range table.IndexList {
		diffIndex := new(DiffIndex)
		diffIndex.Copy(table.IndexList[i], isOld)
		indexes[i] = diffIndex
	}
	diff.DiffIndex = indexes
}

type DiffColumn struct {
	ItemOld *Column
	ItemNew *Column
}

func (diff *DiffColumn) Copy(column *Column, isOld bool) {
	copy(diff, column, isOld)
}

type DiffIndex struct {
	ItemOld *Index
	ItemNew *Index
}

func (diff *DiffIndex) Copy(index *Index, isOld bool) {
	copy(diff, index, isOld)
}

type DiffOption struct {
	ItemOld *Variable
	ItemNew *Variable
}

func (diff *DiffOption) Copy(option *Variable, isOld bool) {
	copy(diff, option, isOld)
}

func copy(owner, item interface{}, isOld bool) {
	v := reflect.ValueOf(owner)
	if isOld {
		v.Elem().FieldByName(defaultItemOldName).Set(reflect.ValueOf(item))
	} else {
		v.Elem().FieldByName(defaultItemNewName).Set(reflect.ValueOf(item))
	}
}
