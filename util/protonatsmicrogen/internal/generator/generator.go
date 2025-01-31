package generator

import (
	"bytes"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type TemplateData struct {
	Package  string
	Services []ServiceData
}

type ServiceData struct {
	Name    string
	Methods []MethodData
}

type MethodData struct {
	Name       string
	InputType  string
	OutputType string
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File) error {
	filename := file.GeneratedFilenamePrefix + "_nats.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("service").Funcs(funcMap).Parse(serviceTmpl)
	if err != nil {
		return err
	}

	data := TemplateData{
		Package:  string(file.GoPackageName),
		Services: make([]ServiceData, 0, len(file.Services)),
	}

	for _, service := range file.Services {
		svcData := ServiceData{
			Name:    service.GoName,
			Methods: make([]MethodData, 0, len(service.Methods)),
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
		return err
	}

	g.P(buf.String())
	return nil
}
