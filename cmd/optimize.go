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

// optimizeCmd represents the optimize command
var optimizeCmd = &cobra.Command{
	Use:                   "optimize [FLAG]... INPUT_FILE",
	Short:                 "Optimize PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		password, _ := cmd.Flags().GetString("password")

		outputFile, _ := cmd.Flags().GetString("output-file")
		if outputFile == "" {
			outputFile = inputFile
		}

		imageQuality, err := cmd.Flags().GetInt("image-quality")
		if err != nil {
			imageQuality = 100
		}

		opts := &pdf.OptimizeOpts{
			ImageQuality: imageQuality,
		}

		err = pdf.OptimizePdf(inputFile, outputFile, password, opts)
		if err != nil {
			fmt.Println("Could not optimize input file")
			return
		}

		fmt.Println("Input file sucessfully optimized")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must provide the input file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)

	optimizeCmd.Flags().StringP("password", "p", "", "File password")
	optimizeCmd.Flags().StringP("output-file", "o", "", "Output file")
	optimizeCmd.Flags().IntP("image-quality", "q", 100, "Optimized image quality")
}
