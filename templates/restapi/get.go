package {{.Name | ToLower}}

import (
	"fmt"
	"net/http"
	"strings"
)

func (api *RestAPI) Get(w http.ResponseWriter, r *http.Request) {
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

	if api.Hooks.PreGet != nil {
		if err := api.Hooks.PreGet(r, &object); err != nil {
			api.writeError(w, r, err.Error())
			return
		}
	}

	var objects []{{.Name}}
	objects, err = api.api.Get(object)
	
	if err != nil {
		code := api.generateErrorCode()
		api.writeError(w, r, fmt.Sprintf("Could not get {{.Name}}, ErrorCode: %s", code))
		api.Logger.Printf("[{{.Name}}API:Get] Error: %v, ErrorCode: %s\n", err, code)
		return
	}
	if api.Hooks.GetResponse != nil {
		response, err := api.Hooks.GetResponse(r, objects)
		if err != nil {
			api.writeError(w, r, err.Error())
		} else {
			api.writeSuccessResponse(w, r, response)
		}
		return
	}
	api.writeSuccessResponse(w, r, objects)
}
