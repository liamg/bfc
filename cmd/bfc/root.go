package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/liamg/bfc/pkg/backend/generators/x86_64"
	"github.com/liamg/bfc/pkg/compiler"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:          "bfc [file]",
	Short:        "Compile Brainfuck to a Linux x64_64 binary",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		sourceFilename := args[0]
		input, err := os.Open(sourceFilename)
		if err != nil {
			return err
		}
		outFilename := strings.TrimSuffix(sourceFilename, filepath.Ext(sourceFilename))
		output, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		defer func() {
			_ = output.Close()
			_ = os.Remove(output.Name())
		}()

		gen := x86_64.New()
		if err := compiler.Compile(input, output, gen); err != nil {
			return err
		}
		fmt.Printf("Compilation successful! Binary written to %s.\n", output.Name())
		return nil
	},
}
