/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"github.com/spf13/cobra"
)

const formCmdDesc = `PDF form operations.`

// formCmd represents the form command.
var formCmd = &cobra.Command{
	Use:   "form [FLAG]... COMMAND",
	Short: "PDF form operations",
	Long:  formCmdDesc,
}

func init() {
	rootCmd.AddCommand(formCmd)
}
