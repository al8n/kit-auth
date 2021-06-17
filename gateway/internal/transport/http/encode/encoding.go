package encode

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func JSONRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func JSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	return json.NewEncoder(w).Encode(response)
}



