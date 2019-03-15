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

var inMemoryServer *fasthttputil.InmemoryListener

var apiService *VAPI
var apiClient http.Client

// DemoAPI area
type DemoAPI struct{}

// Test Method to test
func (h *DemoAPI) Test(ctx *fasthttp.RequestCtx, Args *TestArgs, Reply *TestReply) error {
	Reply.ID = Args.ID
	return nil
}

func TestNewServer(t *testing.T) {
	apiService = NewServer()
}

func TestVAPI_RegisterService(t *testing.T) {
	err := apiService.RegisterService(new(DemoAPI), "demo")
	if err != nil {
		t.Error(err)
	}
}

func TestVAPI_GetServiceMap(t *testing.T) {
	tt, err := apiService.GetServiceMap()
	if err != nil {
		t.Error(err)
	}
	if len(tt) != 0 {
		t.Error(fmt.Errorf("size of service map is higher that expected! Shoud be 0"))
	}
}

func Test(t *testing.T) {
	inMemoryServer = fasthttputil.NewInmemoryListener()

	apiClient = http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return inMemoryServer.Dial()
			},
		},
	}

	reqHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api/demo.Test":
			apiService.CallAPI(ctx, "demo.Test")
		case "/api/demo.ErrorTest":
			apiService.CallAPI(ctx, "demo.ErrorTest")
		default:
			ctx.Error(fmt.Sprintf("Unsupported path: %s", ctx.Path()), fasthttp.StatusNotFound)
		}
	}

	go fasthttp.Serve(inMemoryServer, reqHandler)
}

func TestVAPI_CallAPI(t *testing.T) {

	var jsonStr = []byte(`{"ID":"onomnomnom"}`)

	req, err := http.NewRequest("POST", "http://test/api/demo.Test", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	if res.StatusCode != 200 {
		t.Error(fmt.Sprintf("wrong answer http status code received: %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	if string(body) != "{\"response\":{\"ID\":\"onomnomnom\"}}" {
		t.Error(fmt.Sprintf("wrong answer received: %s", body))
	}
}
