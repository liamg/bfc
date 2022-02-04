package main

import (
	"fmt"
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

func main() {
	rootCmd.Flags().BoolVarP(&assemble, "assemble", "a", assemble, "Assemble the result (requires nasm)")
	rootCmd.Flags().BoolVarP(&link, "link", "l", link, "Link the result (requires ld) (implies setting of -a flag)")
	rootCmd.Flags().BoolVarP(&run, "run", "r", run, "Run the executable (implies setting of -al flags)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bfc [file]",
	Short: "Compile Brainfuck to Linux x64 assembly.",
	Args:  cobra.ExactArgs(1),
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
		binaryFilename := strings.TrimSuffix(sourceFilename, filepath.Ext(sourceFilename))
		objFilename := binaryFilename + ".o"
		outFilename := binaryFilename + ".asm"
		output, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		fmt.Printf("Compiling %s...\n", sourceFilename)
		err = compiler.Compile(input, output, x86_64.New())
		_ = output.Close()
		if err != nil {
			_ = os.Remove(output.Name())
			return err
		}
		fmt.Printf("Compiled x64 assembly written to %s\n", outFilename)
		if assemble {
			fmt.Printf("Assembling %s with nasm...\n", outFilename)
			if _, err := exec.LookPath("nasm"); err != nil {
				return fmt.Errorf("nasm is not installed, or not available on your PATH: %s", err)
			}
			if err := exec.Command("nasm", "-f", "elf64", "-o", objFilename, outFilename).Run(); err != nil {
				return fmt.Errorf("assemble failed: %s", err)
			}
			fmt.Printf("Assembled object file written to %s\n", objFilename)
		}
		if link {
			fmt.Printf("Linking %s with ld...\n", objFilename)
			if _, err := exec.LookPath("ld"); err != nil {
				return fmt.Errorf("ld is not installed, or not available on your PATH: %s", err)
			}
			if err := exec.Command("ld", "-o", binaryFilename, objFilename).Run(); err != nil {
				return fmt.Errorf("linking failed: %s", err)
			}
			fmt.Printf("Linked executable file written to %s\n", binaryFilename)

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
		fmt.Println("Done.")
		return nil
	},
}
