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

const formFDFMergeCmdDesc = `Fill form fields from FDF file.

The field values specified in the FDF file template are used to fill the form
fields in the input PDF files. In addition, the output file form fields can be
flattened by using the --flatten flag. The flattening process makes the form
fields of the output files read-only by appending the form field annotation
XObject Form data to the page content stream, thus making it part of the page
contents.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_filled" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the filled output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.
`

var formFDFMergeCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s form fdfmerge fields.fdf file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s form fdfmerge -O fields.fdf file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s form fdfmerge -O -r -f fields.fdf file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form fdfmerge -t out_dir fields.fdf file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form fdfmerge -t out_dir -r fields.fdf file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form fdfmerge -t out_dir -r -p pass fields.fdf file_1.pdf file_n.pdf dir_1 dir_n", appName),
)

// formFDFMergeCmd represents the form fdfmerge command.
var formFDFMergeCmd = &cobra.Command{
	Use:                   "fdfmerge [FLAG]... FDF_FILE INPUT_FILES...",
	Short:                 "Fill form fields from FDF file",
	Long:                  formFDFMergeCmdDesc,
	Example:               formFDFMergeCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input flags.
		outputDir, _ := cmd.Flags().GetString("target-dir")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		recursive, _ := cmd.Flags().GetBool("recursive")
		password, _ := cmd.Flags().GetString("password")
		flatten, _ := cmd.Flags().GetBool("flatten")

		// Parse input parameters.
		fdfPath := args[0]

		inputPaths, err := parseInputPaths(args[1:], recursive, pdfMatcher)
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

		// Fill form fields.
		for _, inputPath := range inputPaths {
			fmt.Printf("Filling form values for %s\n", inputPath)

			// Generate output path.
			outputPath := generateOutputPath(inputPath, outputDir, "filled", overwrite)

			// Fill input file form fields.
			err := pdf.FormFillFDF(inputPath, fdfPath, outputPath, password, flatten)
			if err != nil {
				printErr("Could not fill form fields: %s\n", err)
			}

			fmt.Printf("Original: %s\n", inputPath)
			fmt.Printf("Filled: %s\n", outputPath)
			fmt.Println("Status: success")
			fmt.Println(strings.Repeat("-", 10))
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide the FDF file and at least one input file")
		}

		return nil
	},
}

func init() {
	formCmd.AddCommand(formFDFMergeCmd)

	formFDFMergeCmd.Flags().StringP("target-dir", "t", "", "output directory")
	formFDFMergeCmd.Flags().BoolP("overwrite", "O", false, "overwrite input files")
	formFDFMergeCmd.Flags().BoolP("recursive", "r", false, "search PDF files in subdirectories")
	formFDFMergeCmd.Flags().StringP("password", "p", "", "input file password")
	formFDFMergeCmd.Flags().BoolP("flatten", "f", false, "flatten form annotations")
}
