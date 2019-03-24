package vapi

import "sync"

var (
	responsePool sync.Pool
	errorsPool   sync.Pool
)

// AcquireResponse returns an empty Response instance from response pool.
//
// The returned Response instance may be passed to ReleaseResponse when it is
// no longer needed. This allows Response recycling, reduces GC pressure
// and usually improves performance.
func acquireResponse() *ServerResponse {
	v := responsePool.Get()
	if v == nil {
		return &ServerResponse{}
	}
	return v.(*ServerResponse)
}

// ReleaseResponse return resp acquired via AcquireResponse to response pool.
//
// It is forbidden accessing resp and/or its' members after returning
// it to response pool.
func releaseResponse(resp *ServerResponse) {
	resp.Response = nil
	resp.Error = nil
	responsePool.Put(resp)
}

// AcquireResponse returns an empty Response instance from response pool.
//
// The returned Response instance may be passed to ReleaseResponse when it is
// no longer needed. This allows Response recycling, reduces GC pressure
// and usually improves performance.
func acquireError() *Error {
	v := errorsPool.Get()
	if v == nil {
		return &Error{}
	}
	return v.(*Error)
}

// ReleaseResponse return resp acquired via AcquireResponse to response pool.
//
// It is forbidden accessing resp and/or its' members after returning
// it to response pool.
func releaseError(err *Error) {
	err.ErrorCode = 0
	err.ErrorMessage = ""
	err.ErrorHTTPCode = 0
	err.Data = nil
	errorsPool.Put(err)
}
