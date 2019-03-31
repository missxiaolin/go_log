package response

import (
	"encoding/json"
	"go_log/collect_log/defs"
	"io"
	"net/http"
)

func ApiErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	w.WriteHeader(errResp.HttpSC)

	resStr, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(resStr))
}

func ApiNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
