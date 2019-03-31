package main

import (
	"github.com/julienschmidt/httprouter"
	"go_log/collect_log/handles"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/log", handles.AddLog)

	return router
}

func main()  {
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}