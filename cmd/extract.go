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

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:                   "extract [FLAG]... INPUT_FILE",
	Short:                 "Extract PDF resources",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		password, _ := cmd.Flags().GetString("password")
		outputFile, _ := cmd.Flags().GetString("output-file")

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		resource, _ := cmd.Flags().GetString("resource")
		switch resource {
		case "text":
			text, err := pdf.ExtractText(inputFile, password, pages)
			if err != nil {
				fmt.Println("Could not extract text")
				return
			}

			fmt.Println(text)
		case "images":
			err := pdf.ExtractImages(inputFile, outputFile, password, pages)
			if err != nil {
				fmt.Println("Could not extract images")
				return
			}
			fmt.Println("Images successfully extracted")
		default:
			fmt.Println("Invalid resource type")
			return
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must provide the input file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringP("user-password", "p", "", "PDF file password")
	extractCmd.Flags().StringP("output-file", "o", "", "Output file")
	extractCmd.Flags().StringP("resource", "r", "", "Resource to extract")
	extractCmd.Flags().StringP("pages", "P", "", "Pages to extract resources from")
}
