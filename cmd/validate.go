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

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:                   "validate [FLAG]... INPUT_FILE",
	Short:                 "Validate PDF files",
	Long:                  `A longer description that spans multiple lines and likely contains`,
	Example:               "this is the example",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		password, _ := cmd.Flags().GetString("password")

		info, err := pdf.GetPDFInfo(inputFile, password)
		if err != nil {
			fmt.Println("Could not validate input file")
			return
		}

		// Print basic PDF info
		fmt.Println("Info")
		fmt.Printf("Name: %s\n", inputFile)
		fmt.Printf("Size: %d bytes\n", info.Size)
		fmt.Printf("Pages: %d\n", info.Pages)
		fmt.Printf("PDF Version: %s\n", info.Version)

		if info.Encrypted {
			fmt.Printf("Encryption: encrypted with %s algorithm\n", info.EncryptionAlgo)
		} else {
			fmt.Println("Encryption: none")
		}

		// Print PDF objects
		fmt.Println("\nObjects")

		var malicious bool
		for key, val := range info.Objects {
			maliciousStr := ""
			if key == "JavaScript" || key == "Flash" || key == "Video" {
				maliciousStr = " (potentially malicious)"
				malicious = true
			}

			fmt.Printf("%s objects: %d%s\n", key, val, maliciousStr)
		}

		if malicious {
			fmt.Println("\nFile contains potentially malicious objects!")
		} else {
			fmt.Println("\nFile is safe")
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must provide the input file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("password", "p", "", "PDF file password")
}
