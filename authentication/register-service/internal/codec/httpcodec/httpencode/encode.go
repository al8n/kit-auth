package httpencode

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/al8n/kit-auth/authentication/register-service/internal/codec/httpcodec"
	"github.com/al8n/kit-auth/models/responses"
	"github.com/al8n/kit-auth/utils"
	utilstransport "github.com/al8n/kit-auth/utils/transport"
	"io/ioutil"
	"net/http"
)

func AuthenticationResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error  {
	res, ok := resp.(responses.AuthenticationResponse)
	if !ok {
		httpcodec.ErrorEncoder(
			ctx,
			utilstransport.ErrorCodecCasting(
				"Authentication",
				utilstransport.Response,
				utilstransport.HTTP),
			w)
		return nil
	}
	if res.Error != "" {
		httpcodec.ErrorEncoder(
			ctx,
			utils.Str2Err(res.Error),
			w)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

func SendOTPResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error  {
	res, ok := resp.(responses.SendOTPResponse)
	if !ok {
		httpcodec.ErrorEncoder(
			ctx,
			utilstransport.ErrorCodecCasting(
				"SendOTP",
				utilstransport.Response,
				utilstransport.HTTP),
			w)
		return nil
	}
	if res.Error != "" {
		httpcodec.ErrorEncoder(
			ctx,
			utils.Str2Err(res.Error),
			w)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

// GenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func GenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
