/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"github.com/spf13/cobra"
)

const extractCmdDesc = `Extract PDF resources.`

// extractCmd represents the extract command.
var extractCmd = &cobra.Command{
	Use:   "extract [FLAG]... COMMAND",
	Short: "Extract PDF resources",
	Long:  extractCmdDesc,
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
