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

// grayscaleCmd represents the grayscale command
var grayscaleCmd = &cobra.Command{
	Use:                   "grayscale [FLAG]... INPUT_FILE OUTPUT_FILE",
	Short:                 "Convert PDF to grayscale",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile := args[1]
		password, _ := cmd.Flags().GetString("password")

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		err = pdf.Grayscale(inputFile, outputFile, password, pages)
		if err != nil {
			fmt.Println("Could not convert input file to grayscale")
			return
		}

		fmt.Println("Successfully converted PDF to grayscale")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Must provide the input file, and output file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(grayscaleCmd)

	grayscaleCmd.Flags().StringP("password", "p", "", "PDF file password")
	grayscaleCmd.Flags().StringP("pages", "P", "", "Pages to convert to grayscale")
}
