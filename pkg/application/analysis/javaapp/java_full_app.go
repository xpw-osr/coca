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
)

type JavaFullApp struct {
}

func NewJavaFullApp() JavaFullApp {
	return JavaFullApp{}
}

func (j *JavaFullApp) AnalysisPath(codeDir string, identNodes []core_domain.CodeDataStruct) []string {
	files := cocafile.GetJavaFiles(codeDir)
	return j.AnalysisFiles(identNodes, files)
}

func (j *JavaFullApp) AnalysisFiles(identNodes []core_domain.CodeDataStruct, files []string) []string {
	// var classes []string = nil
	var identMap = make(map[string]core_domain.CodeDataStruct)
	for _, node := range identNodes {
		// classes = append(classes, node.GetClassFullName())
		identMap[node.GetClassFullName()] = node
	}

	fileCount := len(files)
	fileIndex := 0
	var jsonFiles []string = nil
	for _, file := range files {
		fileIndex += 1
		if needSkip(file) {
			fmt.Printf("- [%d / %d] %s ... skip\n", fileIndex, fileCount, file)
			continue
		}
		fmt.Printf("- [%d / %d] %s ... ", fileIndex, fileCount, file)

		// nodes := j.AnalysisFile(file, identMap, classes)
		nodes := j.AnalysisFile(file, identMap, nil)
		nodesJson, _ := json.MarshalIndent(nodes, "", "\t")

		filename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		depJsonFilename := fmt.Sprintf("%s.json", filename)
		if path, succeed := writeJsonFile(depJsonFilename, "deps", string(nodesJson)); succeed {
			jsonFiles = append(jsonFiles, path)
			fmt.Println("done")
		} else {
			fmt.Println("failed")
		}
	}

	return jsonFiles
}

func (j *JavaFullApp) AnalysisFile(file string, identMap map[string]core_domain.CodeDataStruct, classes []string) []core_domain.CodeDataStruct {
	parser := ast_java.ProcessJavaFile(file)
	context := parser.CompilationUnit()

	listener := ast_java.NewJavaFullListener(identMap, file)
	// listener.AppendClasses(classes)

	antlr.NewParseTreeWalker().Walk(listener, context)

	return listener.GetNodeInfo()
}

var skips = []string{
	// out of memory
	"frameworks/base/core/java/android/app/Notification.java",
	"frameworks/base/core/java/android/content/pm/PackageParser.java",
	"frameworks/base/core/java/android/hardware/soundtrigger/SoundTrigger.java",
	"frameworks/base/core/java/android/provider/ContactsContract.java",
	"frameworks/base/core/java/android/widget/RemoteViews.java",
	"frameworks/base/core/java/com/android/internal/widget/RecyclerView.java",
	"frameworks/base/packages/SystemUI/src/com/android/systemui/globalactions/GlobalActionsDialogLite.java",
	"frameworks/base/services/core/java/com/android/server/notification/NotificationManagerService.java",
	"frameworks/wilhelm/tests/native-media/src/com/example/nativemedia/NativeMedia.java",
	// panic: interface conversion: antlr.Tree is *parser.ReceiverParameterContext, not *parser.FormalParameterListContext
	"frameworks/base/services/core/java/com/android/server/location/eventlog/LocalEventLog.java",
}

func needSkip(file string) bool {
	for _, skip := range skips {
		if file == skip {
			return true
		}
	}
	return false
}
