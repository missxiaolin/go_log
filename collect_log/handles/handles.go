package handles

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	defs2 "go_log/collect_log/defs"
	"go_log/collect_log/response"
	"io/ioutil"
	"net/http"
)

func AddLog(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs2.Log{}
	if err := json.Unmarshal(res, ubody); err != nil {
		response.ApiErrorResponse(w, defs2.ErrorRequestBodyParseFailed)
		return
	}
	fmt.Println(ubody)
}
