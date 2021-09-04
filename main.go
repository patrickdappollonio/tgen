package main

import (
	"os"

	"github.com/spf13/cobra"
)

const appName = "tgen"

var (
	version = "development"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var configs conf

	var root = &cobra.Command{
		Use:          appName,
		Short:        appName + " is a template generator with the power of Go Templates",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command(os.Stdout, configs)
		},
	}

	root.Flags().StringVarP(&configs.environmentFile, "environment", "e", "", "an optional environment file to use (key=value formatted) to perform replacements")
	root.Flags().StringVarP(&configs.templateFile, "file", "f", "", "the template file to process (required)")
	root.Flags().StringVarP(&configs.customDelimiters, "delimiter", "d", "", `template delimiter (default "{{}}")`)
	root.Flags().StringVarP(&configs.rawTemplate, "execute", "x", "", "a raw template to execute directly, without providing --file")
	root.Flags().BoolVarP(&configs.stdin, "stdin", "i", false, "a stdin input to execute directly, without providing --file or --execute")
	root.Flags().BoolVarP(&configs.strictMode, "strict", "s", false, "enables strict mode: if an environment variable in the file is defined but not set, it'll fail")

	return root.Execute()
}
