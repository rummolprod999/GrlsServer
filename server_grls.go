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

type ServerGrls struct {
	Port string
}

func (t *ServerGrls) run() {
	router := mux.NewRouter()
	router.HandleFunc(`/grls/{code:\w+}`, t.grlsToJsonFromCode)
	router.HandleFunc("/grlslist", t.grlsListToJsonFromCode).Methods("POST")
	router.HandleFunc("/grls", t.grlsToJson)
	router.HandleFunc(`/except/{code:\w+}`, t.grlsExceptToJsonFromCode)
	router.HandleFunc("/exceptlist", t.grlsExceptListToJsonFromCode).Methods("POST")
	router.HandleFunc("/except", t.grlsExceptToJson)
	router.HandleFunc("/", indexHandler)
	http.Handle("/", router)
	if err := http.ListenAndServe(t.Port, nil); err != nil {
		Logging(err)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := "API GRLS"
	tmpl, _ := template.New("data").Parse("<h1>{{ .}}</h1>Примеры:<p>GET /grls - возвращает весь список<p>GET /except - возвращает список исключенных<p>GET /grls/{штрих-код} - возвращает список по штрих-коду<p>GET /except/{штрих-код} - возвращает список исключенных по штрихкоду<p>POST /grlslist arr[0] = {code}, до 20 элементов<p>POST /exceptlist arr[0] = {code}, до 20 элементов")
	tmpl.Execute(w, data)
}

func (t *ServerGrls) grlsToJson(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	b, err := queryToJson(db, "SELECT * FROM grls")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	fmt.Fprint(w, b)
}

func (t *ServerGrls) grlsExceptToJson(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	b, err := queryToJson(db, "SELECT * FROM grls_except")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	fmt.Fprint(w, b)
}

func (t *ServerGrls) grlsToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	b, err := queryToJson(db, "SELECT * FROM grls WHERE code = $1", code)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, `{"Error": "Not Found"}`)
	} else {
		fmt.Fprint(w, b)
	}

}

func (t *ServerGrls) grlsExceptToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	b, err := queryToJson(db, "SELECT * FROM grls_except WHERE code = $1", code)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, `{"Error": "Not Found"}`)
	} else {
		fmt.Fprint(w, b)
	}
}

func (t *ServerGrls) grlsListToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	var params = []interface{}{}
	for i := 0; i < 20; i++ {
		value := r.FormValue(fmt.Sprintf("arr[%d]", i))
		if value != "" {
			params = append(params, value)
		}
	}
	if len(params) < 1 {
		fmt.Fprint(w, `{"Error": "Слишком мало агрументов в запросе"}`)
		return
	}
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	query := "SELECT * FROM grls WHERE code IN (?" + strings.Repeat(",?", len(params)-1) + ")"
	args := []interface{}{}
	args = append(args, params...)
	b, err := queryToJson(db, query, args...)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, `{"Error": "Not Found"}`)
	} else {
		fmt.Fprint(w, b)
	}

}

func (t *ServerGrls) grlsExceptListToJsonFromCode(w http.ResponseWriter, r *http.Request) {
	var params = []interface{}{}
	for i := 0; i < 20; i++ {
		value := r.FormValue(fmt.Sprintf("arr[%d]", i))
		if value != "" {
			params = append(params, value)
		}
	}
	if len(params) < 1 {
		fmt.Fprint(w, `{"Error": "Слишком мало агрументов в запросе"}`)
		return
	}
	db, err := sql.Open("sqlite3", "file:grls.db?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	query := "SELECT * FROM grls_except WHERE code IN (?" + strings.Repeat(",?", len(params)-1) + ")"
	args := []interface{}{}
	args = append(args, params...)
	b, err := queryToJson(db, query, args...)
	if err != nil {
		Logging(err)
		fmt.Fprint(w, err.Error())
	}
	if b == "null" {
		w.WriteHeader(404)
		fmt.Fprint(w, `{"Error": "Not Found"}`)
	} else {
		fmt.Fprint(w, b)
	}

}
