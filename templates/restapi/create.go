package {{.PackageName}}

import (
	"fmt"
	"net/http"
	"strings"
)

func (api *{{.Name}}RestAPI) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	var object {{.Name}}

	switch strings.ToLower(r.Method) {
		case "put":
			fallthrough
		case "post":
			err = api.unmarshalBody(r, &object)
			if err != nil {
				code := api.generateErrorCode()
				api.Logger.Printf("[{{.Name}}API:Create] Error: %v, ErrorCode: %s\n", err, code)
				api.writeError(w, r, fmt.Sprintf("Request was malformed, Code: %s", code))
				return
			}
		default:
			api.writeError(w, r, "Must be a POST or PUT request")
			return
	}
	if api.Hooks.PreCreate != nil {
		if err := api.Hooks.PreCreate(r, &object); err != nil {
			if err == StopOperation {
				api.writeSuccessResponse(w, r)
			} else {
				api.writeError(w, r, err.Error())
			}
			return
		}
	}
	err = api.api.Create(object)
	if err != nil {
		code := api.generateErrorCode()
		api.writeError(w, r, fmt.Sprintf("Could not create {{.Name}}, ErrorCode: %s", code))
		api.Logger.Printf("[{{.Name}}API:Create] Error: %v, ErrorCode: %s\n", err, code)
		return
	}
	if api.Hooks.CreateResponse != nil {
		response, err := api.Hooks.CreateResponse(r)
		if err != nil {
			api.writeError(w, r, err.Error())
		} else {
			api.writeSuccessResponse(w, r, response)
		}
		return
	}
	api.writeSuccessResponse(w, r)
}
