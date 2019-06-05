package main

import (
	"github.com/gorilla/mux"
	"net/http"
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
	router.HandleFunc("/dateup", t.grlsDateUpdate)
	router.HandleFunc("/", indexHandler)
	http.Handle("/", router)
	if err := http.ListenAndServe(t.Port, nil); err != nil {
		Logging(err)
		return
	}
}
