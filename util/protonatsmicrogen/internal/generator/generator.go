package generator

import (
	"bytes"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type ServiceData struct {
	Name          string
	LowercaseName string
	Methods       []MethodData
}

type MethodData struct {
	Name       string
	InputType  string
	OutputType string
}

type TemplateData struct {
	Package  string
	Services []ServiceData
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_nats.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	tmpl, err := template.New("service").Parse(serviceTmpl)
	if err != nil {
		gen.Error(err)
		return
	}

	data := TemplateData{
		Package:  string(file.GoPackageName),
		Services: make([]ServiceData, 0, len(file.Services)),
	}

	for _, service := range file.Services {
		svcData := ServiceData{
			Name:          service.GoName,
			LowercaseName: strings.ToLower(service.GoName),
			Methods:       make([]MethodData, 0, len(service.Methods)),
		}

		for _, method := range service.Methods {
			svcData.Methods = append(svcData.Methods, MethodData{
				Name:       method.GoName,
				InputType:  method.Input.GoIdent.String(),
				OutputType: method.Output.GoIdent.String(),
			})
		}

		data.Services = append(data.Services, svcData)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		gen.Error(err)
		return
	}

	g.P(buf.String())
}
