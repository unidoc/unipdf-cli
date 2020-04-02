/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const rotateCmdDesc = `Rotate PDF file pages by a specified angle.
The angle argument is specified in degrees and it must be a multiple of 90.

The command can be configured to rotate only the specified pages
using the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be rotated, while
page number 5 is skipped.
`

var rotateCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s rotate input_file.pdf 90", appName),
	fmt.Sprintf("%s rotate -- input_file.pdf -270", appName),
	fmt.Sprintf("%s rotate -o output_file.pdf input_file.pdf 90", appName),
	fmt.Sprintf("%s rotate -o output_file.pdf -P 1-3 input_file.pdf 90", appName),
	fmt.Sprintf("%s rotate -o output_file.pdf -P 1-3 -p pass input_file.pdf 90", appName),
)

// rotateCmd represents the rotate command.
var rotateCmd = &cobra.Command{
	Use:                   "rotate [FLAG]... INPUT_FILE ANGLE",
	Short:                 "Rotate PDF file pages",
	Long:                  rotateCmdDesc,
	Example:               rotateCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")

		// Parse angle parameter.
		angle, err := strconv.Atoi(args[1])
		if err != nil {
			printUsageErr(cmd, "Invalid rotation angle specified\n")
		}

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

		// Rotate file.
		outputPath, err = pdf.Rotate(inputPath, outputPath, angle, password, pages)
		if err != nil {
			printErr("Could not rotate input file pages: %s\n", err)
		}

		fmt.Printf("Successfully rotated %s\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide the input file and the rotation angle")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)

	rotateCmd.Flags().StringP("pages", "P", "", "pages to rotate")
	rotateCmd.Flags().StringP("output-file", "o", "", "putput file")
	rotateCmd.Flags().StringP("password", "p", "", "input file password")
}
