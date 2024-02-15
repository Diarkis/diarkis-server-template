package main

import (
	"fmt"
	"os"

	"github.com/Diarkis/diarkis/puffer"
)

func main() {
	path := os.Args[1]
	definitionPath, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	definitionPath = fmt.Sprintf("%s/bad_json_definitions/%s", definitionPath, path)

	fmt.Println("Path", definitionPath)

	outputPath := fmt.Sprint(definitionPath, "/src")

	params := &puffer.Params{
		DefinitionPath: definitionPath,
		OutputPath:     outputPath,
		ModulePath:     "github.com/Diarkis/diarkis/puffer/test/sample/src/go",
		TPLPath:        puffer.TplDefaultPath,
	}
	puffer.Generate("go", params)
}
