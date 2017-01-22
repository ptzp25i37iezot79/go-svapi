package vapi

type Error struct {
	Code    int         `json:"code"`    /* required */ // A Number that indicates the error type that occurred.
	Message string      `json:"message"` /* required */ // A String providing a short description of the error.  The message SHOULD be limited to a concise single sentence.
	Data    interface{} `json:"data"`    /* optional */ // A Primitive or Structured value that contains additional information about the error.
}

func (e *Error) Error() string {
	return e.Message
}
