package {{.Name | ToLower}}

import (
	"fmt"
	"net/http"
	"strings"
)

func (api *RestAPI) Delete(w http.ResponseWriter, r *http.Request) {
	var err error
	var object {{.Name}}

	switch strings.ToLower(r.Method) {
		case "delete":
			fallthrough
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
			api.writeError(w, r, "Must be a POST, DELETE or GET request")
			return
	}
	if api.Hooks.PreDelete != nil {
		if err := api.Hooks.PreDelete(r, &object); err != nil {
			api.writeError(w, r, err.Error())
		    return
		}
	}
	err = api.api.Delete(object)
	if err != nil {
		code := api.generateErrorCode()
		api.writeError(w, r, fmt.Sprintf("Could not get {{.Name}}, ErrorCode: %s", code))
        api.Logger.Printf("[{{.Name}}API:Get] Error: %v, ErrorCode: %s\n", err, code)
		return
	}
	if api.Hooks.DeleteResponse != nil {
		response, err := api.Hooks.DeleteResponse(r)
		if err != nil {
			api.writeError(w, r, err.Error())
		} else {
			api.writeSuccessResponse(w, r, response)
		}
		return
	}
	api.writeSuccessResponse(w, r)
}
