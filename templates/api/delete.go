package {{.Name | ToLower}}

import (
	"strings"
)

func (api {{.Name}}API) Delete(deleteObject {{.Name}}) error {
	objects, err := api.Get(deleteObject)
	if err != nil {
		return err
	}
	if api.Hooks.PreDelete != nil {
		if err = api.Hooks.PreDelete(objects); err != nil {
			if err == StopOperation {
				return nil
			}
			return err
		}
	}

	for _, obj := range objects {
		whereQueries, whereValues := api.getSetValues(&obj)
		if whereQueries == nil || whereValues == nil || len(whereQueries) <= 0 || len(whereValues) <= 0 {
			continue
		}
		for i := range whereQueries {
			whereQueries[i] = whereQueries[i] + " = ?"
		}
		_, err = api.DB.Exec("DELETE FROM {{.Name}} WHERE "+ strings.Join(whereQueries, " AND "), whereValues...)
	}
	if api.Hooks.PostDelete != nil {
		api.Hooks.PostDelete(objects)
	}
	return err
}


func (api {{.Name}}API) DeleteWhere(deleteQuery string, whereValues ...interface{}) error {
	objects, err := api.GetWhere(deleteQuery, whereValues...)
	if err != nil {
		return err
	}
	if api.Hooks.PreDelete != nil {
		if err = api.Hooks.PreDelete(objects); err != nil {
			return err
		}
	}

	for _, obj := range objects {
		whereQueries, whereValues := api.getSetValues(&obj)
		if whereQueries == nil || whereValues == nil || len(whereQueries) <= 0 || len(whereValues) <= 0 {
			continue
		}
		for i := range whereQueries {
			whereQueries[i] = whereQueries[i] + " = ?"
		}
		_, err = api.DB.Exec("DELETE FROM {{.Name}} WHERE "+ strings.Join(whereQueries, " AND "), whereValues...)
	}
	if api.Hooks.PostDelete != nil {
		api.Hooks.PostDelete(objects)
	}
	return err
}


func (api {{.Name}}API) DeleteAll() error {
	return api.DeleteWhere("")
}
