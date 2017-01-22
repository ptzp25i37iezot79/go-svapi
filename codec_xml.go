package vapi

import (
	"net/http"
	"encoding/xml"
)


// serverResponse represents a JSON-RPC response returned by the server.
type serverResponseXML struct {
	// The Object that was returned by the invoked method. This must be null
	// in case there was an error invoking the method.
	// As per spec the member will be omitted if there was an error.
	//Result interface{}

	// An Error object if there was an error invoking the method. It must be
	// null if there was no error.
	// As per spec the member will be omitted if there was no error.
	//Error *ErrorXML
}

type ResultXML struct {
	XMLName   xml.Name     `xml:"response"`
	Response  interface{}  `xml:"result"`

}

type ErrorXML struct {
	XMLName   xml.Name     `xml:"error"`
	Code      int          `xml:"error_code"`
	Message   string       `xml:"error_msg"` /* required */ // A Primitive or Structured value that contains additional information about the error.
	Data      interface{}  `xml:"error_data,omitempty"` /* optional */
}

func (e *ErrorXML) Error() string {
	return e.Message
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
func (c *serverResponseXML) WriteResponse(w http.ResponseWriter, reply interface{}) {
	var answer ResultXML
	answer.Response = reply
	c.writeServerResponse(w, 200, answer)
}

func (c *serverResponseXML) WriteError(w http.ResponseWriter, status int, errIn error) {
	res := ErrorXML{Code: errIn.(*Error).Code, Message: errIn.(*Error).Message, Data: errIn.(*Error).Data}
	c.writeServerResponse(w, status, res)
}

func (c *serverResponseXML) writeServerResponse(w http.ResponseWriter, status int, res interface{}) {
	b, err := xml.MarshalIndent(res,"", " ")
	if err != nil {
		WritePureError(w, 500, err.Error())
	}
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	w.Write(b)
}