/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/pdf"
)

const extractCmdDesc = `Extracts PDF resources.

Supported resources:
  - text
  - images

The extracted text is always printed to STDOUT.

The images are extracted in a ZIP file and saved at the destination specified
by the --output-file parameter. If no output file is specified, the ZIP
archive is saved in the same directory as the input file.

The command can be configured to extract resources only from the specified
pages using the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Resources will only be extracted from pages 1,2,3 (1-3), 4 and 6,7 (6-7), while page
number 5 is skipped.
`

var extractCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s extract -r text input_file.pdf", appName),
	fmt.Sprintf("%s extract -r text -P 1-3 input_file.pdf", appName),
	fmt.Sprintf("%s extract -r text -P 1-3 -p pass input_file.pdf", appName),
	fmt.Sprintf("%s extract -r images input_file.pdf", appName),
	fmt.Sprintf("%s extract -r images -o images.zip input_file.pdf", appName),
	fmt.Sprintf("%s extract -r images -P 1-3 -p pass -o images.zip input_file.pdf", appName),
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:                   "extract [FLAG]... INPUT_FILE",
	Short:                 "Extract PDF resources",
	Long:                  extractCmdDesc,
	Example:               extractCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")
		outputPath, _ := cmd.Flags().GetString("output-file")

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			printUsageErr(cmd, "Invalid page range specified\n")
		}

		// Parse resource.
		resource, _ := cmd.Flags().GetString("resource")
		switch resource {
		case "text":
			text, err := pdf.ExtractText(inputPath, password, pages)
			if err != nil {
				printErr("Could not extract text: %s\n", err)
			}

			fmt.Println(text)
		case "images":
			outputPath, err = pdf.ExtractImages(inputPath, outputPath, password, pages)
			if err != nil {
				printErr("Could not extract images: %s\n", err)
				return
			}

			fmt.Printf("Images successfully extracted to %s\n", outputPath)
		default:
			printUsageErr(cmd, "Invalid resource type\n")
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must provide the input file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringP("user-password", "p", "", "Input file password")
	extractCmd.Flags().StringP("output-file", "o", "", "Output file")
	extractCmd.Flags().StringP("resource", "r", "", "Resource to extract")
	extractCmd.Flags().StringP("pages", "P", "", "Pages to extract resources from")
}
