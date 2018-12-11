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

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:                   "split [FLAG]... PAGES OUTPUT_FILE",
	Short:                 "Split PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		pages, err := parsePageRange(args[1])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := pdf.SplitPdf(args[0], args[2], pages); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Must provide the input file, page range and output file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
}
