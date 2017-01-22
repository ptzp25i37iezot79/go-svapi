package vapi

import (
	"log"
	"net/http"
	"os"
	"testing"

	"encoding/xml"
)

func TestInitialize(t *testing.T) {
	os.Setenv("TESTING", "YES")

	Initialize("/v1", middlewareLog)

	Server.RegisterService(new(APITodo), "todo")

	http.ListenAndServe(":8080", Server.GetRouter())
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
