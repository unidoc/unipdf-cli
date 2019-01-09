/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unicli/pdf"
)

const optimizeCmdDesc = `Optimize PDF files.

The quality of the images in the output file can be configured. (see the --image-quality flag)
`

var optimizeCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s optimize input_file.pdf", appName),
	fmt.Sprintf("%s optimize -o output_file input_file.pdf", appName),
	fmt.Sprintf("%s optimize -o output_file -i 75 input_file.pdf", appName),
	fmt.Sprintf("%s optimize -o output_file -i 75 -p pass input_file.pdf", appName),
)

// optimizeCmd represents the optimize command
var optimizeCmd = &cobra.Command{
	Use:                   "optimize [FLAG]... INPUT_FILE",
	Short:                 "Optimize PDF files",
	Long:                  optimizeCmdDesc,
	Example:               optimizeCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")

		// Parse output file.
		outputPath, _ := cmd.Flags().GetString("output-file")
		if outputPath == "" {
			outputPath = inputPath
		}

		// Parse optimization parameters.
		imageQuality, err := cmd.Flags().GetInt("image-quality")
		if err != nil {
			imageQuality = 100
		}

		opts := &pdf.OptimizeOpts{
			ImageQuality: imageQuality,
		}

		// Optimize PDF.
		err = pdf.Optimize(inputPath, outputPath, password, opts)
		if err != nil {
			printErr("Could not optimize input file: %s\n", err)
		}

		fmt.Printf("Input file %s successfully optimized\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must provide the input file")
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
