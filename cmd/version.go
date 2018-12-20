/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unidoc/unicli/pdf"
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

		fmt.Printf("%s %s\n", appName, appVersion)
		fmt.Printf("unidoc %s\n", version.Lib)
		fmt.Printf("\nLicense info\n%s", version.License)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
