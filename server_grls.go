package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

type ServerGrls struct {
	Port string
}

func (t *ServerGrls) run() {
	router := mux.NewRouter()
	router.HandleFunc(`/grls/{code:\w+}`, t.grlsToJsonFromCode)
	router.HandleFunc("/grls", t.grlsToJson)
	router.HandleFunc(`/except/{code:\w+}`, t.grlsExceptToJsonFromCode)
	router.HandleFunc("/except", t.grlsExceptToJson)
	router.HandleFunc("/", indexHandler)
	http.Handle("/", router)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		Logging()
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := "API GRLS"
	tmpl, _ := template.New("data").Parse("<h1>{{ .}}</h1>Примеры:<p>/grls - возвращает весь список<p>/except - возвращает список исключенных<p>/grls/{штрих-код} - возвращает список по штрих-коду<p>/except/{штрих-код} - возвращает список исключенных по штрихкоду")
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
		fmt.Fprint(w, "Not Found")
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
		fmt.Fprint(w, "Not Found")
	} else {
		fmt.Fprint(w, b)
	}
}
