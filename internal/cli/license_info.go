/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
)

const licenseInfoCmdDesc = `Outputs information about the license key.`

var licenseInfoCmdExample = strings.Join([]string{
	fmt.Sprintf("%s license_info", appName),
}, "\n")

// licenseInfoCmd represents the license info command.
var licenseInfoCmd = &cobra.Command{
	Use:                   "license_info",
	Short:                 "Output license key information",
	Long:                  licenseInfoCmdDesc,
	Example:               licenseInfoCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(_ *cobra.Command, _ []string) {
		licenseKey := os.Getenv("UNIDOC_LICENSE_API_KEY")
		if licenseKey != "" {
			// To get your free API key for metered license, sign up on: https://cloud.unidoc.io
			// Make sure to be using UniOffice v1.9.0 or newer for Metered API key support
			lk := pdf.GetLicenseKey()
			fmt.Printf("License: %s\n", lk)

			// GetMeteredState freshly checks the state, contacting the licensing server.
			pdf.GetMeteredState()
			return
		}

		licensePath := os.Getenv("UNIDOC_LICENSE_FILE")
		if licensePath != "" {
			lk := pdf.GetLicenseKey()
			fmt.Printf("License: %s\n", lk)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(licenseInfoCmd)
}
