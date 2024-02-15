package main

import (
	"flag"
	"os"

	"github.com/Diarkis/diarkis/puffer"
)

func main() {
	root, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	langStr := flag.String("lang", "all", "output source language [all;cs;go;cpp]")
	outputPath := flag.String("ouput", "./payloads", "Definition path.")
	definitionPath := flag.String("definitions", "./json_definitions", "input folder for json definition files.")
	tplPath := flag.String("tpl", puffer.TplDefaultPath, "folder containing the tpl files (language syntax parsing format definition)")
	modulePath := flag.String("module", "github.com/Diarkis/diarkis/puffer/test/sample/src/go", "go module path for importing generated module")
	helpFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *helpFlag {
		flag.CommandLine.Usage()
		os.Exit(0)
	}

	params := &puffer.Params{
		DefinitionPath: root + "/" + *definitionPath,
		OutputPath:     root + "/" + *outputPath,
		ModulePath:     *modulePath,
		TPLPath:        *tplPath,
	}

	if *langStr == "all" {
		puffer.Generate("cs", params)
		puffer.Generate("go", params)
		puffer.Generate("cpp", params)
	} else {
		puffer.Generate(*langStr, params)
	}
}
