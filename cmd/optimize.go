/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unidoc/unicli/pdf"
)

const optimizeCmdDesc = `Optimize PDF files.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_optimized" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the optimized output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.

The quality of the images in the output files can be configured through
the --image-quality flag.
`

var optimizeCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s optimize file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s optimize -O file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s optimize -O -r file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t output_dir file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t output_dir -r file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t output_dir -r -i 75 file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t output_dir -r -i 75 -p pass file_1.pdf file_n.pdf dir_1 dir_n", appName),
)

// optimizeCmd represents the optimize command
var optimizeCmd = &cobra.Command{
	Use:                   "optimize [FLAG]... INPUT_FILES...",
	Short:                 "Optimize PDF files",
	Long:                  optimizeCmdDesc,
	Example:               optimizeCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags.
		outputDir, _ := cmd.Flags().GetString("target-dir")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		recursive, _ := cmd.Flags().GetBool("recursive")
		password, _ := cmd.Flags().GetString("password")

		// Parse optimization parameters.
		imageQuality, err := cmd.Flags().GetInt("image-quality")
		if err != nil {
			imageQuality = 100
		}

		opts := &pdf.OptimizeOpts{
			ImageQuality: clampInt(imageQuality, 10, 100),
		}

		// Parse input parameters.
		inputPaths, err := parseInputPaths(args, recursive, pdfMatcher)
		if err != nil {
			printErr("Could not parse input files: %s\n", err)
		}

		// Create output directory, if it does not exist.
		if outputDir != "" {
			if overwrite {
				printErr("The --target-dir and the --overwrite flags are mutually exclusive")
			}

			if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
				printErr("Could not create output directory: %s\n", err)
			}
		}

		// Optimize PDF files.
		for _, inputPath := range inputPaths {
			fmt.Printf("Optimizing %s\n", inputPath)

			// Generate output path.
			outputPath := generateOutputPath(inputPath, outputDir, "optimized", overwrite)

			// Optimize input file.
			err = pdf.Optimize(inputPath, outputPath, password, opts)
			if err != nil {
				printErr("Could not optimize input file: %s\n", err)
			}

			fmt.Println("Status: success")
			fmt.Printf("Output file: %s\n", outputPath)
			fmt.Println(strings.Repeat("-", 10))
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must provide at least one input file")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)

	optimizeCmd.Flags().StringP("target-dir", "t", "", "Output directory")
	optimizeCmd.Flags().BoolP("overwrite", "O", false, "Overwrite input files")
	optimizeCmd.Flags().BoolP("recursive", "r", false, "Search PDF files in subdirectories")
	optimizeCmd.Flags().StringP("password", "p", "", "File password")
	optimizeCmd.Flags().IntP("image-quality", "q", 100, "Optimized image quality")
}
