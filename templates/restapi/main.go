package {{.Name | ToLower}}

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"go/ast"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	

	"github.com/gorilla/mux"
)

type RestAPI struct {
	router *mux.Router
	api    *{{.Name}}API
	Logger *log.Logger
	Hooks  *RestHooks
}

type RestHooks struct {
    PreCreate      func(r *http.Request, object *{{.Name}}) error
	PreUpdate      func(r *http.Request, findObject *{{.Name}}, updateObject *{{.Name}}) error
	PreDelete      func(r *http.Request, object *{{.Name}}) error
	PreGet         func(r *http.Request, object *{{.Name}}) error

	CreateResponse func(r *http.Request) (interface{}, error)
	UpdateResponse func(r *http.Request) (interface{}, error)
	DeleteResponse func(r *http.Request) (interface{}, error)
	GetResponse    func(r *http.Request, objects []{{.Name}}) (interface{}, error)
}

func NewRestAPI(router *mux.Router, api *{{.Name}}API) *RestAPI {
	a := RestAPI {
		router: router,
		api: api,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
		Hooks: new(RestHooks),
	}
	a.router.HandleFunc("/create", a.Create)
	a.router.HandleFunc("/create/", a.Create)
	a.router.HandleFunc("/delete", a.Delete)
	a.router.HandleFunc("/delete/", a.Delete)
	a.router.HandleFunc("/get", a.Get)
	a.router.HandleFunc("/get/", a.Get)
	a.router.HandleFunc("/update", a.Update)
	a.router.HandleFunc("/update/", a.Update)
	return &a
}

func (api *RestAPI) writeError(w http.ResponseWriter, r *http.Request, err string) {
	switch strings.ToLower(r.Header.Get("Content-Type")) {
	case "application/xml":
		fallthrough
	case "text/xml":
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(400)
		xml.NewEncoder(w).Encode(&struct{ Error string }{err})
	case "application/json":
		fallthrough
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&struct{ Error string }{err})
	}
}

func (api *RestAPI) writeSuccessResponse(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	switch strings.ToLower(r.Header.Get("Content-Type")) {
	case "application/xml":
		fallthrough
	case "text/xml":
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(200)
		if len(v) > 0 {
			xml.NewEncoder(w).Encode(v[0])
		}
	case "application/json":
		fallthrough
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		if len(v) > 0 {
			json.NewEncoder(w).Encode(v[0])
		}
	}
}

func (api *RestAPI) generateErrorCode() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	return hex.EncodeToString(uuid)
}

type contextKey int

const bodyKey contextKey = iota

func (api *RestAPI) read(r *http.Request) (bytes []byte, err error) {
	body := r.Context().Value(bodyKey)
	if body == nil {
		bytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		*r = *r.WithContext(context.WithValue(r.Context(), bodyKey, bytes))
		return bytes, err
	}
	return body.([]byte), nil
}

func (api *RestAPI) unmarshal(r *http.Request, result interface{}) error {
	var unmarshal func(data []byte, v interface{}) error
	switch strings.ToLower(r.Header.Get("Content-Type")) {
	case "application/xml":
		fallthrough
	case "text/xml":
		unmarshal = xml.Unmarshal
	case "application/json":
		unmarshal = json.Unmarshal
	default:
		return errors.New("Unknown Content-Type")
	}
	bytes, err := api.read(r)
	if err != nil {
		return err
	}

	err = unmarshal(bytes, result)
	if err != nil {
		return err
	}
	return nil
}

func parseInt(s string, base int, bitSize int) (i int64, err error) {
	return strconv.ParseInt(s, base, bitSize)
}

func (api *RestAPI) unmarshalUrlValues(values url.Values, result interface{}) error {
	reflectValue := reflect.ValueOf(result)
	
	if reflectValue.Kind() != reflect.Ptr || reflectValue.IsNil() {
		return fmt.Errorf("Invalid type %s", reflect.TypeOf(reflectValue).String())
	}

	reflectValue = reflectValue.Elem()

	// ptr in ptr is not supported
	if reflectValue.Kind() == reflect.Ptr {
		return fmt.Errorf("Invalid type %s", reflect.TypeOf(reflectValue).String())
	}

	reflectType := reflectValue.Type()

	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			var ptrField *reflect.Value
			field := reflectValue.Field(i)
			fieldType := fieldStruct.Type
			if fieldType.Kind() == reflect.Ptr {
				ptrField = &field
				fieldType = fieldType.Elem()
			}

			// ptr in ptr is not supported
			if fieldType.Kind() == reflect.Ptr {
				return fmt.Errorf("Invalid type %s in %s", fieldType.String(), fieldStruct.Name)
			}

			for key, value := range values {
				if strings.EqualFold(key, fieldStruct.Name) {
					switch fieldType.Kind() {
					case reflect.String:
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetString(value[0])
					case reflect.Int16:
						val, err := parseInt(value[0], 10, 16)
						if err != nil {
							return err
						}
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetInt(val)
					case reflect.Int:
						fallthrough
					case reflect.Int32:
						val, err := parseInt(value[0], 10, 32)
						if err != nil {
							return err
						}
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetInt(val)
					case reflect.Int64:
						val, err := parseInt(value[0], 10, 64)
						if err != nil {
							return err
						}
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetInt(val)
					case reflect.Uint16:
						val, err := parseInt(value[0], 10, 16)
						if err != nil {
							return err
						}
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetUint(uint64(val))
					case reflect.Uint:
						fallthrough
					case reflect.Uint32:
						val, err := parseInt(value[0], 10, 32)
						if err != nil {
							return err
						}
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetUint(uint64(val))
					case reflect.Uint64:
						val, err := parseInt(value[0], 10, 64)
						if err != nil {
							return err
						}
						if ptrField != nil {
							ptrField.Set(reflect.New(fieldType))
							field = ptrField.Elem()
						}
						field.SetUint(uint64(val))
					}
				}
			}
		}
	}
	
	return nil
}

func (api *RestAPI) customHandler(f func(r *http.Request, object *{{.Name}}) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var object {{.Name}}
	
		switch strings.ToLower(r.Method) {
			case "post":
				err = api.unmarshal(r, &object)
				if err != nil {
					api.writeError(w, r, "Request was malformed")
					return
				}
			case "get":
				err = api.unmarshalUrlValues(r.URL.Query(), &object)
				if err != nil {
					api.writeError(w, r, "Request was malformed")
					return
				}
			default:
				api.writeError(w, r, "Must be a POST or GET request")
				return
		}
		var result interface{}
		result, err = f(r, &object)
		if err != nil {
			api.writeError(w, r, err.Error())
			return
		}
		api.writeSuccessResponse(w, r, result)
	}
}

func (api *RestAPI) HandleFunc(path string, f func(r *http.Request, object *{{.Name}}) (interface{}, error)) {
	api.router.HandleFunc(path, api.customHandler(f))
}

func (api *RestAPI) GetCustomFields(r *http.Request, object interface{}) error {
	switch strings.ToLower(r.Method) {
	case "post":
		err := api.unmarshal(r, object)
		if err != nil {
			return errors.New("Request was malformed")
		}
	case "get":
		err := api.unmarshalUrlValues(r.URL.Query(), object)
		if err != nil {
			return errors.New("Request was malformed")
		}
	default:
		return errors.New("Must be a POST or GET request")
	}
	return nil
}