package filters

import (
	"encoding/json"
	"encoding/xml"
	"reflect"

	// TODO "github.com/gorilla/websocket"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

var websocketType = reflect.TypeOf((*websocket.Conn)(nil))

func SwaggerInvoker(c *revel.Controller, fc []revel.Filter) {
	// Instantiate the method.
	methodValue := reflect.ValueOf(c.AppController).MethodByName(c.MethodType.Name)
	// Collect the values for the method's arguments.
	var methodArgs []reflect.Value
	for _, arg := range c.MethodType.Args {
		// If they accept a websocket connection, treat that arg specially.
		var boundArg reflect.Value
		if arg.Type == websocketType {
			boundArg = reflect.ValueOf(c.Request.Websocket)
		} else if arg.Type.Kind() == reflect.Struct {
			boundArg = reflect.Indirect(BindStruct(c.Request, arg.Type))
		} else if arg.Type.Kind() == reflect.Ptr && arg.Type.Elem().Kind() == reflect.Struct {
			revel.WARN.Println(arg.Type.Elem())
			boundArg = BindStruct(c.Request, arg.Type.Elem())
		} else {
			boundArg = revel.Bind(c.Params, arg.Name, arg.Type)
		}
		methodArgs = append(methodArgs, boundArg)
	}
	var resultValue reflect.Value
	if methodValue.Type().IsVariadic() {
		resultValue = methodValue.CallSlice(methodArgs)[0]
	} else {
		resultValue = methodValue.Call(methodArgs)[0]
	}
	if resultValue.Kind() == reflect.Interface && !resultValue.IsNil() {
		c.Result = resultValue.Interface().(revel.Result)
	}
}

func BindStruct(req *revel.Request, typ reflect.Type) reflect.Value {
	var (
		err error
		obj = reflect.New(typ)
	)
	switch req.ContentType {
	case "application/json":
		err = json.NewDecoder(req.Body).
			Decode(obj.Interface())
	case "application/xml":
		err = xml.NewDecoder(req.Body).
			Decode(obj.Interface())
	}
	if err != nil {
		revel.ERROR.Println(err)
	}
	return obj
}
