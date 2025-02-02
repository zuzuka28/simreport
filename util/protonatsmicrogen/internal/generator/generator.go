package generator

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"google.golang.org/protobuf/compiler/protogen"
)

type TemplateData struct {
	Package      string
	PackageAlias string
	PackagePath  string
	Services     []ServiceData
}

type ServiceData struct {
	Name    string
	Methods []MethodData
}

type MethodData struct {
	Name       string
	Reciever   string
	InputType  string
	OutputType string
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File) error {
	filename := file.GeneratedFilenamePrefix + "_nats.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	funcMap := template.FuncMap{}

	for k, v := range sprig.FuncMap() {
		funcMap[k] = v
	}

	tmpl, err := template.New(".").Funcs(funcMap).Parse(serviceTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, file); err != nil {
		return err
	}

	g.P(buf.String())
	return nil
}
