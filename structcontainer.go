package gogenapi

type structField struct {
	Name string
	Type string
}

type structContainer struct {
	PackageName string
	Name        string
	Fields      []structField
}
