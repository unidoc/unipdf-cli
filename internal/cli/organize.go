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

const organizeCmdDesc = `Split PDF files.

The command is used to organize one or more page ranges from the input file
and save the result as the output file.
If no page range is specified, all the pages from the input file will be
copied to the output file.

An example of the pages parameter: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be present in the output file,
while page number 5 is skipped.
`

var organizeCmdExample = fmt.Sprintf("%s\n%s\n",
	fmt.Sprintf("%s organize input_file.pdf output_file.pdf 1-2", appName),
	fmt.Sprintf("%s organize -p pass input_file.pd output_file.pdf 1-2,4", appName),
)

// organizeCmd represents the split command.
var organizeCmd = &cobra.Command{
	Use:                   "organize [FLAG]... INPUT_FILE OUTPUT_FILE [PAGES]",
	Short:                 "Organize PDF files",
	Long:                  organizeCmdDesc,
	Example:               organizeCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputPath := args[1]
		password, _ := cmd.Flags().GetString("password")

		// Parse page range.
		var err error
		var pages []int

		if len(args) > 2 {
			if pages, err = parsePageRangeUnsorted(args[2]); err != nil {
				printUsageErr(cmd, "Invalid page range specified\n")
			}
		}

		if err := pdf.Organize(inputPath, outputPath, password, pages); err != nil {
			printErr("Error: %s\n", err)
		}

		fmt.Printf("Successfully organized file %s\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide at least the input and output files")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(organizeCmd)

	organizeCmd.Flags().StringP("password", "p", "", "input file password")
}
