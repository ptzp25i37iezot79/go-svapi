# go-vapi
vk styled api package for easy api developing

[Website](https://www.riftbit.com) | [Contributing](https://www.riftbit.com/How-to-Contribute)

[![license](https://img.shields.io/github/license/riftbit/go-vapi.svg)](LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/riftbit/go-vapi)
[![Coverage Status](https://coveralls.io/repos/github/riftbit/go-vapi/badge.svg?branch=master)](https://coveralls.io/github/riftbit/go-vapi?branch=master)
[![Build Status](https://travis-ci.org/riftbit/go-vapi.svg?branch=master)](https://travis-ci.org/riftbit/go-vapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/riftbit/go-vapi)](https://goreportcard.com/report/github.com/riftbit/go-vapi)
[![Release](https://img.shields.io/badge/release-v1.0.0-blue.svg?style=flat)](https://github.com/riftbit/go-vapi/releases)

## Installation

```
go get -u github.com/riftbit/go-vapi
```

## Usage
This is a minimal example.

```go
import (
	"net/http"
	"log"
	"github.com/riftbit/go-vapi"
)

func main() {
	//Initializing VAPI (required)
	vapi.Initialize("/v1") // This method can receive middlewares as additional params

	//Add Apis to VAPI
	vapi.Server.RegisterService(new(ApiTodo), "todo")

	//Add Routes to VAPI
	vapi.Server.AddRoute("GET", "/", http.FileServer(http.Dir("./views/")))
	vapi.Server.AddRoute("GET", "/uploads", http.FileServer(http.Dir("./uploads/")))
	vapi.Server.AddRoute("GET", "/static/js/", http.FileServer(http.Dir("./static/js/")))
	vapi.Server.AddRoute("GET", "/static/css/", http.FileServer(http.Dir("./static/css/")))
	vapi.Server.AddRoute("GET", "/static/img/", http.FileServer(http.Dir("./static/img/")))

	log.Println("Started server on port",":8080")
	log.Fatal(http.ListenAndServe(":8080", vapi.Server.GetRouter()))
}

type ApiTodo struct {}

type ApiTodo_arg struct {
	XMLName    xml.Name    `xml:"todo" json:"-"`
	Title      string      `json:"title" xml:"title"`
	Body       string      `schema:"-" json:"body" xml:"body"`
	Tags       []string    `schema:"tags[]" json:"tags" xml:"tags"`
}

func (self *ApiTodo) Get(r *http.Request, Args *ApiTodo_arg, Reply *ApiTodo_arg) error {
	Reply.Tags=Args.Tags
	Reply.Title=Args.Title
	Reply.Body=Args.Body
	return nil
}

```

*Initialize* method can receive additional middlewares as params

```go
    vapi.Initialize("/v1", middleware_log)
```

And middleware_log example is:

```go
    
func middleware_log(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL Raw Query %v", r.URL.RawQuery)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

```

And make request by curl (for json output)

```
curl -X POST -d 'title=New Todo' -d 'body=This is a new todo' -d 'tags[]=Todo' -d 'tags[]=Tag' http://127.0.0.1:8080/v1/todo.get
```

Or for XML output

```
curl -X POST -d 'title=New Todo' -d 'body=This is a new todo' -d 'tags[]=Todo' -d 'tags[]=Tag' http://127.0.0.1:8080/v1/todo.get?format=xml
```