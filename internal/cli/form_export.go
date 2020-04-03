/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const formExportCmdDesc = `Export JSON representation of form fields.

By default, the resulting JSON content is printed to STDOUT. The output can be
saved to a file by using the --output-file flag (see the Examples section).

The exported JSON template can be used to fill PDF forms using the
"form fill" command.
`

var formExportCmdExample = fmt.Sprintf("%s\n%s\n%s\n",
	fmt.Sprintf("%s form export in_file.pdf", appName),
	fmt.Sprintf("%s form export in_file.pdf > out_file.json", appName),
	fmt.Sprintf("%s form export -o out_file.json in_file.pdf", appName),
)

// formExportCmd represents the form export command.
var formExportCmd = &cobra.Command{
	Use:                   "export [FLAG]... INPUT_FILE",
	Short:                 "Export form fields as JSON",
	Long:                  formExportCmdDesc,
	Example:               formExportCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		outputPath, _ := cmd.Flags().GetString("output-file")

		// Export form fields.
		json, err := pdf.FormExport(inputPath)
		if err != nil {
			printErr("Could not export form fields: %s\n", err)
			return
		}
		if json == "" {
			fmt.Println("Could not find any form fields to export.")
			return
		}

		// Write exported data.
		if outputPath == "" {
			fmt.Println(json)
			return
		}

		err = ioutil.WriteFile(outputPath, []byte(json), os.ModePerm)
		if err != nil {
			printErr("Could not export form fields: %s\n", err)
		}

		fmt.Printf("Form fields successfully exported from %s\n", inputPath)
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
	formCmd.AddCommand(formExportCmd)

	formExportCmd.Flags().StringP("output-file", "o", "", "output file")
}
