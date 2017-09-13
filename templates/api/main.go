package {{.Name | ToLower}}

import (
	"database/sql"
)

type {{.Name}} struct {
	{{- range $i, $e := .Fields}}
	{{$e.Name}} *{{$e.Type -}}
	{{end}}
}

type {{.Name}}API struct {
	DB      *sql.DB
	Hooks   *Hooks
}

type Hooks struct {
	PreCreate  func(object *{{.Name}}) error
	PreUpdate  func(objects []{{.Name}}) error
	PreDelete  func(objects []{{.Name}}) error
	PostCreate func(object *{{.Name}})
	PostUpdate func(objects []{{.Name}})
	PostDelete func(objects []{{.Name}})
}

func New(db *sql.DB) *{{.Name}}API {
	return &{{.Name}}API {
		DB: db,
		Hooks: new(Hooks),
	}
}


func (api {{.Name}}API) scanRows(rows *sql.Rows) ([]{{.Name}}, error){
	var objects []{{.Name}}
	for rows.Next() {
		object := {{.Name}}{
			{{range $i, $e := .Fields}}
			{{$e.Name}}: new({{$e.Type}}),
			{{- end}}
		}
		if err := rows.Scan({{range $i, $e := .Fields}}{{if $i}}, {{end}}object.{{$e.Name}}{{end}}); err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return objects, nil
}

func (api {{.Name}}API) getSetValues(object *{{.Name}}) ([]string, []interface{}) {
	var queries []string
	var queryValues []interface{}
	{{range $i, $e := .Fields -}}
	if object.{{$e.Name}} != nil {
		queries = append(queries, "{{$e.Name}}")
		queryValues = append(queryValues, *object.{{$e.Name}})
	}
	{{end}}
	if len(queries) <= 0 {
		return nil, nil
	}
	return queries, queryValues
}

{{- range $i, $e := .Fields}}
func (object *{{$.Name}}) Set{{$e.Name}}(value {{$e.Type}}) *{{$.Name}} {
	object.{{$e.Name}} = &value
	return object
}
{{end}}