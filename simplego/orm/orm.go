package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type ormer struct {
	db *sql.DB
}

func (this *ormer) Select(object interface{}) error {
	if object == nil {
		return errors.New("can not select nil value")
	}

	paramsMap, tableName, pkName := getStructInformation(object)

	var whereSql string
	for key, value := range paramsMap {
		if key == pkName {
			whereSql = fmt.Sprintf("%s=%d", pkName, value)
		}
	}

	sql := fmt.Sprintf("select * from %s where %s", tableName, whereSql)
	fmt.Println(sql)

	rows, _ := this.db.Query(sql)
	cols, _ := rows.Columns()

	buff := make([]interface{}, len(cols))
	data := make([]string, len(cols))
	for i, _ := range buff {
		buff[i] = &data[i]
	}
	for rows.Next() {
		rows.Scan(buff...)
	}

	v := reflect.ValueOf(object)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		kind := value.Kind()

		trueValue := data[i]
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(trueValue, 10, 64)
			if err == nil {
				value.SetInt(val)
			}
		case reflect.String:
			value.SetString(trueValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(trueValue, 10, 64)
			if err == nil {
				value.SetUint(val)
			}
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(trueValue, 64)
			if err == nil {
				value.SetFloat(val)
			}
		case reflect.Bool:
			val, err := strconv.ParseBool(trueValue)
			if err == nil {
				value.SetBool(val)
			}
		}
	}

	return nil
}

func (this *ormer) Insert(object interface{}) (int64, error) {
	if object == nil {
		return 0, errors.New("can not insert nil value")
	}
	paramsMap, tableName, _ := getStructInformation(object)

	keys := make([]string, 0)
	values := make([]interface{}, 0)
	preValues := make([]string, 0)
	for key, value := range paramsMap {
		keys = append(keys, key)
		values = append(values, value)
		preValues = append(preValues, "?")
	}

	nameSql := strings.Join(keys, ",")
	preValueSql := strings.Join(preValues, ",")

	sql := fmt.Sprintf("insert into %s (%s)values(%s)", tableName, nameSql, preValueSql)
	fmt.Println(sql, values)

	stmt, err := this.db.Prepare(sql)

	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(values...)

	if err != nil {
		return 0, err
	}

	affected, err := res.LastInsertId()
	return affected, err
}

func getStructInformation(object interface{}) (paramsMap map[string]interface{}, tableName, pkName string) {
	if object == nil {
		return
	}

	paramsMap = make(map[string]interface{}, 0)

	t := reflect.TypeOf(object)
	v := reflect.ValueOf(object)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	tableName = t.Name()

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		tag := t.Field(i).Tag
		value := v.FieldByName(name)
		kind := value.Kind()

		if tag == "PK" {
			pkName = name
		}

		var trueValue interface{}
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			trueValue = value.Int()
		case reflect.String:
			trueValue = value.String()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			trueValue = value.Uint()
		case reflect.Float32, reflect.Float64:
			trueValue = value.Float()
		case reflect.Bool:
			trueValue = value.Bool()
		}
		paramsMap[name] = trueValue
	}

	return
}
