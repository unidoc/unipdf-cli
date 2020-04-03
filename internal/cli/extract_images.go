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

	uniextractor "github.com/unidoc/unipdf/v3/extractor"
)

const extractImagesCmdDesc = `Extracts PDF images.

The images are extracted in a ZIP file and saved at the destination specified
by the --output-file parameter. If no output file is specified, the ZIP
archive is saved in the same directory as the input file.

The command can be configured to extract images only from the specified
pages using the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Images will only be extracted from pages 1,2,3 (1-3), 4 and 6,7 (6-7), while page
number 5 is skipped.
`

var extractImagesCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s extract images input_file.pdf", appName),
	fmt.Sprintf("%s extract images -o images.zip input_file.pdf", appName),
	fmt.Sprintf("%s extract images -P 1-3 -p pass -o images.zip input_file.pdf", appName),
	fmt.Sprintf("%s extract images -P 1-3 -p pass -o images.zip -S input_file.pdf", appName),
)

// extractImagesCmd represents the extract images command.
var extractImagesCmd = &cobra.Command{
	Use:                   "images [FLAG]... INPUT_FILE",
	Short:                 "Extract PDF images",
	Long:                  extractImagesCmdDesc,
	Example:               extractImagesCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")
		outputPath, _ := cmd.Flags().GetString("output-file")

		// Parse image extraction options.
		includeSM, _ := cmd.Flags().GetBool("include-inline-stencil-masks")

		extractOptions := &uniextractor.ImageExtractOptions{
			IncludeInlineStencilMasks: includeSM,
		}

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			printUsageErr(cmd, "Invalid page range specified\n")
		}

		// Extract images.
		outputPath, count, err := pdf.ExtractImages(
			inputPath,
			outputPath,
			password,
			pages,
			extractOptions,
		)
		if err != nil {
			printErr("Could not extract images: %s\n", err)
			return
		}

		if count == 0 {
			fmt.Printf("%s does not contain any images to extract\n", inputPath)
		} else {
			fmt.Printf("Images successfully extracted to %s\n", outputPath)
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
	extractCmd.AddCommand(extractImagesCmd)

	extractImagesCmd.Flags().StringP("password", "p", "", "input file password")
	extractImagesCmd.Flags().StringP("output-file", "o", "", "output file")
	extractImagesCmd.Flags().StringP("pages", "P", "", "pages to extract images from")
	extractImagesCmd.Flags().BoolP("include-inline-stencil-masks", "S", false, "include inline stencil masks")
}
