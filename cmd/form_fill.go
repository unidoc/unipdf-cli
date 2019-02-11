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

const formFillCmdDesc = `Fill form fields from JSON file.
`

var formFillCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s form fill in_file.pdf fields.json", appName),
	fmt.Sprintf("%s form fill -o out_file.pdf in_file.pdf fields.json", appName),
	fmt.Sprintf("%s form fill -o out_file.pdf -f in_file.pdf fields.json", appName),
	fmt.Sprintf("%s form fill -o out_file.pdf -f -p pass in_file.pdf fields.json", appName),
)

// formFillCmd represents the form export command
var formFillCmd = &cobra.Command{
	Use:                   "fill [FLAG]... JSON_FILE INPUT_FILES...",
	Short:                 "Fill form fields from JSON file",
	Long:                  formFillCmdDesc,
	Example:               formFillCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		jsonPath := args[1]

		// Parse input flags.
		flatten, _ := cmd.Flags().GetBool("flatten")
		password, _ := cmd.Flags().GetString("password")

		// Parse output path.
		outputPath, _ := cmd.Flags().GetString("output-file")
		if outputPath == "" {
			outputPath = inputPath
		}

		// Fill form fields.
		err := pdf.FormFill(inputPath, jsonPath, outputPath, password, flatten)
		if err != nil {
			printErr("Could not fill form fields: %s\n", err)
		}

		//fmt.Printf("Form fields successfully exported from %s\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide the input PDF and the JSON fields files")
		}

		return nil
	},
}

func init() {
	formCmd.AddCommand(formFillCmd)

	formFillCmd.Flags().StringP("password", "p", "", "input file password")
	formFillCmd.Flags().StringP("output-file", "o", "", "output file")
	formFillCmd.Flags().BoolP("flatten", "f", false, "flatten form annotations")
}
