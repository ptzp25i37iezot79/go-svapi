package vapi

import (
	"net/http"
	"os"
	"testing"

	"encoding/xml"
	"io/ioutil"
	"net/http/httptest"
)

var server *httptest.Server

func TestInitialize(t *testing.T) {
	os.Setenv("TESTING", "YES")

	Initialize("/v1", middlewareLog)
}

func TestRegisterService(t *testing.T) {
	Server.RegisterService(new(APITodo), "todo")
}

func TestRunService(t *testing.T) {
	server = httptest.NewServer(Server.GetRouter())
}

func TestGetNotFoundError(t *testing.T) {
	resp, err := http.Get(server.URL + "/404")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("Received non-404 response: %d\n", resp.StatusCode)
	}
}

func TestGetSimpleRequestWOParams(t *testing.T) {
	resp, err := http.Get(server.URL + "/v1/todo.get")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expected := "{\"response\":{\"title\":\"\",\"body\":\"\",\"tags\":null}}"
	if expected != string(actual) {
		t.Errorf("Expected the message '%s'\n", expected)
	}
}

func TestGetSimpleRequestWithParams(t *testing.T) {
	resp, err := http.Get(server.URL + "/v1/todo.get?title=Title_For_Todo&body=This_is_a_body_for_todo&tags[]=tag1&tags[]=tag2")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	actualURLParams, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedURLParams := "{\"response\":{\"title\":\"Title_For_Todo\",\"body\":\"This_is_a_body_for_todo\",\"tags\":[\"tag1\",\"tag2\"]}}"
	if expectedURLParams != string(actualURLParams) {
		t.Errorf("Expected the message '%s'\n", expectedURLParams)
	}
}

func TestGetSimpleRequestWithUnknownParams(t *testing.T) {
	resp, err := http.Get(server.URL + "/v1/todo.get?unknown_param=OMFGERROR&title=Title_For_Todo&body=This_is_a_body_for_todo&tags[]=tag1&tags[]=tag2")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Received non-400 response: %d\n", resp.StatusCode)
	}
	actualWrongURLParams, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedWrongURLParams := "{\"error\":{\"error_code\":444,\"error_msg\":\"schema: invalid path \\\"unknown_param\\\"\"}}"
	if expectedWrongURLParams != string(actualWrongURLParams) {
		t.Errorf("Expected the message '%s'\n", expectedWrongURLParams)
	}
}

func TestStopServer(t *testing.T) {
	server.Close()
}

/*
	ADDITIONAL FUNCTIONS FOR TESTS
*/
func middlewareLog(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//It will print RawQuery of every request
		//log.Printf("URL Raw Query %v", r.URL.RawQuery)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type APITodo struct{}

type APITodoArg struct {
	XMLName xml.Name `xml:"todo" json:"-"`
	Title   string   `schema:"title" json:"title" xml:"title"`
	Body    string   `schema:"body" json:"body" xml:"body"`
	Tags    []string `schema:"tags[]" json:"tags" xml:"tags"`
}

func (a *APITodo) Get(r *http.Request, Args *APITodoArg, Reply *APITodoArg) error {
	Reply.Tags = Args.Tags
	Reply.Title = Args.Title
	Reply.Body = Args.Body
	return nil
}
