package {{.Name | ToLower}}

import (
	"fmt"
	"strings"
)

func (api {{.Name}}API) Create(newObject {{.Name}}) error {
	if api.Hooks.PreCreate != nil {
		if err := api.Hooks.PreCreate(&newObject); err != nil {
			if err == StopOperation {
				return nil
			}
			return err
		}
	}
	setQueries, setValues := api.getSetValues(&newObject)
	if setQueries == nil || setValues == nil || len(setQueries) <= 0 || len(setValues) <= 0 {
		return nil
	}
	_, err := api.DB.Exec(fmt.Sprintf("INSERT INTO {{.Name}} (%s) VALUES (%s)", strings.Join(setQueries, ", "), strings.Repeat(", ?", len(setQueries))[2:]), setValues...)
	if err != nil {
		if api.Hooks.PostCreate != nil {
			api.Hooks.PostCreate(&newObject)
		}
	}
	return err
}
