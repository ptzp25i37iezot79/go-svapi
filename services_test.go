package svapi

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

var inMemoryServer *fasthttputil.InmemoryListener

var apiService *SVAPI
var apiClient http.Client

func testErrorHandler(ctx *fasthttp.RequestCtx, err error) {
	WriteResponseString(ctx, fasthttp.StatusConflict, ContentTypeJSON, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
}

// DemoAPI area
type DemoAPI struct{}

// TestXml Method to test xml content type
func (h *DemoAPI) TestXml(ctx *fasthttp.RequestCtx) error {
	WriteResponseString(ctx, fasthttp.StatusOK, ContentTypeXML, "ok")
	return nil
}

// TestRss Method to test rss content type
func (h *DemoAPI) TestRss(ctx *fasthttp.RequestCtx) error {
	WriteResponseString(ctx, fasthttp.StatusOK, ContentTypeRssXML, "ok")
	return nil
}

// TestAtom Method to test content atom type
func (h *DemoAPI) TestAtom(ctx *fasthttp.RequestCtx) error {
	WriteResponseString(ctx, fasthttp.StatusOK, ContentTypeAtomXML, "ok")
	return nil
}

// TestJson Method to test content json type
func (h *DemoAPI) TestJson(ctx *fasthttp.RequestCtx) error {
	WriteResponseString(ctx, fasthttp.StatusOK, ContentTypeJSON, "ok")
	return nil
}

// TestHtml Method to test content html type
func (h *DemoAPI) TestHtml(ctx *fasthttp.RequestCtx) error {
	WriteResponseString(ctx, fasthttp.StatusOK, ContentTypeHTML, "ok")
	return nil
}

// TestProtobuf Method to test content protobuf type
func (h *DemoAPI) TestProtobuf(ctx *fasthttp.RequestCtx) error {
	WriteResponseString(ctx, fasthttp.StatusOK, ContentTypeProtobuf, "ok")
	return nil
}

// ErrorTest Method to test error response
func (h *DemoAPI) ErrorTest(ctx *fasthttp.RequestCtx) error {
	return fmt.Errorf("test error")
}

func TestNewServer(t *testing.T) {
	apiService = NewServer()
}

func TestVAPI_RegisterService(t *testing.T) {
	err := apiService.RegisterService(new(DemoAPI), "demo")
	assert.NoError(t, err)
}

func TestVAPI_GetServiceMap(t *testing.T) {
	serviceMap := apiService.GetServiceMap()
	assert.Equal(t, 7, len(serviceMap))
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
		case "/api/demo.TestXml":
			apiService.CallAPI(ctx, "demo.TestXml")
		case "/api/demo.TestJson":
			apiService.CallAPI(ctx, "demo.TestJson")
		case "/api/demo.TestHtml":
			apiService.CallAPI(ctx, "demo.TestHtml")
		case "/api/demo.TestRss":
			apiService.CallAPI(ctx, "demo.TestRss")
		case "/api/demo.TestAtom":
			apiService.CallAPI(ctx, "demo.TestAtom")
		case "/api/demo.TestProtobuf":
			apiService.CallAPI(ctx, "demo.TestProtobuf")
		case "/api/demo.ErrorTest":
			apiService.CallAPI(ctx, "demo.ErrorTest")
		default:
			ctx.Error(fmt.Sprintf("Unsupported path: %s", ctx.Path()), fasthttp.StatusNotFound)
		}
	}

	go fasthttp.Serve(inMemoryServer, reqHandler)
}

func TestVAPI_CallAPI_WrongAnswer(t *testing.T) {

	var jsonStr = []byte(`{"ID":"onomnomnom"}`)

	req, err := http.NewRequest("POST", "http://testerr/api/demo.ErrorTest", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusInternalServerError, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "svapi: error test error", string(bodyB))
	assert.Equal(t, ContentTypeHTML, res.Header.Get("Content-type"))
}

func TestVAPI_CallAPI_WrongAnswer_WithCustomErrorHandler(t *testing.T) {

	apiService.SetErrorHandlerFunction(testErrorHandler)

	req, err := http.NewRequest("POST", "http://testerr/api/demo.ErrorTest", bytes.NewBuffer([]byte(`{"ID":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusConflict, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"error": "test error"}`, string(bodyB))
	assert.Equal(t, ContentTypeJSON, res.Header.Get("Content-type"))

	apiService.SetErrorHandlerFunction(defaultErrorHandler)
}

func TestVAPI_CallAPI_Json(t *testing.T) {
	req, err := http.NewRequest("POST", "http://test/api/demo.TestJson", bytes.NewBuffer([]byte(`{"id":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", string(bodyB))
	assert.Equal(t, ContentTypeJSON, res.Header.Get("Content-type"))
}

func TestVAPI_CallAPI_Html(t *testing.T) {
	req, err := http.NewRequest("POST", "http://test/api/demo.TestHtml", bytes.NewBuffer([]byte(`{"id":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", string(bodyB))
	assert.Equal(t, ContentTypeHTML, res.Header.Get("Content-type"))
}

func TestVAPI_CallAPI_Xml(t *testing.T) {
	req, err := http.NewRequest("POST", "http://test/api/demo.TestXml", bytes.NewBuffer([]byte(`{"id":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", string(bodyB))
	assert.Equal(t, ContentTypeXML, res.Header.Get("Content-type"))
}

func TestVAPI_CallAPI_Rss(t *testing.T) {
	req, err := http.NewRequest("POST", "http://test/api/demo.TestRss", bytes.NewBuffer([]byte(`{"id":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", string(bodyB))
	assert.Equal(t, ContentTypeRssXML, res.Header.Get("Content-type"))
}

func TestVAPI_CallAPI_Atom(t *testing.T) {
	req, err := http.NewRequest("POST", "http://test/api/demo.TestAtom", bytes.NewBuffer([]byte(`{"id":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", string(bodyB))
	assert.Equal(t, ContentTypeAtomXML, res.Header.Get("Content-type"))
}

func TestVAPI_CallAPI_Protobuf(t *testing.T) {
	req, err := http.NewRequest("POST", "http://test/api/demo.TestProtobuf", bytes.NewBuffer([]byte(`{"id":"onomnomnom"}`)))
	if err != nil {
		t.Error(err)
	}
	res, err := apiClient.Do(req)

	assert.Equal(t, fasthttp.StatusOK, res.StatusCode)

	bodyB, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", string(bodyB))
	assert.Equal(t, ContentTypeProtobuf, res.Header.Get("Content-type"))
}
