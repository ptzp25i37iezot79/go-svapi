//go:generate easyjson $GOFILE

package vapi

import (
	"encoding/json"
)

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// ServerResponse represents a JSON-RPC response returned by the server.
//easyjson:json
type ServerResponse struct {
	// The Object that was returned by the invoked method. This must be null
	// in case there was an error invoking the method.
	// As per spec the member will be omitted if there was an error.
	Response json.RawMessage `json:"response,omitempty"`

	// An Error object if there was an error invoking the method. It must be
	// null if there was no error.
	// As per spec the member will be omitted if there was no error.
	Error *Error `json:"error,omitempty"`
}

//Error ...
//easyjson:json
type Error struct {
	// A Number that indicates the error type that occurred.
	ErrorHTTPCode int `json:"-"`

	// A Number that indicates the error type that occurred.
	ErrorCode int `json:"error_code"`

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	ErrorMessage string `json:"error_msg"`

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"`
}

func (e *Error) Error() string {
	return e.ErrorMessage
}

//TestArgs args for tests
//easyjson:json
type TestArgs struct {
	ID string
}

//TestReply reply for tests
//easyjson:json
type TestReply struct {
	ID string
}
