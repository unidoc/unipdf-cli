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

const extractTextCmdDesc = `Extracts PDF text.

The extracted text is always printed to STDOUT.

The command can be configured to extract text only from the specified pages
using the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Text will only be extracted from pages 1,2,3 (1-3), 4 and 6,7 (6-7), while page
number 5 is skipped.
`

var extractTextCmdExample = fmt.Sprintf("%s\n%s\n%s\n",
	fmt.Sprintf("%s extract text input_file.pdf", appName),
	fmt.Sprintf("%s extract text -P 1-3 input_file.pdf", appName),
	fmt.Sprintf("%s extract text -P 1-3 -p pass input_file.pdf", appName),
)

// extractTextCmd represents the extract text command.
var extractTextCmd = &cobra.Command{
	Use:                   "text [FLAG]... INPUT_FILE",
	Short:                 "Extract PDF text",
	Long:                  extractTextCmdDesc,
	Example:               extractTextCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			printUsageErr(cmd, "Invalid page range specified\n")
		}

		// Extract text.
		text, err := pdf.ExtractText(inputPath, password, pages)
		if err != nil {
			printErr("Could not extract text: %s\n", err)
		}

		fmt.Println(text)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must provide the input file")
		}

		return nil
	},
}

func init() {
	extractCmd.AddCommand(extractTextCmd)

	extractTextCmd.Flags().StringP("password", "p", "", "input file password")
	extractTextCmd.Flags().StringP("pages", "P", "", "pages to extract text from")
}
