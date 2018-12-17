/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"github.com/spf13/cobra"
)

const appName = "unipdf"
const appVersion = "0.1"

const rootCmdDesc = ` is a CLI application for working with PDF files.
It supports the most common PDF operations. A full list of the supported
operations can be found below.`

var rootCmd = &cobra.Command{
	Use:  appName,
	Long: appName + rootCmdDesc,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printErr("%s\n", err)
	}
}
