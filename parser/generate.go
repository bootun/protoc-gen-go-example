package parser

import (
	"bytes"
	"fmt"
	"text/template"

	example_tmpl "github.com/bootun/protoc-gen-go-example/template"

	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) error {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + ".example.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	return NewFile(g, file).Generate()
}

type File struct {
	gen *protogen.GeneratedFile
	FileDescription
}

// FileDescription 描述了一个解析过后的proto文件的信息
// 为我们后边的代码生成做准备
type FileDescription struct {
	// PackageName 代表我们生成后的example.pb.go文件的包名
	PackageName string

	// Services 代表我们生成后的example.pb.go文件中的所有服务
	// 我们在proto文件中写的每个server都会转化为一个 Service 实体
	Services []*Service
}

type Service struct {
	// Service 的名称
	Name string

	// Service 里具有哪些方法
	Methods []*Method
}

type Method struct {
	// 方法名称
	Name string
	// 请求类型
	RequestType string
	// 响应类型
	ResponseType string
}

func NewFile(gen *protogen.GeneratedFile, protoFile *protogen.File) *File {
	f := &File{
		gen: gen,
	}

	f.PackageName = string(protoFile.GoPackageName)

	for _, s := range protoFile.Services {
		f.ParseService(s)
	}

	return f
}

func (f *File) ParseService(protoSvc *protogen.Service) {
	s := &Service{
		Name:    protoSvc.GoName,
		Methods: make([]*Method, 0, len(protoSvc.Methods)),
	}

	for _, m := range protoSvc.Methods {
		s.Methods = append(s.Methods, f.ParseMethod(m))
	}

	f.FileDescription.Services = append(f.FileDescription.Services, s)
}

func (f *File) ParseMethod(m *protogen.Method) *Method {
	return &Method{
		Name:         m.GoName,
		RequestType:  m.Input.GoIdent.GoName,
		ResponseType: m.Output.GoIdent.GoName,
	}
}

func (f *File) Generate() error {
	tmpl, err := template.New("example-template").Parse(example_tmpl.HTTP)
	if err != nil {
		return fmt.Errorf("failed to parse example template: %w", err)
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, f.FileDescription); err != nil {
		return fmt.Errorf("failed to execute example template: %w", err)
	}
	f.gen.P(buf.String())
	return nil
}
