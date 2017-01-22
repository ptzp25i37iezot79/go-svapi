package vapi

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

// ----------------------------------------------------------------------------
// Codec
// ----------------------------------------------------------------------------

// Codec creates a CodecRequest to process each request.
type Codec struct{}

// NewRequest returns a CodecRequest.
func (c *CodecRequest) NewRequest(r *http.Request, cResp codecServerResponseInterface) *CodecRequest {
	return newCodecRequest(r, cResp)
}

// ----------------------------------------------------------------------------
// CodecRequest
// ----------------------------------------------------------------------------

// CodecRequest decodes and encodes a single request.
type CodecRequest struct {
	Responser codecServerResponseInterface
	request   *serverRequest
	err       error
}

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// serverRequest represents a request received by the server.
type serverRequest struct {
	Method string
	Params url.Values
}

// newCodecRequest returns a new CodecRequest.
func newCodecRequest(r *http.Request, cResp codecServerResponseInterface) *CodecRequest {
	// Decode the request body and check if RPC method is valid.
	defer r.Body.Close()
	req := new(serverRequest)
	req.Method = r.Context().Value(KeyMethodID).(string)

	var errr error
	if r.Method == "POST" {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			errr = errors.New("Wrong Content-Type in request")
		} else {
			err := r.ParseForm()
			req.Params = r.PostForm
			if err != io.EOF {
				errr = err
			}
		}
	} else if r.Method == "GET" {
		req.Params = r.URL.Query()
	}
	return &CodecRequest{Responser: cResp, request: req, err: errr}
}

// Method returns the RPC method for the current request.
//
// The method uses a dotted notation as in "Service.Method".
func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

// ReadRequest fills the request object for the RPC method.
func (c *CodecRequest) ReadRequest(args interface{}) error {
	if c.err == nil {
		if c.request.Params != nil {
			c.err = schemaDecoder.Decode(args, c.request.Params)
		} else {
			c.err = errors.New("api: method request ill-formed: missing params field")
		}
	}
	if c.err != nil {
		c.err = &Error{Code: 444, Message: c.err.Error()}
	}
	return c.err
}

type codecServerResponseInterface interface {
	WriteResponse(http.ResponseWriter, interface{})
	WriteError(http.ResponseWriter, int, error)
}
