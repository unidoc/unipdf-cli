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
	Use:                   "split [FLAG]... INPUT_FILE PAGES OUTPUT_FILE",
	Short:                 "Split PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputPath := args[2]
		password, _ := cmd.Flags().GetString("password")

		pages, err := parsePageRange(args[1])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		err = pdf.Split(inputPath, outputPath, password, pages)
		if err != nil {
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

	splitCmd.Flags().StringP("password", "p", "", "PDF file password")
}
