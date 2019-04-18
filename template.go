package dbdiff

import (
	"database/sql"
	"reflect"
)

type Operations interface {
	queryListByRowMapper(sql string, rowMapper RowMapper, out interface{}) error

	QuerySingle(sql string, out interface{}) error

	QuerySingleByMapper(sql string, rowMapper RowMapper, out interface{}) error

	QueryList(sql string, out interface{}) error

	QueryListByMapper(sql string, rowMapper RowMapper, out interface{}) error
}

type RowMapperResultSetExtractor struct {
	rowMapper RowMapper
}

func (extractor *RowMapperResultSetExtractor) ExtractData(rs *sql.Rows, out interface{}) error {
	if !AssertTypePtrOfSlice(out) {
		return &DataAccessError{Message: "input param out must be a slice", Err: nil}
	}

	var (
		rowNum = 0
		vArr   = reflect.ValueOf(out)
		elType = vArr.Type().Elem().Elem()
		vArrPt = vArr.Elem()
	)
	for rs.Next() {
		newEl := reflect.New(elType)
		err := extractor.rowMapper.MapRow(rs, rowNum, newEl.Interface())
		if err != nil {
			return err
		}
		vArrPt.Set(reflect.Append(vArrPt, reflect.Indirect(newEl)))
	}
	return nil
}

type DBTemplate struct {
	db *sql.DB
}

func NewDBTemplate(db *sql.DB) *DBTemplate {
	return &DBTemplate{db: db}
}

func (tpl *DBTemplate) queryListByRowMapper(sql string, rowMapper RowMapper, out interface{}) error {
	rs, err := tpl.db.Query(sql)
	if err != nil {
		return &DataAccessError{Message: "Db query error", Err: err}
	}
	extractor := RowMapperResultSetExtractor{rowMapper: rowMapper}
	defer rs.Close()
	err = extractor.ExtractData(rs, out)
	if err != nil {
		return &DataAccessError{Message: "row mapper result set extractor error", Err: err}
	}
	return nil
}

const COL_TAG_NAME = "col"

type defaultRowMapper4Struct struct {
	fieldMap map[string]string
	init     bool
}

func (drm *defaultRowMapper4Struct) MapRow(rs *sql.Rows, rowNum int, out interface{}) error {
	if !drm.init {
		drm.fieldMap = make(map[string]string)
		t := reflect.ValueOf(out).Type().Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := field.Tag.Get(COL_TAG_NAME)
			if len(name) != 0 {
				drm.fieldMap[name] = field.Name
			}
		}
		drm.init = true
	}

	colNames, err := rs.Columns()
	if err != nil {
		return &DataAccessError{Message: "get columns from result set error", Err: err}
	}
	var (
		//point on point
		values      = make([]interface{}, len(drm.fieldMap))
		v           = reflect.ValueOf(out)
		resetFields = make([]reflect.Value, len(drm.fieldMap))
	)
	for i, name := range colNames {
		fileName := drm.fieldMap[name]
		fieldV := v.Elem().FieldByName(fileName)
		if fieldV.Type().Kind() == reflect.Ptr {
			if fieldV.IsNil() {
				fieldV.Set(reflect.New(fieldV.Type().Elem()))
			}
			//field is point,use it's addr
			values[i] = fieldV.Addr().Interface()
		} else {
			//new point on point wiht field
			reflectValue := reflect.New(reflect.PtrTo(fieldV.Type()))
			reflectValue.Elem().Set(fieldV.Addr())
			values[i] = reflectValue.Interface()
			resetFields[i] = fieldV
		}
	}
	err = rs.Scan(values...)
	if err != nil {
		return &DataAccessError{Message: "error when scan", Err: err}
	}

	for i, field := range resetFields {
		if v := reflect.ValueOf(values[i]).Elem().Elem(); v.IsValid() {
			//reset value on field which is not point
			field.Set(v)
		}
	}

	if err != nil {
		return &DataAccessError{Message: "scan result set error", Err: err}
	}
	return nil
}

func (tpl *DBTemplate) QuerySingle(sql string, out interface{}) error {
	if !AssertTypePtrOfStruct(out) {
		return &DataAccessError{Message: "out param must be a ptr of struct"}
	}
	return tpl.QuerySingleByMapper(sql, &defaultRowMapper4Struct{}, out)
}

func (tpl *DBTemplate) QuerySingleByMapper(sql string, rowMapper RowMapper, out interface{}) error {
	var (
		v        = reflect.ValueOf(out)
		slv      = sliceByType(v.Type().Elem())
		outSlice = slv.Interface()
		err      = tpl.queryListByRowMapper(sql, rowMapper, outSlice)
	)
	if err != nil {
		return &DataAccessError{Message: "query by mapper error", Err: err}
	}
	outSliceV := reflect.ValueOf(outSlice)
	len := outSliceV.Elem().Len()
	if len == 0 {
		return nil
	}

	if len != 1 {
		return &DataAccessError{Message: "incorrect result size", Err: err}
	}
	v.Elem().Set(outSliceV.Elem().Index(0))
	return nil
}

func sliceByType(tp reflect.Type) reflect.Value {
	return reflect.New(reflect.SliceOf(tp))
}

func (tpl *DBTemplate) QueryList(sql string, out interface{}) error {
	if !AssertTypePtrOfSliceWithStruct(out) {
		return &DataAccessError{Message: "out param must be a ptr of slice with struct"}
	}
	return tpl.queryListByRowMapper(sql, &defaultRowMapper4Struct{}, out)
}

func (tpl *DBTemplate) QueryListByMapper(sql string, rowMapper RowMapper, out interface{}) error {
	return tpl.queryListByRowMapper(sql, rowMapper, out)
}
