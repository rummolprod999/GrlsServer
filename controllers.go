package main

import (
	"html/template"
	"net/http"
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

func (t *ServerGrls) grlsDBUpdate(w http.ResponseWriter, r *http.Request) {
	t.updateDB(w, r)
}
