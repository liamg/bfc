package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/liamg/bfc/pkg/backend/generator/x86_64"
	"github.com/liamg/bfc/pkg/compiler"
	"github.com/spf13/cobra"
)

var assemble bool
var link bool
var run bool
var output string

func main() {
	rootCmd.Flags().BoolVarP(&assemble, "assemble", "a", assemble, "Assemble the result (requires nasm)")
	rootCmd.Flags().BoolVarP(&link, "link", "l", link, "Link the result (requires ld) (implies setting of -a flag)")
	rootCmd.Flags().BoolVarP(&run, "run", "r", run, "Run the executable (implies setting of -al flags)")
	rootCmd.Flags().StringVarP(&output, "output", "o", output, "Output file")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:           "bfc [file]",
	Short:         "Compile Brainfuck to Linux x64 assembly.",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if run {
			link = true
		}
		if link {
			assemble = true
		}

		cmd.SilenceUsage = true
		sourceFilename := args[0]
		input, err := os.Open(sourceFilename)
		if err != nil {
			return err
		}
		baseName := strings.TrimSuffix(filepath.Base(sourceFilename), filepath.Ext(sourceFilename))

		tmp, err := os.MkdirTemp(os.TempDir(), baseName)
		if err != nil {
			return err
		}

		asmFilename := filepath.Join(tmp, fmt.Sprintf("%s.asm", baseName))
		objFilename := filepath.Join(tmp, fmt.Sprintf("%s.o", baseName))
		binaryFilename := filepath.Join(tmp, baseName)
		finalFile := asmFilename
		abs, err := filepath.Abs(sourceFilename)
		if err != nil {
			return err
		}
		publishDir := filepath.Dir(abs)

		asmFile, err := os.Create(asmFilename)
		if err != nil {
			return err
		}
		err = compiler.Compile(input, asmFile, x86_64.New())
		_ = asmFile.Close()
		if err != nil {
			return err
		}
		if assemble {
			if _, err := exec.LookPath("nasm"); err != nil {
				return fmt.Errorf("nasm is not installed, or not available on your PATH: %s", err)
			}
			if err := exec.Command("nasm", "-f", "elf64", "-o", objFilename, asmFilename).Run(); err != nil {
				return fmt.Errorf("assemble failed: %s", err)
			}
			finalFile = objFilename
		}
		if link {
			if _, err := exec.LookPath("ld"); err != nil {
				return fmt.Errorf("ld is not installed, or not available on your PATH: %s", err)
			}
			if err := exec.Command("ld", "-o", binaryFilename, objFilename).Run(); err != nil {
				return fmt.Errorf("linking failed: %s", err)
			}
			finalFile = binaryFilename
		}

		if finalFile != "" {
			published := filepath.Join(publishDir, filepath.Base(finalFile))
			if output != "" {
				published = output
			}
			buffer, err := ioutil.ReadFile(finalFile)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(published, buffer, 0700); err != nil {
				return err
			}
			binaryFilename = published
		}

		if run {
			cmd := exec.Command(binaryFilename)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("executable returned non-zero status: %s", err)
			}
		}
		return nil
	},
}
