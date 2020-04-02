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

const encryptCmdDesc = `Encrypts the input file using the specified owner password.

The algorithm used for the file encryption is configurable.

Supported encryption algorithms:
  - rc4 (default)
  - aes128
  - aes256

A user password along with a set of permissions can also be specified.

Supported user permissions:
  - all (default)
  - none
  - print-low-res
  - print-high-res
  - modify
  - extract
  - extract-graphics
  - annotate
  - fill-forms
  - rotate
`

var encryptCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s encrypt input_file.pdf owner_pass", appName),
	fmt.Sprintf("%s encrypt input_file.pdf owner_pass user_pass", appName),
	fmt.Sprintf("%s encrypt -o output_file.pdf -m aes256 input_file.pdf owner_pass user_pass", appName),
	fmt.Sprintf("%s encrypt -o output_file.pdf -P none -m aes256 input_file.pdf owner_pass user_pass", appName),
	fmt.Sprintf("%s encrypt -o output_file.pdf -P modify,annotate -m aes256 input_file.pdf owner_pass user_pass", appName),
)

// encryptCmd represents the encrypt command.
var encryptCmd = &cobra.Command{
	Use:                   "encrypt [FLAG]... INPUT_FILE OWNER_PASSWORD [USER_PASSWORD]",
	Short:                 "Encrypt PDF files",
	Long:                  encryptCmdDesc,
	Example:               encryptCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		ownerPassword := args[1]

		// Parse user password.
		var userPassword string
		if len(args) > 2 {
			userPassword = args[2]
		}

		// Parse output file.
		outputPath, _ := cmd.Flags().GetString("output-file")
		if outputPath == "" {
			outputPath = inputPath
		}

		// Parse encryption mode.
		mode, _ := cmd.Flags().GetString("mode")

		algorithm, err := parseEncryptionMode(mode)
		if err != nil {
			printUsageErr(cmd, "Invalid encryption mode\n")
		}

		// Parse user permissions.
		permList, _ := cmd.Flags().GetString("perms")

		perms, err := parsePermissionList(permList)
		if err != nil {
			printUsageErr(cmd, "Invalid user permission values\n")
		}

		opts := &pdf.EncryptOpts{
			OwnerPassword: ownerPassword,
			UserPassword:  userPassword,
			Algorithm:     algorithm,
			Permissions:   perms,
		}

		// Encrypt file.
		if err := pdf.Encrypt(inputPath, outputPath, opts); err != nil {
			printErr("Could not encrypt file: %s\n", err)
		}

		fmt.Printf("File %s successfully encrypted\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide the input file and the owner password")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("output-file", "o", "", "output file")
	encryptCmd.Flags().StringP("perms", "P", "all", "user permissions")
	encryptCmd.Flags().StringP("mode", "m", "rc4", "algorithm to use for encrypting the file")
}
