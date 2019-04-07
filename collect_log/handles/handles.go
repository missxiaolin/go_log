package handles

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
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
	//if err := dbops.AddLog(ubody); err != nil {
	//	response.ApiErrorResponse(w, defs.ErrorDBError)
	//	return
	//}




	go func(body interface{}) {
		data, err := json.Marshal(ubody)
		if err != nil {
			fmt.Println("json 解析错误")
		}
		// mq 推送
		if err := util.Push("myPusher", "myQueue", []byte(data)); err != nil {
			fmt.Println(err)
		}
	}(ubody)

	//if err := util.Fini(); err != nil {
	//	fmt.Println(err)
	//}

	response.ApiNormalResponse(w, "success", 200)
}
