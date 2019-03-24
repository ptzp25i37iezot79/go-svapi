package svapi

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

var (
	// Precompute the reflect.Type of error and fasthttp.RequestCtx
	typeOfRequest = reflect.TypeOf((*fasthttp.RequestCtx)(nil)).Elem()
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
)

// SVAPI - main structure
type SVAPI struct {
	mutex         sync.RWMutex
	services      map[string]bool
	errorCallback ErrorHandlerFunction
	methods       map[string]*serviceMethod
	serviceMap    map[string][]string
}

// serviceMethod - sub struct
type serviceMethod struct {
	rcvr     reflect.Value  // receiver of methods for the service
	rcvrType reflect.Type   // type of the receiver
	method   reflect.Method // receiver method
}

// RegisterService adds a new service to the api server.
//
// The name parameter is optional: if empty it will be inferred from
// the receiver type name.
//
// Methods from the receiver will be extracted if these rules are satisfied:
//
//    - The receiver is exported (begins with an upper case letter) or local
//      (defined in the package registering the service).
//    - The method name is exported.
//    - The method has one argument: *fasthttp.RequestCtx.
//    - All arguments are pointers.
//    - The method has return type error that will be processed by error handler method.
//
// All other methods are ignored.
func (as *SVAPI) RegisterService(receiver interface{}, name string) error {
	return as.register(receiver, name)
}

// register adds a new service using reflection to extract its methods.
func (as *SVAPI) register(rcvr interface{}, serviceName string) error {

	rcvrValue := reflect.ValueOf(rcvr)
	rcvrType := reflect.TypeOf(rcvr)

	if serviceName == "" {
		serviceName = reflect.Indirect(rcvrValue).Type().Name()

		if !isExported(serviceName) {
			return fmt.Errorf("svapi: type %q is not exported", serviceName)
		}
	}

	if serviceName == "" {
		return fmt.Errorf("svapi: no service name for type %q", rcvrType.String())
	}

	as.mutex.RLock()
	defer as.mutex.RUnlock()

	if _, ok := as.services[serviceName]; ok {
		return fmt.Errorf("svapi: service already defined: %q", serviceName)
	}

	as.services[serviceName] = true

	var tmpMethodList []string

	// Setup methods.
	for i := 0; i < rcvrType.NumMethod(); i++ {

		method := rcvrType.Method(i)
		mtype := method.Type

		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}

		// Method needs four ins: receiver, *fasthttp.RequestCtx.
		if mtype.NumIn() != 2 {
			continue
		}

		// First argument must be a pointer and must be fasthttp.RequestCtx.
		reqType := mtype.In(1)
		if reqType.Kind() != reflect.Ptr || reqType.Elem() != typeOfRequest {
			continue
		}

		// Method needs one out: error.
		if mtype.NumOut() != 1 {
			continue
		}

		// Method out should be an error type.
		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}

		as.methods[fmt.Sprintf("%s.%s", serviceName, method.Name)] = &serviceMethod{
			rcvr:     rcvrValue,
			rcvrType: rcvrType,
			method:   method,
		}

		tmpMethodList = append(tmpMethodList, method.Name)

	}

	if len(tmpMethodList) == 0 {
		return fmt.Errorf("svapi: %q has no exported methods of suitable type", serviceName)
	}

	as.serviceMap[serviceName] = tmpMethodList

	return nil
}

// get returns a registered service method by given name.
//
// The method name uses a dotted notation as in "Service.Method".
func (as *SVAPI) get(serviceWithMethod string) (*serviceMethod, error) {

	parts := strings.Split(serviceWithMethod, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("svapi: service/method request ill-formed: %q", serviceWithMethod)
	}

	if _, ok := as.services[parts[0]]; !ok {
		return nil, fmt.Errorf("svapi: service not found: %q", parts[0])
	}

	as.mutex.Lock()
	serviceMethod, okMethod := as.methods[serviceWithMethod]
	as.mutex.Unlock()

	if !okMethod {
		return nil, fmt.Errorf("svapi: can't find method %q", parts[1])
	}

	return serviceMethod, nil
}

// GetServiceMap returns an api methods list
func (as *SVAPI) GetServiceMap() map[string][]string {
	return as.serviceMap
}

// CallAPI call api method and process it.
// Modifying body after this function not recommended
func (as *SVAPI) CallAPI(ctx *fasthttp.RequestCtx, method string) {

	methodSpec, err := as.get(method)

	if err != nil {
		as.errorCallback(ctx, err)
		return
	}

	// Call the service method.
	errValue := methodSpec.method.Func.Call([]reflect.Value{methodSpec.rcvr, reflect.ValueOf(ctx)})

	errInter := errValue[0].Interface()
	if errInter != nil {
		as.errorCallback(ctx, errInter.(error))
		return
	}

}

// NewServer returns a new RPC server.
func NewServer() *SVAPI {
	return &SVAPI{
		services:      make(map[string]bool),
		methods:       make(map[string]*serviceMethod),
		serviceMap:    make(map[string][]string),
		errorCallback: defaultErrorHandler,
	}
}

// SetErrorHandlerFunction allows to set custom error processing for api functions
func (as *SVAPI) SetErrorHandlerFunction(errHndl ErrorHandlerFunction) {
	as.errorCallback = errHndl
}
