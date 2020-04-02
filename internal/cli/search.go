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

const searchCmdDesc = `Search text in PDF files`

var searchCmdExample = fmt.Sprintf("%s\n%s\n",
	fmt.Sprintf("%s search input_file.pdf text_to_search", appName),
	fmt.Sprintf("%s search -p pass input_file.pdf text_to_search", appName),
)

// searchCmd represents the search command.
var searchCmd = &cobra.Command{
	Use:                   "search [FLAG]... INPUT_FILE TEXT",
	Short:                 "Search text in PDF files",
	Long:                  searchCmdDesc,
	Example:               searchCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		text := args[1]
		password, _ := cmd.Flags().GetString("password")

		// Search text.
		results, err := pdf.Search(inputPath, text, password)
		if err != nil {
			printErr("Could not search the specified text: %s\n", err)
		}

		// Print results.
		fmt.Printf("Search results for term: %s\n", text)

		totalOccurrences := 0
		for _, result := range results {
			totalOccurrences += result.Occurrences
			fmt.Printf("Page %d: %d occurrences\n", result.Page, result.Occurrences)
		}

		fmt.Printf("Total occurrences: %d\n", totalOccurrences)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide a PDF file and the text to search")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("password", "p", "", "input file password")
}
