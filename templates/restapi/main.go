package {{.Name | ToLower}}

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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


func (api *RestAPI) read(r *http.Request) (bytes []byte, err error) {
	bytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return bytes, err
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

func (api *RestAPI) unmarshalUrlValues(values url.Values, result *{{.Name}}) error {
	for key, value := range values {
		if len(value) <= 0 {
			continue
		}
		switch strings.ToLower(key) {
		{{range $i, $e := .Fields -}}
		case "{{$e.Name | ToLower}}":
			{{if eq $e.Type "string" -}}
			result.{{$e.Name}} = &value[0]
		{{else if eq $e.Type "int16" -}}
			val, err := parseInt(value[0], 10, 16)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &int16(val)
		{{else if eq $e.Type "int32" -}}
			val, err := parseInt(value[0], 10, 32)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &int32(val)
		{{else if eq $e.Type "int64" -}}
			val, err := parseInt(value[0], 10, 64)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &val
		{{else if eq $e.Type "int" -}}
			val, err := parseInt(value[0], 10, 32)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &int(val)
		{{else if eq $e.Type "uint16" -}}
			val, err := parseInt(value[0], 10, 16)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &uint16(val)
		{{else if eq $e.Type "uint32" -}}
			val, err := parseInt(value[0], 10, 32)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &uint32(val)
		{{else if eq $e.Type "uint64" -}}
			val, err := parseInt(value[0], 10, 64)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &uint64(val)
		{{else if eq $e.Type "uint" -}}
			val, err := parseInt(value[0], 10, 32)
			if err != nil {
				return err
			}
			result.{{$e.Name}} = &uint(val)
		{{end -}}
		{{end -}}
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