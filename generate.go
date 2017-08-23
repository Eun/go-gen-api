package gogenapi

import (
	"errors"
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

type generator struct {
	Config        *Config
	apiTemplates  *template.Template
	restTemplates *template.Template
}

func Generate(config *Config) error {
	if config.Structs == nil || len(config.Structs) <= 0 {
		return errors.New("Structs is empty")
	}
	if len(config.OutputPath) <= 0 {
		return errors.New("Output is empty")
	}

	generator := &generator{
		Config: config,
	}

	if len(generator.Config.TemplatePath) == 0 {
		gopath := os.Getenv("GOPATH")
		for _, p := range filepath.SplitList(gopath) {
			if !strings.HasPrefix(p, "~") && filepath.IsAbs(p) {
				generator.Config.TemplatePath = filepath.Join(p, "src/github.com/Eun/go-gen-api/templates")
				if fi, err := os.Stat(generator.Config.TemplatePath); err == nil && fi.IsDir() {
					break
				}
				generator.Config.TemplatePath = ""
			}
		}
		if len(generator.Config.TemplatePath) == 0 {
			return fmt.Errorf("Unable to find TemplatePath, make sure GOPATH is set accordingly")
		}
	} else {
		if fi, err := os.Stat(generator.Config.TemplatePath); err != nil || !fi.IsDir() {
			return fmt.Errorf("TemplatePath (%s) is invalid: %v", generator.Config.TemplatePath, err)
		}
	}

	err := generator.PrepareOutputDir()
	if err != nil {
		return fmt.Errorf("Unable to prepare output directory '%s': %v", config.OutputPath, err)
	}

	generator.apiTemplates, err = generator.ParseTemplates(filepath.Join(generator.Config.TemplatePath, "api"))
	if err != nil {
		return err
	}

	if generator.Config.SkipGenerateRestAPI == false {
		generator.restTemplates, err = generator.ParseTemplates(filepath.Join(generator.Config.TemplatePath, "restapi"))
		if err != nil {
			return err
		}
	}

	var structContainers []*structContainer

	for _, strct := range config.Structs {
		c, err := generator.GenerateStruct(strct)
		if err != nil {
			return err
		}
		if c != nil {
			structContainers = append(structContainers, c)
		}
	}

	for _, stct := range structContainers {

		if generator.apiTemplates != nil {
			for _, t := range generator.apiTemplates.Templates() {
				generator.GenerateFile(stct, filepath.Join(generator.Config.OutputPath, strings.ToLower(stct.Name), t.Name()), t, t.Name())
			}
		}

		if generator.restTemplates != nil {
			for _, t := range generator.restTemplates.Templates() {
				generator.GenerateFile(stct, filepath.Join(generator.Config.OutputPath, strings.ToLower(stct.Name), "rest_"+t.Name()), t, t.Name())
			}
		}
	}

	return nil
}

func (generator *generator) PrepareOutputDir() error {
	fi, err := os.Stat(generator.Config.OutputPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(generator.Config.OutputPath, 0600)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unable to access '%s'", generator.Config.OutputPath)
		}
	} else if fi.IsDir() == false {
		return fmt.Errorf("%s is not a directory", generator.Config.OutputPath)
	}

	fd, err := os.Open(generator.Config.OutputPath)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = nil
	for {
		names, err1 := fd.Readdirnames(100)
		for _, name := range names {
			err1 := os.RemoveAll(filepath.Join(generator.Config.OutputPath, name))
			if err == nil {
				err = err1
			}
		}
		if err1 == io.EOF {
			break
		}
		if err == nil {
			err = err1
		}
		if len(names) == 0 {
			break
		}
	}
	return err
}

func (generator *generator) GenerateStruct(strct interface{}) (*structContainer, error) {
	if strct == nil {
		log.Printf("Warning: '%v' is null => ignored", strct)
		return nil, nil
	}

	reflectType := reflect.ValueOf(strct).Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Scope value need to be a struct
	if reflectType.Kind() != reflect.Struct {
		log.Printf("Warning: '%v' is not a struct => ignored", strct)
		return nil, nil
	}

	if len(reflectType.Name()) <= 0 {
		log.Printf("Warning: '%v' has no name => ignored", strct)
		return nil, nil
	}

	structContainer := structContainer{Name: reflectType.Name()}
	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			fieldType := fieldStruct.Type
			for fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			structContainer.Fields = append(structContainer.Fields, structField{
				Name: fieldStruct.Name,
				Type: fieldType.String(),
			})
		}
	}

	return &structContainer, nil
}

func (generator *generator) GenerateFile(structContainer *structContainer, filePath string, template *template.Template, templateName string) error {
	if err := os.MkdirAll(filePath[:len(filePath)-len(filepath.Base(filePath))], 0600); err != nil {
		return err
	}

	outputFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("Could not open %s: %v", filePath, err)
	}
	defer outputFile.Close()

	if err := template.ExecuteTemplate(outputFile, templateName, structContainer); err != nil {
		return err
	}
	return nil
}

func (generator *generator) ParseTemplates(dir string) (*template.Template, error) {
	return template.New("").
		Funcs(template.FuncMap{
			"ToLower": strings.ToLower,
			"ToUpper": strings.ToUpper,
		}).
		ParseGlob(filepath.Join(dir, "*.go"))
}
