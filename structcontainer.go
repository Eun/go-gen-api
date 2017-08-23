package gogenapi

type structField struct {
	Name string
	Type string
}

type structContainer struct {
	Name   string
	Fields []structField
}
