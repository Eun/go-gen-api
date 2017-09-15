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
			err = api.unmarshalBody(r, &object)
		case "get":
			err = api.unmarshalUrlValues(r.URL.Query(), &object)
		default:
			api.writeError(w, r, "Must be a POST or GET request")
			return
	}

	if err != nil {
		code := api.generateErrorCode()
		api.Logger.Printf("[{{.Name}}API:Update] Error: %v, ErrorCode: %s\n", err, code)
		api.writeError(w, r, fmt.Sprintf("Request was malformed, Code: %s", code))
		return
	}

	if api.Hooks.PreGet != nil {
		if err := api.Hooks.PreGet(r, &object); err != nil {
			if err == StopOperation {
				api.writeSuccessResponse(w, r)
			} else {
				api.writeError(w, r, err.Error())
			}
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
