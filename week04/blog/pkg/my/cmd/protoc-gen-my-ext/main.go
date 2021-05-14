package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	"os"
	"path/filepath"
)

const version = "v1.0.0"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Fprintf(os.Stdout, "--my-ext_out=ext=(errors|https|grpc),paths=source_relative:..\n")
		os.Exit(0)
	}

	var flags flag.FlagSet
	extErr := flags.Bool("errors", false, "gen my.api errors")
	extHttp := flags.Bool("http", false, "gen my.api http")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if *extErr {
				generateErrorFile(gen,f)
			}
			if *extHttp {
				generateHttpFile(gen,f)
			}
		}
		return nil
	})
}
