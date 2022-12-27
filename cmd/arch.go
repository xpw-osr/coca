package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/modernizing/coca/cmd/cmd_util"
	"github.com/modernizing/coca/cmd/config"
	"github.com/modernizing/coca/pkg/application/arch"
	"github.com/modernizing/coca/pkg/application/arch/tequila"
	"github.com/modernizing/coca/pkg/application/visual"
	"github.com/modernizing/coca/pkg/domain/core_domain"
	"github.com/spf13/cobra"
)

type ArchCmdConfig struct {
	DependencePath string
	IsMergePackage bool
	FilterString   string
	IsMergeHeader  bool
	WithVisual     bool
}

var (
	archCmdConfig ArchCmdConfig
)

var archCmd = &cobra.Command{
	Use:   "arch",
	Short: "project package visualization",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("- loading identifies ... ")
		identifiers = cmd_util.LoadIdentify(apiCmdConfig.DependencePath)
		fmt.Println("done")

		fmt.Printf("- generating identifiers map ...")
		identifiersMap = core_domain.BuildIdentifierMap(identifiers)
		fmt.Print("done")

		fmt.Printf("- loading deps ... ")
		parsedDeps := cmd_util.GetDepsFromJson(archCmdConfig.DependencePath)
		fmt.Println("done")

		fmt.Printf("- analysising ... ")
		archApp := arch.NewArchApp()
		result := archApp.Analysis(parsedDeps, identifiersMap)
		fmt.Println("done")

		filter := strings.Split(archCmdConfig.FilterString, ",")
		var nodeFilter = func(key string) bool {
			for _, f := range filter {
				if strings.Contains(key, f) {
					return true
				}
			}
			return false
		}

		if archCmdConfig.WithVisual {
			fmt.Printf("- generating visual ... ")
			output := visual.FromDeps(parsedDeps)
			out, _ := json.Marshal(output)
			cmd_util.WriteToCocaFile("visual.json", string(out))
			fmt.Println("done")
		}

		if archCmdConfig.IsMergeHeader {
			fmt.Printf("- merging header ... ")
			result = result.MergeHeaderFile(tequila.MergeHeaderFunc)
			fmt.Println("done")
		}

		if archCmdConfig.IsMergePackage {
			fmt.Printf("- merging package ... ")
			result = result.MergeHeaderFile(tequila.MergePackageFunc)
			fmt.Println("done")
		}

		fmt.Printf("- generating dot file ... ")
		graph := result.ToMapDot(nodeFilter)
		f, _ := os.Create("coca_reporter/arch.dot")
		w := bufio.NewWriter(f)
		_, _ = w.WriteString("di" + graph.String())
		_ = w.Flush()
		fmt.Println("done")

		// fmt.Printf("- generate svg ... ")
		// cmd_util.ConvertToSvg("arch")
		// fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(archCmd)

	archCmd.PersistentFlags().StringVarP(&archCmdConfig.DependencePath, "dependence", "d", config.CocaConfig.ReporterPath+"/deps.json", "get dependence file")
	archCmd.PersistentFlags().BoolVarP(&archCmdConfig.IsMergePackage, "mergePackage", "P", false, "merge package")
	archCmd.PersistentFlags().BoolVarP(&archCmdConfig.IsMergeHeader, "mergeHeader", "H", false, "merge header")
	archCmd.PersistentFlags().BoolVarP(&archCmdConfig.WithVisual, "showVisual", "v", false, "build visual json")
	archCmd.PersistentFlags().StringVarP(&archCmdConfig.FilterString, "filter", "x", "", "filter -x com.phodal")
}
