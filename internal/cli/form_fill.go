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

const formFillCmdDesc = `Fill form fields from JSON file.

The field values specified in the JSON file template are used to fill the form
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

The "form export" command can be used to generate the JSON form fields template
for a PDF file.
`

var formFillCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s form fill fields.json file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s form fill -O fields.json file_1.pdf file_n.pdf", appName),
	fmt.Sprintf("%s form fill -O -r -f fields.json file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form fill -t out_dir fields.json file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form fill -t out_dir -r fields.json file_1.pdf file_n.pdf dir_1 dir_n", appName),
	fmt.Sprintf("%s form fill -t out_dir -r -p pass fields.json file_1.pdf file_n.pdf dir_1 dir_n", appName),
)

// formFillCmd represents the form fill command.
var formFillCmd = &cobra.Command{
	Use:                   "fill [FLAG]... JSON_FILE INPUT_FILES...",
	Short:                 "Fill form fields from JSON file",
	Long:                  formFillCmdDesc,
	Example:               formFillCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input flags.
		outputDir, _ := cmd.Flags().GetString("target-dir")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		recursive, _ := cmd.Flags().GetBool("recursive")
		password, _ := cmd.Flags().GetString("password")
		flatten, _ := cmd.Flags().GetBool("flatten")

		// Parse input parameters.
		jsonPath := args[0]

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
			err := pdf.FormFillJSON(inputPath, jsonPath, outputPath, password, flatten)
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
			return errors.New("must provide the JSON file and at least one input file")
		}

		return nil
	},
}

func init() {
	formCmd.AddCommand(formFillCmd)

	formFillCmd.Flags().StringP("target-dir", "t", "", "output directory")
	formFillCmd.Flags().BoolP("overwrite", "O", false, "overwrite input files")
	formFillCmd.Flags().BoolP("recursive", "r", false, "search PDF files in subdirectories")
	formFillCmd.Flags().StringP("password", "p", "", "input file password")
	formFillCmd.Flags().BoolP("flatten", "f", false, "flatten form annotations")
}
