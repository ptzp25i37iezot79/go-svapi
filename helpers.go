package svapi

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/valyala/fasthttp"
)

const ContentTypeJson = "application/json; charset=utf-8"
const XmlContentTypeXml = "application/xml; charset=utf-8"
const XmlContentTypeXmlRss = "application/rss+xml; charset=utf-8"
const XmlContentTypeXmlAtom = "application/atom+xml; charset=utf-8"
const ContentTypeHtml = "text/html; charset=utf-8"
const ContentTypeProtobuf = "application/protobuf"

// isExported returns true of a string is an exported (upper case) name.
func isExported(name string) bool {
	runez, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(runez)
}

// SendPushFunction type that define function that used to send push
type ErrorHandlerFunction func(ctx *fasthttp.RequestCtx, err error)

func defaultErrorHandler(ctx *fasthttp.RequestCtx, err error) {
	WriteResponseString(ctx, fasthttp.StatusInternalServerError, ContentTypeHtml, fmt.Sprintf("go-svapi error: %v", err))
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
