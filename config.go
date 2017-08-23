package gogenapi

type Config struct {
	Structs             []interface{}
	OutputPath          string
	TemplatePath        string
	SkipGenerateRestAPI bool
}
