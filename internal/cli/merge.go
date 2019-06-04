/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const mergeCmdDesc = `Merge the provided input files and save the result to the
specified output file.`

var mergeCmdExample = fmt.Sprintf("%s\n",
	fmt.Sprintf("%s merge output_file.pdf input_file1.pdf input_file2.pdf", appName),
)

var mergeCmd = &cobra.Command{
	Use:                   "merge [FLAG]... OUTPUT_FILE INPUT_FILE...",
	Short:                 "Merge PDF files",
	Long:                  mergeCmdDesc,
	Example:               mergeCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := args[0]
		inputPaths := args[1:]

		if err := pdf.Merge(inputPaths, outputPath); err != nil {
			printErr("Could not merge the input files: %s\n", err)
		}

		fmt.Printf("Successfully merged input files\n")
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("must provide the output file and at least two input files")
		}

		return nil
	},
}

func init() {
	// Add current command to parent.
	rootCmd.AddCommand(mergeCmd)
}
