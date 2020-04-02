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

const explodeCmdDesc = `Splits the input file into separate single page PDF files.

The resulting PDF files are saved in a ZIP archive at the location specified
by the --output-file parameter. If no output file is specified, the ZIP file
is saved in the same directory as the input file.

The command can be configured to extract only the specified pages using
the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be extracted, while page
number 5 is skipped.
`

var explodeCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s explode input_file.pdf", appName),
	fmt.Sprintf("%s explode -o pages.zip input_file.pdf", appName),
	fmt.Sprintf("%s explode -o pages.zip -P 1-3 input_file.pdf", appName),
	fmt.Sprintf("%s explode -o pages.zip -P 1-3 -p pass input_file.pdf", appName),
)

// explodeCmd represents the explode command.
var explodeCmd = &cobra.Command{
	Use:                   "explode [FLAG]... INPUT_FILE",
	Short:                 "Explodes the input file into separate single page PDF files",
	Long:                  explodeCmdDesc,
	Example:               explodeCmdExample,
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

		// Explode file.
		outputPath, err = pdf.Explode(inputPath, outputPath, password, pages)
		if err != nil {
			printErr("Could not explode input file: %s\n", err)
			return
		}

		fmt.Printf("File %s successfully exploded\n", inputPath)
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
	rootCmd.AddCommand(explodeCmd)

	explodeCmd.Flags().StringP("password", "p", "", "input file password")
	explodeCmd.Flags().StringP("output-file", "o", "", "output file")
	explodeCmd.Flags().StringP("pages", "P", "", "pages to extract from the input file")
}
