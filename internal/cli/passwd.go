/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const passwdCmdDesc = `Change owner and user passwords of PDF files.`

var passwdCmdExample = fmt.Sprintf("%s\n%s\n%s\n",
	fmt.Sprintf("%s passwd -p pass input_file.pdf new_owner_pass", appName),
	fmt.Sprintf("%s passwd -p pass -o output_file.pdf input_file.pdf new_owner_pass", appName),
	fmt.Sprintf("%s passwd -p pass -o output_file.pdf input_file.pdf new_owner_pass new_user_pass", appName),
)

// passwdCmd represents the passwd command.
var passwdCmd = &cobra.Command{
	Use:                   "passwd [FLAG]... INPUT_FILE NEW_OWNER_PASSWORD [NEW_USER_PASSWORD]",
	Short:                 "Change PDF passwords",
	Long:                  passwdCmdDesc,
	Example:               passwdCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		newOwnerPassword := args[1]
		ownerPassword, _ := cmd.Flags().GetString("password")

		newUserPassword := ""
		if len(args) > 2 {
			newUserPassword = args[2]
		}

		// Parse output file.
		outputPath, _ := cmd.Flags().GetString("output-file")
		if outputPath == "" {
			outputPath = inputPath
		}

		// Change input file password.
		err := pdf.Passwd(inputPath, outputPath, ownerPassword, newOwnerPassword, newUserPassword)
		if err != nil {
			printErr("Could not change input file password: %s\n", err)
		}

		fmt.Printf("Password successfully changed\n")
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide the input file and the new owner password")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(passwdCmd)

	passwdCmd.Flags().StringP("output-file", "o", "", "output file")
	passwdCmd.Flags().StringP("password", "p", "", "input file password")
}
