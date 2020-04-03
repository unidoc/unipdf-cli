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

const splitCmdDesc = `Split PDF files.

The command is used to extract one or more page ranges from the input file
and save the result as the output file.
If no page range is specified, all the pages from the input file will be
copied to the output file.

An example of the pages parameter: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be present in the output file,
while page number 5 is skipped.
`

var splitCmdExample = fmt.Sprintf("%s\n%s\n",
	fmt.Sprintf("%s split input_file.pdf output_file.pdf 1-2", appName),
	fmt.Sprintf("%s split -p pass input_file.pd output_file.pdf 1-2,4", appName),
)

// splitCmd represents the split command.
var splitCmd = &cobra.Command{
	Use:                   "split [FLAG]... INPUT_FILE OUTPUT_FILE [PAGES]",
	Short:                 "Split PDF files",
	Long:                  splitCmdDesc,
	Example:               splitCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputPath := args[1]
		password, _ := cmd.Flags().GetString("password")

		// Parse page range.
		var err error
		var pages []int

		if len(args) > 2 {
			if pages, err = parsePageRange(args[2]); err != nil {
				printUsageErr(cmd, "Invalid page range specified\n")
			}
		}

		err = pdf.Split(inputPath, outputPath, password, pages)
		if err != nil {
			printErr("Error: %v\n", err)
		}

		fmt.Printf("Successfully split file %s\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide at least the input and output files")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().StringP("password", "p", "", "input file password")
}
