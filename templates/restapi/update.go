package {{.PackageName}}

import (
	"fmt"
	"net/http"
	"strings"
)

func (api *{{.Name}}RestAPI) Update(w http.ResponseWriter, r *http.Request) {
	var err error
	var update struct {
		Find {{.Name}}
		Update {{.Name}}
	}

	switch strings.ToLower(r.Method) {
		case "put":
			fallthrough
		case "post":
			err = api.unmarshalBody(r, &update)
			if err != nil {
				code := api.generateErrorCode()
				api.Logger.Printf("[{{.Name}}API:Update] Error: %v, ErrorCode: %s\n", err, code)
				api.writeError(w, r, fmt.Sprintf("Request was malformed, Code: %s", code))
				return
			}
		default:
			api.writeError(w, r, "Must be a POST or PUT request")
			return
	}


	if api.Hooks.PreUpdate != nil {
		if err := api.Hooks.PreUpdate(r, &update.Find, &update.Update); err != nil {
			if err == StopOperation {
				api.writeSuccessResponse(w, r)
			} else {
				api.writeError(w, r, err.Error())
			}
			return
		}
	}
	err = api.api.Update(update.Find, update.Update)
	if err != nil {
		code := api.generateErrorCode()
		api.writeError(w, r, fmt.Sprintf("Could not update {{.Name}}, ErrorCode: %s", code))
		api.Logger.Printf("[{{.Name}}API:Update] Error: %v, ErrorCode: %s\n", err, code)
		return
	}
	if api.Hooks.UpdateResponse != nil {
		response, err := api.Hooks.UpdateResponse(r)
		if err != nil {
			api.writeError(w, r, err.Error())
		} else {
			api.writeSuccessResponse(w, r, response)
		}
		return
	}
	api.writeSuccessResponse(w, r)
}
