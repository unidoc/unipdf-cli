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

const infoCmdDesc = `Outputs information about the input file.
Also provides basic validation.
`

var infoCmdExample = fmt.Sprintf("%s\n%s\n",
	fmt.Sprintf("%s info input_file.pdf", appName),
	fmt.Sprintf("%s info -p pass input_file.pdf", appName),
)

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:                   "info [FLAG]... INPUT_FILE",
	Short:                 "Output PDF information",
	Long:                  infoCmdDesc,
	Example:               infoCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		password, _ := cmd.Flags().GetString("password")

		info, err := pdf.Info(inputFile, password)
		if err != nil {
			printErr("Could not retrieve input file information: %s\n", err)
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
			return errors.New("must provide the input file")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringP("password", "p", "", "input file password")
}
