package main

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

func queryToJson(db *sql.DB, query string, args ...interface{}) (string, error) {
	var objects []map[string]interface{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		columns, err := rows.ColumnTypes()
		if err != nil {
			return "", err
		}
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return "", err
		}

		objects = append(objects, object)
	}
	b, err := json.MarshalIndent(objects, "", "\t")
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func StringToJson(st map[string]string) string {
	b, err := json.MarshalIndent(st, "", "\t")
	if err != nil {
		return err.Error()
	} else {
		return string(b[:])
	}
}
