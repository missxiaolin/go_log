package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_log/collect_log/defs"
	"go_log/collect_log/handles"
	"go_log/collect_log/util"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/log", handles.AddLog)

	return router
}

func main() {
	if err := util.Init(defs.CONFIG_PATH + "./rmq.json"); err != nil {
		fmt.Println(err)
	}

	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}
