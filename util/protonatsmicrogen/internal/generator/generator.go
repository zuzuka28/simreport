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

	tmpl, err := template.New(".").Funcs(sprig.FuncMap()).Parse(serviceTmpl)
	if err != nil {
		return err
	}

	packageAlias := string(file.GoPackageName)
	packagePath := string(file.GoImportPath)

	data := TemplateData{
		Package:      packageAlias,
		PackageAlias: packageAlias,
		PackagePath:  packagePath,
		Services:     make([]ServiceData, 0, len(file.Services)),
	}

	for _, service := range file.Services {
		svcData := ServiceData{
			Name:    service.GoName,
			Methods: make([]MethodData, 0, len(service.Methods)),
		}

		for _, method := range service.Methods {
			svcData.Methods = append(svcData.Methods, MethodData{
				Name:       method.GoName,
				Reciever:   svcData.Name,
				InputType:  method.Input.GoIdent.GoName,
				OutputType: method.Output.GoIdent.GoName,
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
