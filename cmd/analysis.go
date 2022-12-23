package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/modernizing/coca/cmd/cmd_util"
	"github.com/modernizing/coca/pkg/application/analysis/javaapp"
	"github.com/modernizing/coca/pkg/domain/core_domain"
	"github.com/spf13/cobra"
)

type AnalysisCmdConfig struct {
	Path           string
	UpdateIdentify bool
}

var (
	analysisCmdConfig AnalysisCmdConfig
)

var analysisCmd = &cobra.Command{
	Use:   "analysis",
	Short: "analysis code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// var outputName string
		// var ds []core_domain.CodeDataStruct

		// ds = AnalysisJava()
		// outputName = "deps.json"

		// cModel, _ := json.MarshalIndent(ds, "", "\t")
		// cmd_util.WriteToCocaFile(outputName, string(cModel))

		callNodes := AnalysisJava()
		callNodesJson, _ := json.MarshalIndent(callNodes, "", "\t")
		cmd_util.WriteToCocaFile("deps.json", string(callNodesJson))
	},
}

func AnalysisJava() []map[string]interface{} {
	importPath := analysisCmdConfig.Path
	var iNodes []core_domain.CodeDataStruct

	if analysisCmdConfig.UpdateIdentify {
		fmt.Println("# Updating identify ... ")

		identifierApp := javaapp.NewJavaIdentifierApp()
		jsonFiles := identifierApp.AnalysisPath(importPath)

		for _, file := range jsonFiles {
			if contents, succeed := javaapp.ReadJsonFile(file); succeed {
				var subNodes []core_domain.CodeDataStruct
				json.Unmarshal(contents, &subNodes)

				iNodes = append(iNodes, subNodes...)
			}
		}

		identModel, _ := json.MarshalIndent(iNodes, "", "\t")
		cmd_util.WriteToCocaFile("identify.json", string(identModel))
	} else {
		fmt.Println("use local identify")
		identContent := cmd_util.ReadCocaFile("identify.json")
		_ = json.Unmarshal(identContent, &iNodes)
	}

	fmt.Println("# Generating deps ... ")

	callApp := javaapp.NewJavaFullApp()
	callJsonFiles := callApp.AnalysisPath(importPath, iNodes)

	var callNodes []map[string]interface{}
	for _, file := range callJsonFiles {
		if contents, succeed := javaapp.ReadJsonFile(file); succeed {
			var subNodes []map[string]interface{}
			json.Unmarshal(contents, &subNodes)

			callNodes = append(callNodes, subNodes...)
		}
	}

	return callNodes
}

func init() {
	rootCmd.AddCommand(analysisCmd)

	analysisCmd.PersistentFlags().StringVarP(&analysisCmdConfig.Path, "path", "p", ".", "example -p core/main")
	//analysisCmd.PersistentFlags().StringVarP(&analysisCmdConfig.Lang, "lang", "l", "java", "example coca analysis -l java, typescript, python")
	analysisCmd.PersistentFlags().BoolVarP(&analysisCmdConfig.UpdateIdentify, "identify", "i", true, "use current identify")
}
