package vapi

import (
	"log"
	"net/http"
	"os"
	"testing"

	"encoding/xml"
	"net/http/httptest"
)

func TestInitialize(t *testing.T) {
	os.Setenv("TESTING", "YES")

	Initialize("/v1", middlewareLog)

	Server.RegisterService(new(APITodo), "todo")
	server := httptest.NewServer(Server.GetRouter())

	resp, err := http.Get(server.URL + "/404")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("Received non-404 response: %d\n", resp.StatusCode)
	}

	//resp, err = http.Get(server.URL + "/v1/todo.get")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if resp.StatusCode != 404 {
	//	t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	//}

	//expected := fmt.Sprintf("Visitor count: %d.", i)
	//actual, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if expected != string(actual) {
	//	t.Errorf("Expected the message '%s'\n", expected)
	//}

	defer server.Close()
}

func middlewareLog(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL Raw Query %v", r.URL.RawQuery)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type APITodo struct{}

type APITodoArg struct {
	XMLName xml.Name `xml:"todo" json:"-"`
	Title   string   `json:"title" xml:"title"`
	Body    string   `schema:"-" json:"body" xml:"body"`
	Tags    []string `schema:"tags[]" json:"tags" xml:"tags"`
}

func (a *APITodo) Get(r *http.Request, Args *APITodoArg, Reply *APITodoArg) error {
	Reply.Tags = Args.Tags
	Reply.Title = Args.Title
	Reply.Body = Args.Body
	return nil
}
