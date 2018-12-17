/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/pdf"
)

var versionCmdExample = fmt.Sprintf("%s\n",
	fmt.Sprintf("%s version", appName),
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:                   "version",
	Short:                 "Output version information and exit",
	Example:               versionCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		version := pdf.Version()

		fmt.Printf("unipdf %s\n", version.App)
		fmt.Printf("unidoc %s\n", version.Lib)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
