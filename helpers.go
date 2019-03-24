package vapi

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

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

// WriteResponse write response to client with status code and server response struct
func WriteResponse(ctx *fasthttp.RequestCtx, status int, resp ServerResponse) {
	body, err := resp.MarshalJSON()
	if err != nil {
		ctx.SetBody([]byte(fmt.Sprintf(`{"error": {"error_code": 0, "error_msg": "can't marshal response", "dara": "%s"}}`, err.Error())))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		ctx.SetBody(body)
		ctx.SetStatusCode(status)
	}
	ctx.Response.Header.Set("x-content-type-options", "nosniff")
	ctx.SetContentType("application/json; charset=utf-8")
}
