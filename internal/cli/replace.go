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

const replaceCmdDesc = `Replace text in PDF files`

var replaceCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s replace input_file.pdf text_to_search", appName),
	fmt.Sprintf("%s replace -o output_file input_file.pdf text_to_search", appName),
	fmt.Sprintf("%s replace -o output_file -r new_text input_file.pdf text_to_search", appName),
	fmt.Sprintf("%s replace -o output_file  -r new_text -p pass input_file.pdf text_to_search", appName),
)

// replaceCmd represents the replace command.
var replaceCmd = &cobra.Command{
	Use:                   "replace [FLAG]... INPUT_FILE TEXT",
	Short:                 "Replace text in PDF files",
	Long:                  replaceCmdDesc,
	Example:               replaceCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		text := args[1]
		password, _ := cmd.Flags().GetString("password")

		// Parse output file.
		outputPath, _ := cmd.Flags().GetString("output-file")
		if outputPath == "" {
			outputPath = inputPath
		}

		// Parse replaceText.
		replaceText, _ := cmd.Flags().GetString("replace-text")
		if replaceText == "" {
			replaceText = text
		}

		// Search text.
		err := pdf.Replace(inputPath, outputPath, text, replaceText, password)
		if err != nil {
			printErr("Could not replace the specified text: %s\n", err)
		}

		fmt.Printf("Successfully replaced text %s with %s\n", text, replaceText)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide a PDF file and the text to search")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	replaceCmd.Flags().StringP("output-file", "o", "", "output file")
	replaceCmd.Flags().StringP("replace-text", "r", "", "replacement text")
	replaceCmd.Flags().StringP("password", "p", "", "input file password")
}
