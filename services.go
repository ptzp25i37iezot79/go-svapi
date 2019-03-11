package vapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

var (
	// Precompute the reflect.Type of error and fasthttp.RequestCtx
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfArgs    = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	typeOfReply   = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	typeOfRequest = reflect.TypeOf((*fasthttp.RequestCtx)(nil)).Elem()
)

// VAPI - main structure
type VAPI struct {
	mutex    sync.RWMutex
	services map[string]bool
	methods  map[string]*serviceMethod
}

// serviceMethod - sub struct
type serviceMethod struct {
	//name      string         // name of service
	rcvr      reflect.Value  // receiver of methods for the service
	rcvrType  reflect.Type   // type of the receiver
	method    reflect.Method // receiver method
	argsType  reflect.Type   // type of the request argument
	replyType reflect.Type   // type of the response argument
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
//    - The method has three arguments: *fasthttp.RequestCtx, *args, *reply.
//    - All three arguments are pointers.
//    - The second and third arguments are exported or local.
//    - The method has return type error.
//
// All other methods are ignored.
func (as *VAPI) RegisterService(receiver interface{}, name string) error {
	return as.register(receiver, name)
}

// register adds a new service using reflection to extract its methods.
func (as *VAPI) register(rcvr interface{}, serviceName string) error {

	rcvrValue := reflect.ValueOf(rcvr)
	rcvrType := reflect.TypeOf(rcvr)

	if serviceName == "" {
		serviceName = reflect.Indirect(rcvrValue).Type().Name()

		if !isExported(serviceName) {
			return fmt.Errorf("vapi: type %q is not exported", serviceName)
		}
	}

	if serviceName == "" {
		return fmt.Errorf("vapi: no service name for type %q", rcvrType.String())
	}

	as.mutex.RLock()
	defer as.mutex.RUnlock()

	if _, ok := as.services[serviceName]; ok {
		return fmt.Errorf("vapi: service already defined: %q", serviceName)
	}

	as.services[serviceName] = true

	addedMethodCounter := 0

	// Setup methods.
	for i := 0; i < rcvrType.NumMethod(); i++ {

		method := rcvrType.Method(i)
		mtype := method.Type

		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}

		// Method needs four ins: receiver, *fasthttp.RequestCtx, *args, *reply.
		if mtype.NumIn() != 4 {
			continue
		}

		// First argument must be a pointer and must be fasthttp.RequestCtx.
		reqType := mtype.In(1)
		if reqType.Kind() != reflect.Ptr || reqType.Elem() != typeOfRequest {
			continue
		}

		// Second argument is Args must be a pointer, must be exported and must implement Unmarshaller interface.
		args := mtype.In(2)
		if args.Kind() != reflect.Ptr || !isExportedOrBuiltin(args) || !args.Implements(typeOfArgs) {
			continue
		}

		// Third argument must be a pointer, must be exported and must implement Marshaller interface.
		reply := mtype.In(3)
		if reply.Kind() != reflect.Ptr || !isExportedOrBuiltin(reply) || !reply.Implements(typeOfReply) {
			continue
		}

		// Method needs one out: error.
		if mtype.NumOut() != 1 {
			continue
		}

		// Method out should be an error type.
		if returnType := mtype.Out(0); returnType.Elem() != typeOfError {
			continue
		}

		as.methods[fmt.Sprintf("%s.%s", serviceName, method.Name)] = &serviceMethod{
			rcvr:      rcvrValue,
			rcvrType:  rcvrType,
			method:    method,
			argsType:  args.Elem(),
			replyType: reply.Elem(),
		}

		addedMethodCounter++
	}

	if addedMethodCounter == 0 {
		return fmt.Errorf("vapi: %q has no exported methods of suitable type", serviceName)
	}

	return nil
}

// get returns a registered service method by given name.
//
// The method name uses a dotted notation as in "Service.Method".
func (as *VAPI) get(serviceWithMethod string) (*serviceMethod, error) {

	parts := strings.Split(serviceWithMethod, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("vapi: service/method request ill-formed: %q", serviceWithMethod)
	}

	if _, ok := as.services[parts[0]]; !ok {
		return nil, fmt.Errorf("vapi: service not found: %q", parts[0])
	}

	as.mutex.Lock()
	serviceMethod, okMethod := as.methods[serviceWithMethod]
	as.mutex.Unlock()

	if !okMethod {
		return nil, fmt.Errorf("vapi: can't find method %q", parts[1])
	}

	return serviceMethod, nil
}

// GetServiceMap returns an json api schema
// todo realize this function
func (as *VAPI) GetServiceMap() (map[string]*serviceMethod, error) {
	methods := as.methods
	return methods, nil
}

// CallAPI call api method and process it.
// Modifying body after this function not recommended
func (as *VAPI) CallAPI(ctx *fasthttp.RequestCtx, method string) {

	methodSpec, err := as.get(method)

	srvResponse := responsePool.Get().(*ServerResponse)
	defer responsePool.Put(srvResponse)

	if err != nil {

		errAPI := errorsPool.Get().(*Error)
		defer errorsPool.Put(errAPI)
		errAPI.ErrorHTTPCode = fasthttp.StatusNotFound
		errAPI.ErrorCode = 0
		errAPI.ErrorMessage = err.Error()

		srvResponse.Error = errAPI
		WriteResponse(ctx, errAPI.ErrorHTTPCode, *srvResponse)
		return
	}

	// Decode the args.
	args := reflect.New(methodSpec.argsType)
	err = args.Interface().(json.Unmarshaler).UnmarshalJSON(ctx.Request.Body())
	if err != nil {

		errAPI := errorsPool.Get().(*Error)
		defer errorsPool.Put(errAPI)
		errAPI.ErrorHTTPCode = fasthttp.StatusInternalServerError
		errAPI.ErrorCode = 0
		errAPI.ErrorMessage = err.Error()

		srvResponse.Error = errAPI
		WriteResponse(ctx, errAPI.ErrorHTTPCode, *srvResponse)
		return
	}

	// Call the service method.
	reply := reflect.New(methodSpec.replyType)
	errValue := methodSpec.method.Func.Call([]reflect.Value{
		methodSpec.rcvr,
		reflect.ValueOf(ctx),
		args,
		reply,
	})

	errInter := errValue[0].Interface()
	if errInter != nil {
		errResult := errorsPool.Get().(*Error)
		defer errorsPool.Put(errResult)

		srvResponse.Error = errInter.(*Error)
		WriteResponse(ctx, srvResponse.Error.ErrorHTTPCode, *srvResponse)
		return
	}

	repBytes, err := reply.Interface().(json.Marshaler).MarshalJSON()
	if err != nil {

		errAPI := errorsPool.Get().(*Error)
		defer errorsPool.Put(errAPI)
		errAPI.ErrorHTTPCode = fasthttp.StatusInternalServerError
		errAPI.ErrorCode = 0
		errAPI.ErrorMessage = err.Error()

		srvResponse.Error = errAPI
		WriteResponse(ctx, errAPI.ErrorHTTPCode, *srvResponse)
		return
	}

	srvResponse.Response = repBytes
	WriteResponse(ctx, fasthttp.StatusOK, *srvResponse)
	return
}

// NewServer returns a new RPC server.
func NewServer() *VAPI {
	return &VAPI{
		services: make(map[string]bool),
		methods:  make(map[string]*serviceMethod),
	}
}
