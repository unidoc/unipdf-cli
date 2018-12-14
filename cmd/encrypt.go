/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/pdf"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:                   "encrypt [FLAG]... INPUT_FILE OWNER_PASSWORD",
	Short:                 "Encrypt PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		ownerPassword := args[1]
		userPassword, _ := cmd.Flags().GetString("user-password")

		// Parse output file.
		outputFile, _ := cmd.Flags().GetString("output-file")
		if outputFile == "" {
			outputFile = inputFile
		}

		// Parse encryption mode
		mode, _ := cmd.Flags().GetString("mode")

		algorithm, err := parseEncryptionMode(mode)
		if err != nil {
			fmt.Println("Invalid encryption mode")
			return
		}

		// Parse user permissions
		permList, _ := cmd.Flags().GetString("perms")

		perms, err := parsePermissionList(permList)
		if err != nil {
			fmt.Println("Invalid user permissions")
			return
		}

		opts := &pdf.EncryptOpts{
			OwnerPassword: ownerPassword,
			UserPassword:  userPassword,
			Algorithm:     algorithm,
			Permissions:   perms,
		}

		// Encrypt input file.
		if err := pdf.Encrypt(inputFile, outputFile, opts); err != nil {
			fmt.Println("Could not encrypt input file")
			spew.Dump(err)
			return
		}

		fmt.Println("File successfully encrypted")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Must provide the input file and the owner password\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("user-password", "p", "", "PDF file password")
	encryptCmd.Flags().StringP("output-file", "o", "", "Output file")
	encryptCmd.Flags().StringP("perms", "P", "all", "User permissions")
	encryptCmd.Flags().StringP("mode", "m", "rc4", "Algorithm to use for encrypting the file")
}
