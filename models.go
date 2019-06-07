package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
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

func (t *ServerGrls) grlsAll(w http.ResponseWriter, r *http.Request, table string) {
	defer SaveStack()
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	defer db.Close()
	b, err := queryToJson(db, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	fmt.Fprint(w, b)
}

func (t *ServerGrls) grlsListFromCode(w http.ResponseWriter, r *http.Request, table string) {
	defer SaveStack()
	w.Header().Set("Content-Type", "application/json")
	var params = []interface{}{}
	for i := 0; i < 20; i++ {
		value := r.FormValue(fmt.Sprintf("arr[%d]", i))
		if value != "" {
			params = append(params, strings.TrimSpace(value))
		}
	}
	if len(params) < 1 {
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Слишком мало агрументов в запросе"}))
		return
	}
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return

	}
	defer db.Close()
	query := "SELECT * FROM " + table + " WHERE code IN (?" + strings.Repeat(",?", len(params)-1) + ")"
	args := []interface{}{}
	args = append(args, params...)
	b, err := queryToJson(db, query, args...)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Not found"}))
	} else {
		fmt.Fprint(w, b)
	}
}

func (t *ServerGrls) grlsFromCode(w http.ResponseWriter, r *http.Request, table string) {
	defer SaveStack()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	code := strings.TrimSpace(vars["code"])
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	defer db.Close()
	b, err := queryToJson(db, fmt.Sprintf("SELECT * FROM %s WHERE code = $1", table), code)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Not found"}))
	} else {
		fmt.Fprint(w, b)
	}
}

func (t *ServerGrls) grlsDate(w http.ResponseWriter, r *http.Request, table string) {
	defer SaveStack()
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	defer db.Close()
	b, err := queryToJson(db, fmt.Sprintf("SELECT date_pub FROM %s LIMIT 1", table))
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Not found"}))
	} else {
		fmt.Fprint(w, b)
	}
}

func (t *ServerGrls) updateDB(w http.ResponseWriter, r *http.Request) {
	defer SaveStack()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	pass := strings.TrimSpace(vars["pass"])
	if pass == "" {
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Пуcтой параметр пароль"}))
		return
	}
	if pass != SecretKey {
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Неправильный пароль"}))
		return
	}
	fileExec := strings.TrimSpace(vars["file"])
	if fileExec == "" {
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Пустой параметр файл"}))
		return
	}
	cmd := exec.Command(fileExec)
	err := cmd.Start()
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	err = cmd.Wait()
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
		return
	}
	fmt.Fprint(w, StringToJson(map[string]string{"Ok": "Завершено успешно"}))
}
