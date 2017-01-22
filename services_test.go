package vapi

import (
	"log"
	"net/http"
	"os"
	"testing"

	"encoding/xml"

	"github.com/riftbit/vapi"
)

func TestInitialize(t *testing.T) {
	os.Setenv("TESTING", "YES")

	vapi.Initialize("/v1", middleware_log)

	vapi.Server.RegisterService(new(ApiTodo), "todo")

	http.ListenAndServe(":8080", vapi.Server.GetRouter())
}

func middleware_log(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL Raw Query %v", r.URL.RawQuery)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type ApiTodo struct{}

type ApiTodo_arg struct {
	XMLName xml.Name `xml:"todo" json:"-"`
	Title   string   `json:"title" xml:"title"`
	Body    string   `schema:"-" json:"body" xml:"body"`
	Tags    []string `schema:"tags[]" json:"tags" xml:"tags"`
}

func (self *ApiTodo) Get(r *http.Request, Args *ApiTodo_arg, Reply *ApiTodo_arg) error {
	Reply.Tags = Args.Tags
	Reply.Title = Args.Title
	Reply.Body = Args.Body
	return nil
}
