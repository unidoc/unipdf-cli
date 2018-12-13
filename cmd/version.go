/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	unicommon "github.com/unidoc/unidoc/common"
)

const appVersion = "0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output version information and exit",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("unipdf %s\nunidoc %s\n", appVersion, unicommon.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
