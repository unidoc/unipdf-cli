/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unidoc/unicli/pdf"
	unicommon "github.com/unidoc/unidoc/common"
)

const appName = "unicli"
const appVersion = "0.1"

const rootCmdDesc = ` is a CLI application for working with PDF files.
It supports the most common PDF operations. A full list of the supported
operations can be found below.

If you have a license for Unidoc, you can set it through the
UNIDOC_LICENSE_FILE and UNIDOC_LICENSE_CUSTOMER environment variables.

EXPORT UNIDOC_LICENSE_FILE="PATH_TO_LICENSE_FILE"
EXPORT UNIDOC_LICENSE_CUSTOMER="CUSTOMER_NAME"

By default, the application only displays error messages on command execution
failure. To change the verbosity of the output, set the UNIDOC_LOG_LEVEL
environment variable.

EXPORT UNIDOC_LOG_LEVEL="DEBUG"

Supported log levels: trace, debug, info, notice, warning, error (default)
`

var rootCmd = &cobra.Command{
	Use:  appName,
	Long: appName + rootCmdDesc,
}

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
