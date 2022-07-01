package main

import (
	"os"

	"github.com/go-swagger/go-swagger/generator"
)

func usage() {
	println("A minimal OpenAPI generator")
	println("Usage:")
	println("rest-gen [model|markdown] <open-api-spec.yaml> [output.md for markdown generation]")
}

func main() {
	if len(os.Args) == 1 || len(os.Args) > 4 {
		usage()
		return
	}
	mode := os.Args[1]
	specFile := os.Args[2]

	modelNames := []string{}
	opts := generator.GenOpts{}
	opts.EnsureDefaults()
	opts.Spec = specFile
	//opts.StructTags = append(opts.StructTags, "bson")

	if mode == "model" {
		generator.GenerateModels(modelNames, &opts)
	} else if mode == "markdown" {
		outFile := os.Args[3]
		operationIDs := []string{}
		generator.GenerateMarkdown(outFile, modelNames, operationIDs, &opts)
	} else {
		panic("Unrecognized param: " + mode)
	}
}
