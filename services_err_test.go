package vapi

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

var inMemoryServerWithError *fasthttputil.InmemoryListener

var apiServiceWithError *VAPI
var apiClientWithError http.Client

// DemoAPI area
type DemoAPIErr struct{}

// ErrorTest Method to test
func (h *DemoAPI) ErrorTest(ctx *fasthttp.RequestCtx, Args *TestArgs, Reply *TestReply) error {

	errs := &Error{
		ErrorHTTPCode: fasthttp.StatusFailedDependency,
		ErrorCode:     606,
		ErrorMessage:  "Test Wrong answer",
		Data:          nil,
	}

	return errs
}

func TestErrNewServer(t *testing.T) {
	apiServiceWithError = NewServer()
}

func TestErrVAPI_RegisterService(t *testing.T) {
	err := apiServiceWithError.RegisterService(new(DemoAPI), "demoerr")
	if err != nil {
		t.Error(err)
	}
}

func TestErrVAPI_GetServiceMap(t *testing.T) {
	tt, err := apiServiceWithError.GetServiceMap()
	if err != nil {
		t.Error(err)
	}
	if len(tt) != 0 {
		t.Error(fmt.Errorf("size of service map is higher that expected! Shoud be 0"))
	}
}

func TestErr(t *testing.T) {
	inMemoryServerWithError = fasthttputil.NewInmemoryListener()

	apiClientWithError = http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return inMemoryServerWithError.Dial()
			},
		},
	}

	reqHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api/demoerr.Test":
			apiServiceWithError.CallAPI(ctx, "demoerr.Test")
		case "/api/demoerr.ErrorTest":
			apiServiceWithError.CallAPI(ctx, "demoerr.ErrorTest")
		default:
			ctx.Error(fmt.Sprintf("Unsupported path: %s", ctx.Path()), fasthttp.StatusNotFound)
		}
	}

	go fasthttp.Serve(inMemoryServerWithError, reqHandler)
}

func TestErrVAPI_CallAPI_WrongAnswer(t *testing.T) {

	var jsonStr = []byte(`{"ID":"onomnomnom"}`)

	req, err := http.NewRequest("POST", "http://testerr/api/demoerr.ErrorTest", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error(err)
	}
	ress, err := apiClientWithError.Do(req)

	if ress.StatusCode != fasthttp.StatusFailedDependency {
		t.Error(fmt.Sprintf("wrong answer http status code received: %d", ress.StatusCode))
	}

	bodyS, err := ioutil.ReadAll(ress.Body)
	if err != nil {
		t.Error(err)
	}

	bodyStr := string(bodyS)

	if bodyStr != "{\"error\":{\"error_code\":606,\"error_msg\":\"Test Wrong answer\",\"data\":null}}" {
		t.Error(fmt.Sprintf("wrong answer received: %s", bodyStr))
	}
}
