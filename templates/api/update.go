package {{.Name | ToLower}}

import (
	"fmt"
	"strings"
)

func (api {{.Name}}API) Update(findObject {{.Name}}, updateObject {{.Name}}) error {
	queries, values := api.getSetValues(&updateObject)
	if queries == nil || values == nil || len(queries) <= 0 || len(values) <= 0 {
		return nil
	}

	for i := range queries {
		queries[i] = queries[i] + " = ?"
	}

	objects, err := api.Get(findObject)
	if err != nil {
		return err
	}
	if api.Hooks.PreUpdate != nil {
		if err = api.Hooks.PreUpdate(objects); err != nil {
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
		_, err = api.DB.Exec(fmt.Sprintf("UPDATE {{.Name}} SET %s WHERE %s", strings.Join(queries, ", "), strings.Join(whereQueries, " AND ")), append(values, whereValues...)...)
	}
	if api.Hooks.PostUpdate != nil {
		api.Hooks.PostUpdate(objects)
	}
	return err
}


func (api {{.Name}}API) UpdateWhere(updateObject {{.Name}}, whereQuery string, whereValues ...interface{}) error {
	queries, values := api.getSetValues(&updateObject)
	if queries == nil || values == nil || len(queries) <= 0 || len(values) <= 0 {
		return nil
	}

	objects, err := api.GetWhere(whereQuery, whereValues...)
	if err != nil {
		return err
	}
	if api.Hooks.PreUpdate != nil {
		if err = api.Hooks.PreUpdate(objects); err != nil {
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
		_, err = api.DB.Exec(fmt.Sprintf("UPDATE {{.Name}} SET %s WHERE %s", strings.Join(queries, ", "), strings.Join(whereQueries, " AND ")), append(values, whereValues...)...)
	}
	if api.Hooks.PostUpdate != nil {
		api.Hooks.PostUpdate(objects)
	}
	return err
}
