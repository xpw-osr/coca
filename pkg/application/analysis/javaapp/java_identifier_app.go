package javaapp

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/modernizing/coca/pkg/adapter/cocafile"
	"github.com/modernizing/coca/pkg/domain/core_domain"
	"github.com/modernizing/coca/pkg/infrastructure/ast/ast_java"
	"github.com/modernizing/coca/pkg/infrastructure/ast/ast_java/java_identify"
)

type JavaIdentifierApp struct {
}

func NewJavaIdentifierApp() JavaIdentifierApp {
	return JavaIdentifierApp{}
}

func (j *JavaIdentifierApp) AnalysisPath(codeDir string) []string {
	files := cocafile.GetJavaFiles(codeDir)
	fmt.Printf("File Count: %d\n", len(files))

	return j.AnalysisFiles(files)
}

func (j *JavaIdentifierApp) AnalysisFiles(files []string) []string {
	var identJsonFiles []string = nil

	fileCount := len(files)
	fileIndex := 0
	for _, file := range files {
		fileIndex += 1
		fmt.Printf("- [%d / %d] %s ... ", fileIndex, fileCount, file)

		identifiers := j.AnalysisFile(file)
		nodesJson, _ := json.MarshalIndent(identifiers, "", "\t")

		filename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		identJsonFilename := fmt.Sprintf("%s.json", filename)
		if path, succeed := writeJsonFile(identJsonFilename, "ident", string(nodesJson)); succeed {
			identJsonFiles = append(identJsonFiles, path)
			fmt.Println("done")
		} else {
			fmt.Println("failed")
		}
	}

	return identJsonFiles
}

func (j *JavaIdentifierApp) AnalysisFile(file string) []core_domain.CodeDataStruct {
	parser := ast_java.ProcessJavaFile(file)
	context := parser.CompilationUnit()
	listener := java_identify.NewJavaIdentifierListener()

	antlr.NewParseTreeWalker().Walk(listener, context)

	return listener.GetNodes()
}
