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

const watermarkCmdDesc = `Add watermark to PDF files.

The command can be configured to apply the watermark image only to the specified
pages using the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Watermark will only be applied to pages 1,2,3 (1-3), 4 and 6,7 (6-7), while page
number 5 is skipped.
`

var watermarkCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s watermark input_file.pdf watermark.png", appName),
	fmt.Sprintf("%s watermark -o output file.png input_file.pdf watermark.png", appName),
	fmt.Sprintf("%s watermark -o output file.png -P 1-3 input_file.pdf watermark.png", appName),
	fmt.Sprintf("%s watermark -o output file.png -P 1-3 -p pass input_file.pdf watermark.png", appName),
)

// watermarkCmd represents the watermark command.
var watermarkCmd = &cobra.Command{
	Use:                   "watermark [FLAG]... INPUT_FILE WATERMARK_IMAGE",
	Short:                 "Add watermark to PDF files",
	Long:                  watermarkCmdDesc,
	Example:               watermarkCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		watermarkPath := args[1]
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

		// Apply watermark.
		err = pdf.Watermark(inputPath, outputPath, watermarkPath, password, pages)
		if err != nil {
			printErr("Could not apply watermark to the input file: %s\n", err)
		}

		fmt.Printf("Watermark successfully applied to %s\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide the input file and the watermark image")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(watermarkCmd)

	watermarkCmd.Flags().StringP("output-file", "o", "", "output file")
	watermarkCmd.Flags().StringP("password", "p", "", "input file password")
	watermarkCmd.Flags().StringP("pages", "P", "", "pages on which to add watermark")
}
