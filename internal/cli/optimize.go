/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const optimizeCmdDesc = `Optimize PDF files by optimizing structure, compression and image quality.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_optimized" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the optimized output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.

The quality of the images in the output files can be configured through
the --image-quality flag (default 90).
The resolution of the output images can be controlled using the --image-ppi flag.
Common pixels per inch values are 100 (screen), 150-300 (print), 600 (art). If
not specified, the PPI of the output images is 100.
`

var optimizeCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s optimize file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s optimize -O file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s optimize -O -r file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t out_dir file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t out_dir -r file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t out_dir -r -q 75 file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t out_dir -r -q 75 -P 100 file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s optimize -t out_dir -r -q 75 -P 100 -p pass file_1.pdf file_n.pdf dir_1 dir_n", appName),
)

// optimizeCmd represents the optimize command.
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
			imageQuality = 90
		}

		imagePPI, err := cmd.Flags().GetFloat64("image-ppi")
		if err != nil {
			imagePPI = 100
		}

		opts := &pdf.OptimizeOpts{
			ImageQuality: clampInt(imageQuality, 10, 100),
			ImagePPI:     imagePPI,
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
			res, err := pdf.Optimize(inputPath, outputPath, password, opts)
			if err != nil {
				printErr("Could not optimize input file: %s\n", err)
			}

			inSize := res.Original.Size
			outSize := res.Optimized.Size
			ratio := 100.0 - (float64(outSize) / float64(inSize) * 100.0)
			duration := float64(res.Duration) / float64(time.Millisecond)

			fmt.Printf("Original: %s\n", res.Original.Name)
			fmt.Printf("Original size: %d bytes\n", inSize)
			fmt.Printf("Optimized: %s\n", res.Optimized.Name)
			fmt.Printf("Optimized size: %d bytes\n", outSize)
			fmt.Printf("Compression ratio: %.2f%%\n", ratio)
			fmt.Printf("Processing time: %.2f ms\n", duration)
			fmt.Println("Status: success")
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

	optimizeCmd.Flags().StringP("target-dir", "t", "", "output directory")
	optimizeCmd.Flags().BoolP("overwrite", "O", false, "overwrite input files")
	optimizeCmd.Flags().BoolP("recursive", "r", false, "search PDF files in subdirectories")
	optimizeCmd.Flags().StringP("password", "p", "", "file password")
	optimizeCmd.Flags().IntP("image-quality", "q", 90, "output JPEG image quality")
	optimizeCmd.Flags().Float64P("image-ppi", "P", 100, "output images pixels per inch")
}
