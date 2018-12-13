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

var mergeCmd = &cobra.Command{
	Use:                   "merge [FLAG]... OUTPUT_FILE INPUT_FILE...",
	Short:                 "Merge PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pdf.Merge(args[1:], args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Must provide the output file and at least two input files\n")
		}

		return nil
	},
}

func init() {
	// Add current command to parent.
	rootCmd.AddCommand(mergeCmd)

	// Add flags.
	mergeCmd.Flags().StringP("password", "p", "", "Help message for toggle")
}
