/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/pdf"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:                   "decrypt [FLAG]... INPUT_FILE",
	Short:                 "Decrypt PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		password, _ := cmd.Flags().GetString("password")

		// Parse output file.
		outputFile, _ := cmd.Flags().GetString("output-file")
		if outputFile == "" {
			outputFile = inputFile
		}

		if err := pdf.Decrypt(inputFile, outputFile, password); err != nil {
			fmt.Println("Could not decrypt input file")
			return
		}

		fmt.Println("Successfully decrypted input file")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must provide the input file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringP("password", "p", "", "PDF file password")
	decryptCmd.Flags().StringP("output-file", "o", "", "Output file")
}
