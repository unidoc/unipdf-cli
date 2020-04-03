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

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const formFlattenCmdDesc = `Flatten PDF file form annotations.

The flattening process makes the form fields of the output files read-only by
appending the form field annotation XObject Form data to the page content
stream, thus making it part of the page contents.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_flattened" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the flattened output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.
`

var formFlattenCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s form flatten file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s form flatten -O file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s form flatten -O -r file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form flatten -t out_dir file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form flatten -t out_dir -r file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form flatten -t out_dir -r -p pass file_1.pdf file_n.pdf dir_1 dir_n", appName),
)

// formFlattenCmd represents the form flatten command.
var formFlattenCmd = &cobra.Command{
	Use:                   "flatten [FLAG]... INPUT_FILES...",
	Short:                 "Flatten form annotations",
	Long:                  formFlattenCmdDesc,
	Example:               formFlattenCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags.
		outputDir, _ := cmd.Flags().GetString("target-dir")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		recursive, _ := cmd.Flags().GetBool("recursive")
		password, _ := cmd.Flags().GetString("password")

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

		// Flatten PDF files form annotations.
		for _, inputPath := range inputPaths {
			fmt.Printf("Flattening %s\n", inputPath)

			// Generate output path.
			outputPath := generateOutputPath(inputPath, outputDir, "flattened", overwrite)

			// Flatten input file form fields.
			err := pdf.FormFlatten(inputPath, outputPath, password)
			if err != nil {
				printErr("Could not flatten input file form annotations: %s\n", err)
			}

			fmt.Printf("Original: %s\n", inputPath)
			fmt.Printf("Flattened: %s\n", outputPath)
			fmt.Println("Status: success")
			fmt.Println(strings.Repeat("-", 10))
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must provide the at least on input file or directory")
		}

		return nil
	},
}

func init() {
	formCmd.AddCommand(formFlattenCmd)

	formFlattenCmd.Flags().StringP("target-dir", "t", "", "output directory")
	formFlattenCmd.Flags().BoolP("overwrite", "O", false, "overwrite input files")
	formFlattenCmd.Flags().BoolP("recursive", "r", false, "search PDF files in subdirectories")
	formFlattenCmd.Flags().StringP("password", "p", "", "input file password")
}
