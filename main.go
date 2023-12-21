package main

import (
	"github.com/bootun/protoc-gen-go-example/parser"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if err := parser.GenerateFile(gen, f); err != nil {
				return err
			}
		}
		return nil
	})
}
