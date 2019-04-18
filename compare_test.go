package dbdiff

import (
	"fmt"
	"testing"
)

func TestComparableSlice_Compare(t *testing.T) {
	columns1 := []*Column{
		&Column{ColumnScheme: ColumnScheme{TableName: "student", ColumnName: "1"}},
		&Column{ColumnScheme: ColumnScheme{TableName: "student", ColumnName: "2"}},
	}

	columns2 := []*Column{
		&Column{ColumnScheme: ColumnScheme{TableName: "student", ColumnName: "0"}},
		&Column{ColumnScheme: ColumnScheme{TableName: "student", ColumnName: "2"}},
		&Column{ColumnScheme: ColumnScheme{TableName: "student", ColumnName: "3"}},
	}
	diffColumns := []*DiffColumn{}
	aa := &diffItems{
		items: &diffColumns,
	}
	compSlice := KeySlice{
		keyCompareAction: aa,
		keyComparator:    SchemeKeyComparator,
	}
	compSlice.Compare(&columns1, &columns2)
	printColumns(diffColumns)

}

func printColumns(columns []*DiffColumn) {
	for i, _ := range columns {
		fmt.Println(columns[i].ItemNew, columns[i].ItemOld)
	}
	fmt.Println("---------")
}
