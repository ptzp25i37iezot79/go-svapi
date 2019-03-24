package svapi

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/valyala/fasthttp"
)

const (
	// ContentTypeJson Content type for JSON
	ContentTypeJson = "application/json; charset=utf-8"

	// ContentTypeXml Content type for XML
	ContentTypeXml = "application/xml; charset=utf-8"

	// ContentTypeRss Content type for RSS
	ContentTypeRssXml = "application/rss+xml; charset=utf-8"

	// ContentTypeAtom Content type for ATOM
	ContentTypeAtomXml = "application/atom+xml; charset=utf-8"

	// ContentTypeHtml Content type for HTML
	ContentTypeHtml = "text/html; charset=utf-8"

	// ContentTypeProtobuf Content type for ProtoBuf
	ContentTypeProtobuf = "application/protobuf"
)

// isExported returns true of a string is an exported (upper case) name.
func isExported(name string) bool {
	runez, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(runez)
}

// SendPushFunction type that define function that used to send push
type ErrorHandlerFunction func(ctx *fasthttp.RequestCtx, err error)

func defaultErrorHandler(ctx *fasthttp.RequestCtx, err error) {
	WriteResponseString(ctx, fasthttp.StatusInternalServerError, ContentTypeHtml, fmt.Sprintf("svapi: error %v", err))
}

// WriteResponseBytes write response to client with status code, body and content type
func WriteResponseBytes(ctx *fasthttp.RequestCtx, status int, contentType string, resp []byte) {
	ctx.SetStatusCode(status)
	ctx.SetContentType(contentType)
	ctx.SetBody(resp)
}

// WriteResponseString write response to client with status code, body and content type
func WriteResponseString(ctx *fasthttp.RequestCtx, status int, contentType string, resp string) {
	ctx.SetStatusCode(status)
	ctx.SetContentType(contentType)
	ctx.SetBodyString(resp)
}
