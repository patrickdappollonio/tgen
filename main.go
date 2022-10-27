package main

import (
	"os"

	"github.com/spf13/cobra"
)

const appName = "tgen"

var version = "development"

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var configs conf
	var withValues bool

	var root = &cobra.Command{
		Use:          appName,
		Short:        appName + " is a template generator with the power of Go Templates",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if withValues {
				configs.valuesFile = "values.yaml"
			}

			return command(os.Stdout, configs)
		},
	}

	root.Flags().StringVarP(&configs.environmentFile, "environment", "e", "", "an optional environment file to use (key=value formatted) to perform replacements")
	root.Flags().StringVarP(&configs.templateFilePath, "file", "f", "", "the template file to process")
	root.Flags().StringVarP(&configs.customDelimiters, "delimiter", "d", "", `template delimiter (default "{{}}")`)
	root.Flags().StringVarP(&configs.stdinTemplateFile, "execute", "x", "", "a raw template to execute directly, without providing --file")
	root.Flags().StringVarP(&configs.valuesFile, "values", "v", "", "a file containing values to use for the template, a la Helm")
	root.Flags().BoolVar(&withValues, "with-values", false, "automatically include a values.yaml file from the current working directory")
	root.Flags().BoolVarP(&configs.strictMode, "strict", "s", false, "strict mode: if an environment variable or value is used in the template but not set, it fails rendering")

	root.Flags().SortFlags = false

	return root.Execute()
}
