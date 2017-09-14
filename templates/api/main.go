package {{.Name | ToLower}}

import (
	"errors"
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

var StopOperation = errors.New("StopOperation")

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
		tempObject := struct {
		{{- range $i, $e := .Fields}}
			{{if eq $e.Type "string"}}
				{{- $e.Name}} sql.NullString
			{{- else if eq $e.Type "int"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "int16"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "int32"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "int64"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "uint"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "uint16"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "uint32"}}
				{{- $e.Name}} sql.NullInt64
			{{- else if eq $e.Type "uint64"}}
				{{- $e.Name}} sql.NullInt64
			{{end}}
		{{- end}}
		} {}
		if err := rows.Scan({{range $i, $e := .Fields}}{{if $i}}, {{end}}&tempObject.{{$e.Name}}{{end}}); err != nil {
			return nil, err
		}

		var object {{.Name}}
		{{range $i, $e := .Fields}}
		if tempObject.{{$e.Name}}.Valid {
			object.{{$e.Name}} = new({{$e.Type}})
		{{- if eq $e.Type "string"}}
			*object.{{$e.Name}} = tempObject.{{$e.Name}}.String
		{{- else if eq $e.Type "int"}}
			*object.{{$e.Name}} = int(tempObject.{{$e.Name}}.Int64)
		{{- else if eq $e.Type "int16"}}
			*object.{{$e.Name}} = int16(tempObject.{{$e.Name}}.Int64)
		{{- else if eq $e.Type "int32"}}
			*object.{{$e.Name}} = int32(tempObject.{{$e.Name}}.Int64)
		{{- else if eq $e.Type "int64"}}
			*object.{{$e.Name}} = tempObject.{{$e.Name}}.Int64
		{{- else if eq $e.Type "uint"}}
			*object.{{$e.Name}} = uint(tempObject.{{$e.Name}}.Int64)
		{{- else if eq $e.Type "uint16"}}
			*object.{{$e.Name}} = uint16(tempObject.{{$e.Name}}.Int64)
		{{- else if eq $e.Type "uint32"}}
			*object.{{$e.Name}} = uint32(tempObject.{{$e.Name}}.Int64)
		{{- else if eq $e.Type "uint64"}}
			*object.{{$e.Name}} = uint64(tempObject.{{$e.Name}}.Int64)
		{{- end}}
		}
		{{- end}}
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