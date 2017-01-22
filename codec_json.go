package vapi

import (
	"encoding/json"
	"net/http"
)

// serverResponse represents a JSON-RPC response returned by the server.
type serverResponseJSON struct {
	// The Object that was returned by the invoked method. This must be null
	// in case there was an error invoking the method.
	// As per spec the member will be omitted if there was an error.
	//Result interface{} `json:"response,omitempty"`

	// An Error object if there was an error invoking the method. It must be
	// null if there was no error.
	// As per spec the member will be omitted if there was no error.
	//Error *errorJSON `json:"error,omitempty"`
}

type errorJSON struct {
	Code    int         `json:"error_code"`
	Message string      `json:"error_msg"`            /* required */ // A Primitive or Structured value that contains additional information about the error.
	Data    interface{} `json:"error_data,omitempty"` /* optional */
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
func (c *serverResponseJSON) WriteResponse(w http.ResponseWriter, reply interface{}) {
	res := struct {
		Response interface{} `json:"response"`
	}{reply}
	c.writeServerResponse(w, 200, res)
}

// WriteError encodes the error response and writes it to the ResponseWriter.
func (c *serverResponseJSON) WriteError(w http.ResponseWriter, status int, err error) {
	res := struct {
		Error interface{} `json:"error"`
	}{errorJSON{Code: err.(*Error).Code, Message: err.(*Error).Message, Data: err.(*Error).Data}}
	c.writeServerResponse(w, status, res)
}

func (c *serverResponseJSON) writeServerResponse(w http.ResponseWriter, status int, res interface{}) {
	b, err := json.Marshal(res)
	if err != nil {
		writePureError(w, 400, err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(b)
}
