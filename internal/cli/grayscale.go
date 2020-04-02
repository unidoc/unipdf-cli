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

const grayscaleCmdDesc = `Converts the input file to grayscale.

The command can be configured to convert only the specified
pages to grayscale using the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be converted to grayscale, while
page number 5 is skipped.
`

var grayscaleCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s grayscale input_file.pdf", appName),
	fmt.Sprintf("%s grayscale -o output_file input_file.pdf", appName),
	fmt.Sprintf("%s grayscale -o output_file -P 1-3 input_file.pdf", appName),
	fmt.Sprintf("%s grayscale -o output_file -P 1-3 -p pass input_file.pdf", appName),
)

// grayscaleCmd represents the grayscale command.
var grayscaleCmd = &cobra.Command{
	Use:                   "grayscale [FLAG]... INPUT_FILE",
	Short:                 "Convert PDF to grayscale",
	Long:                  grayscaleCmdDesc,
	Example:               grayscaleCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")

		// Parse output file.
		outputPath, _ := cmd.Flags().GetString("output-file")
		if outputPath == "" {
			outputPath = inputPath
		}

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			printUsageErr(cmd, "Invalid page range specified\n")
		}

		// Convert file to grayscale.
		err = pdf.Grayscale(inputPath, outputPath, password, pages)
		if err != nil {
			printErr("Could not convert input file to grayscale: %s\n", err)
		}

		fmt.Printf("Successfully converted %s to grayscale\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must provide the input file")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(grayscaleCmd)

	grayscaleCmd.Flags().StringP("output-file", "o", "", "output file")
	grayscaleCmd.Flags().StringP("password", "p", "", "input file password")
	grayscaleCmd.Flags().StringP("pages", "P", "", "pages to convert to grayscale")
}
