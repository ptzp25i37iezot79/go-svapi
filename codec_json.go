//go:generate easyjson $GOFILE

package vapi

import (
	"encoding/json"
)

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// ServerResponse represents a JSON-RPC response returned by the server.
// easyjson:json
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

// Error ...
// easyjson:json
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

// Marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// TestArgs args for tests
// easyjson:json
type TestArgs struct {
	ID  string `json:"id,omitempty"`
	Ttt string `json:"ttt,omitempty"`
}

// TestReply reply for tests
// easyjson:json
type TestReply struct {
	ID  string `json:"id,omitempty"`
	Ttt string `json:"ttt,omitempty"`
}
