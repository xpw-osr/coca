package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/modernizing/coca/cmd/cmd_util"
	"github.com/modernizing/coca/cmd/config"

	. "github.com/modernizing/coca/pkg/application/call"
	"github.com/modernizing/coca/pkg/domain/core_domain"
	"github.com/spf13/cobra"
)

type CallCmdConfig struct {
	Path       string
	ClassName  string
	RemoveName string
	Lookup     bool
}

var (
	callCmdConfig CallCmdConfig
)

var callGraphCmd = &cobra.Command{
	Use:   "call",
	Short: "show call graph with specific method",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var parsedDeps []core_domain.CodeDataStruct
		dependence := callCmdConfig.Path

		if dependence != "" {
			fmt.Printf("ClassName: %s, Path: %s, Remove Name: %s\n", callCmdConfig.ClassName, callCmdConfig.Path, callCmdConfig.RemoveName)

			analyser := NewCallGraph()
			file := cmd_util.ReadFile(dependence)
			if file == nil {
				log.Fatal("lost file:" + dependence)
			}

			_ = json.Unmarshal(file, &parsedDeps)

			depCount := len(parsedDeps)
			depIndex := 0
			var contents string
			for _, dep := range parsedDeps {
				depIndex += 1
				fmt.Printf("[%d / %d] Package: %s, Type: %s, NodeName: %s\n", depIndex, depCount, dep.Package, dep.Type, dep.NodeName)

				funcCount := len(dep.Functions)
				funcIndex := 0
				for _, fun := range dep.Functions {
					funcIndex += 1
					className := fmt.Sprintf("%s.%s.%s", dep.Package, dep.NodeName, fun.Name)
					fmt.Printf("  [%d / %d] %s ... ", funcIndex, funcCount, className)

					content := analyser.Analysis(className, parsedDeps, callCmdConfig.Lookup)
					contents = fmt.Sprintf("%s\n%s", contents, content)

					fmt.Println("done")
				}
			}

			dotFileName := fmt.Sprintf("%s.dot", callCmdConfig.ClassName)
			cmd_util.WriteToCocaFile(dotFileName, contents)

			// --------------------------------------------------

			// content := analyser.Analysis(callCmdConfig.ClassName, parsedDeps, callCmdConfig.Lookup)
			// if callCmdConfig.RemoveName != "" {
			// 	content = strings.ReplaceAll(content, callCmdConfig.RemoveName, "")
			// }

			// cmd_util.WriteToCocaFile("call.dot", content)
			// cmd_util.ConvertToSvg("call")
		}
	},
}

func init() {
	rootCmd.AddCommand(callGraphCmd)

	callGraphCmd.PersistentFlags().StringVarP(&callCmdConfig.ClassName, "className", "c", "", "class")
	callGraphCmd.PersistentFlags().StringVarP(&callCmdConfig.Path, "dependence", "d", config.CocaConfig.ReporterPath+"/deps.json", "get dependence file")
	callGraphCmd.PersistentFlags().StringVarP(&callCmdConfig.RemoveName, "remove", "r", "", "remove package ParamName")
	callGraphCmd.PersistentFlags().BoolVarP(&callCmdConfig.Lookup, "lookup", "l", false, "call with rcall")
}
