package {{.PackageName}}

import (
	"strings"
	"database/sql"
)

func (api {{.Name}}API) Get(findObject {{.Name}}) ([]{{.Name}}, error) {
	queries, values := api.getSetValues(&findObject)
	if queries == nil || values == nil {
		return api.GetAll()
	}
	for i := range queries {
		queries[i] = queries[i] + " = ?"
	}
	return api.GetWhere(strings.Join(queries, " AND "), values...)
}

func (api {{.Name}}API) GetWhere(whereQuery string, whereValues ...interface{}) ([]{{.Name}}, error) {
	var rows *sql.Rows
	var err error
	if len(whereQuery) > 0 {
		rows, err = api.DB.Query("SELECT {{range $i, $e := .Fields}}{{if $i}}, {{end}}{{$e.Name}}{{end}} FROM {{.Name}} WHERE "+whereQuery, whereValues...)
	} else {
		rows, err = api.DB.Query("SELECT {{range $i, $e := .Fields}}{{if $i}}, {{end}}{{$e.Name}}{{end}} FROM {{.Name}}")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return api.scanRows(rows)
}

func (api {{.Name}}API) GetFirst(findObject {{.Name}}) (*{{.Name}}, error) {
	result, err := api.Get(findObject)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, nil
}


func (api {{.Name}}API) GetFirstWhere(whereQuery string, whereValues ...interface{}) (*{{.Name}}, error) {
	result, err := api.GetWhere(whereQuery, whereValues...)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, nil
}

func (api {{.Name}}API) GetAll() ([]{{.Name}}, error) {
	return api.GetWhere("")
}

func (api {{.Name}}API) Exists(findObject {{.Name}}) (bool, error) {
	object, err := api.GetFirst(findObject)
	if err != nil {
		return false, err
	}
	if object == nil {
		return false, nil
	}
	return true, nil
}

func (api {{.Name}}API) ExistsWhere(whereQuery string, whereValues ...interface{}) (bool, error) {
	result, err := api.GetWhere(whereQuery, whereValues...)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, nil
	}
	return true, nil
}