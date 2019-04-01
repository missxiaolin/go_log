package handles

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go_log/collect_log/dbops"
	"go_log/collect_log/defs"
	"go_log/collect_log/response"
	"go_log/collect_log/util"
	"io/ioutil"
	"net/http"
)

func AddLog(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.Log{}
	if err := json.Unmarshal(res, ubody); err != nil {
		response.ApiErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	ubody.Ip = util.RemoteIp(r)

	// 记录到数据库
	if err := dbops.AddLog(ubody); err != nil {
		response.ApiErrorResponse(w, defs.ErrorDBError)
		return
	}
	response.ApiNormalResponse(w, "success", 201)
}
