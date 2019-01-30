package vapi

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

// isExported returns true of a string is an exported (upper case) name.
func isExported(name string) bool {
	runez, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(runez)
}

// isExportedOrBuiltin returns true if a type is exported or a builtin.
func isExportedOrBuiltin(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// ReadRequestParams getting request parametrs
func readRequestParams(ctx *fasthttp.RequestCtx, args interface{}) error {
	return ffjson.Unmarshal(ctx.Request.Body(), args)
}

// WriteResponse write response to client with status code and server response struct
func WriteResponse(ctx *fasthttp.RequestCtx, status int, resp ServerResponse) {
	body, _ := ffjson.Marshal(resp)
	ctx.SetBody(body)
	ffjson.Pool(body)
	ctx.Response.Header.Set("x-content-type-options", "nosniff")
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(status)
}
