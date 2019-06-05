package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := "API GRLS"
	tmpl, _ := template.New("data").Parse("<h1>{{ .}}</h1>Примеры:<p>GET /grls - возвращает весь список<p>GET /except - возвращает список исключенных<p>GET /grls/{штрих-код} - возвращает список по штрих-коду<p>GET /except/{штрих-код} - возвращает список исключенных по штрихкоду<p>POST /grlslist arr[0] = {code}, до 20 элементов<p>POST /exceptlist arr[0] = {code}, до 20 элементов<p>GET /dateup - дата последнего обновления списка грлс")
	tmpl.Execute(w, data)
}

func (t *ServerGrls) grlsToJson(w http.ResponseWriter, r *http.Request) {
	t.grlsAll(w, r, GrlsTable)
}

func (t *ServerGrls) grlsExceptToJson(w http.ResponseWriter, r *http.Request) {
	t.grlsAll(w, r, GrlsExceptTable)
}

func (t *ServerGrls) grlsToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	t.grlsFromCode(w, r, GrlsTable)
}

func (t *ServerGrls) grlsExceptToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	t.grlsFromCode(w, r, GrlsExceptTable)
}

func (t *ServerGrls) grlsListToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	t.grlsListFromCode(w, r, GrlsTable)
}

func (t *ServerGrls) grlsExceptListToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	t.grlsListFromCode(w, r, GrlsExceptTable)
}

func (t *ServerGrls) grlsDateUpdate(w http.ResponseWriter, r *http.Request) {
	t.grlsDate(w, r, GrlsTable)
}

func (t *ServerGrls) grlsAll(w http.ResponseWriter, r *http.Request, table string) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	defer db.Close()
	b, err := queryToJson(db, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	fmt.Fprint(w, b)
}

func (t *ServerGrls) grlsListFromCode(w http.ResponseWriter, r *http.Request, table string) {
	w.Header().Set("Content-Type", "application/json")
	var params = []interface{}{}
	for i := 0; i < 20; i++ {
		value := r.FormValue(fmt.Sprintf("arr[%d]", i))
		if value != "" {
			params = append(params, value)
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

	}
	defer db.Close()
	query := "SELECT * FROM " + table + " WHERE code IN (?" + strings.Repeat(",?", len(params)-1) + ")"
	args := []interface{}{}
	args = append(args, params...)
	b, err := queryToJson(db, query, args...)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Not found"}))
	} else {
		fmt.Fprint(w, b)
	}
}

func (t *ServerGrls) grlsFromCode(w http.ResponseWriter, r *http.Request, table string) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	code := vars["code"]
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	defer db.Close()
	b, err := queryToJson(db, fmt.Sprintf("SELECT * FROM %s WHERE code = $1", table), code)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Not found"}))
	} else {
		fmt.Fprint(w, b)
	}
}

func (t *ServerGrls) grlsDate(w http.ResponseWriter, r *http.Request, table string) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	code := vars["code"]
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	defer db.Close()
	b, err := queryToJson(db, fmt.Sprintf("SELECT date_pub FROM %s LIMIT 1", table), code)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": err.Error()}))
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, StringToJson(map[string]string{"Error": "Not found"}))
	} else {
		fmt.Fprint(w, b)
	}
}
