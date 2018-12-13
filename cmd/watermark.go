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

// watermarkCmd represents the watermark command
var watermarkCmd = &cobra.Command{
	Use:                   "watermark [FLAG]... INPUT_FILE WATERMARK_IMAGE OUTPUT_FILE",
	Short:                 "Add watermark to PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		watermarkFile := args[1]
		outputFile := args[2]

		password, _ := cmd.Flags().GetString("password")
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		err = pdf.Watermark(inputFile, outputFile, watermarkFile, password, pages)
		if err != nil {
			fmt.Println("Could not apply watermark to the input file")
			return
		}

		fmt.Println("Watermark sucessfully applied to the output file")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Must provide the input file, watermark image and output file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(watermarkCmd)

	watermarkCmd.Flags().StringP("password", "p", "", "PDF file password")
	watermarkCmd.Flags().StringP("pages", "P", "", "Pages on which to add watermark")
}
