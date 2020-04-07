/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf-cli/pkg/pdf"
	unicommon "github.com/unidoc/unipdf/v3/common"
)

const appName = "unipdf"
const appVersion = "0.3.2"

const rootCmdDesc = ` is a CLI application for working with PDF files.
It supports the most common PDF operations. A full list of the supported
operations can be found below.

If you have a license for Unidoc, you can set it through the
UNIDOC_LICENSE_FILE and UNIDOC_LICENSE_CUSTOMER environment variables.

export UNIDOC_LICENSE_FILE="PATH_TO_LICENSE_FILE"
export UNIDOC_LICENSE_CUSTOMER="CUSTOMER_NAME"

By default, the application only displays error messages on command execution
failure. To change the verbosity of the output, set the UNIDOC_LOG_LEVEL
environment variable.

export UNIDOC_LOG_LEVEL="DEBUG"

Supported log levels: trace, debug, info, notice, warning, error (default)
`

var rootCmd = &cobra.Command{
	Use:  appName,
	Long: appName + rootCmdDesc,
}

// Execute represents the entry point of the application.
// The method parses the command line arguments and executes the appropriate
// action.
func Execute() {
	readEnv()

	if err := rootCmd.Execute(); err != nil {
		printErr("%s\n", err)
	}
}

func readEnv() {
	// Set license key.
	licensePath := os.Getenv("UNIDOC_LICENSE_FILE")
	licenseCustomer := os.Getenv("UNIDOC_LICENSE_CUSTOMER")

	if licensePath != "" {
		pdf.SetLicense(licensePath, licenseCustomer)
	}

	// Set log level.
	logLevel, err := parseLogLevel(os.Getenv("UNIDOC_LOG_LEVEL"))
	if err != nil {
		logLevel = unicommon.LogLevelError
	}

	pdf.SetLogLevel(logLevel)
}
