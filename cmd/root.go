package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kylejryan/go-vuln-scan/internal/analyzer"
	"github.com/kylejryan/go-vuln-scan/internal/fileutils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosecscan [path_to_code]",
	Short: "GoSecScan scans code for security vulnerabilities using AI.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		files, err := fileutils.FindFiles(path)
		if err != nil {
			log.Fatalf("Error finding files: %v", err)
		}

		for _, file := range files {
			analysis, err := analyzer.Analyze(file)
			if err != nil {
				fmt.Printf("Error analyzing file %s: %v\n", file, err)
				continue
			}
			fmt.Printf("Results for %s:\n%s\n---------------------------\n\n", file, analysis)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
