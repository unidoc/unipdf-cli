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
)

const renderCmdDesc = `Renders the pages of the input file to image targets.

The rendered image files are saved in a ZIP archive at the location specified
by the --output-file parameter. If no output file is specified, the ZIP file
is saved in the same directory as the input file.

The command can be configured to render only the specified pages using
the --pages parameter.

An example of the pages parameter: 1-3,4,6-7
Pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be rendered, while page
number 5 is skipped.

The format of the rendered image files can be specified using
the --image-format flag (default jpeg).

Supported image formats:
  - jpeg (default)
  - png

The quality of the rendered image files can be configured through
the --image-quality flag (default 100). Only applies to JPEG images.
`

var renderCmdExample = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("%s render input_file.pdf", appName),
	fmt.Sprintf("%s render -o images.zip input_file.pdf", appName),
	fmt.Sprintf("%s render -o images.zip -P 1-3 input_file.pdf", appName),
	fmt.Sprintf("%s render -o images.zip -P 1-3 -p pass input_file.pdf", appName),
	fmt.Sprintf("%s render -o images.zip -P 1-3 -p pass -f jpeg -q 100 input_file.pdf", appName),
)

// renderCmd represents the render command.
var renderCmd = &cobra.Command{
	Use:                   "render [FLAG]... INPUT_FILE",
	Short:                 "Render PDF pages to images",
	Long:                  renderCmdDesc,
	Example:               renderCmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		inputPath := args[0]
		password, _ := cmd.Flags().GetString("password")
		outputPath, _ := cmd.Flags().GetString("output-file")

		// Parse page range.
		pageRange, _ := cmd.Flags().GetString("pages")

		pages, err := parsePageRange(pageRange)
		if err != nil {
			printUsageErr(cmd, "Invalid page range specified\n")
		}

		// Parse render options.
		imageFormat, _ := cmd.Flags().GetString("image-format")
		if _, ok := imageFormats[imageFormat]; !ok {
			imageFormat = "jpeg"
		}

		imageQuality, err := cmd.Flags().GetInt("image-quality")
		if err != nil {
			imageQuality = 100
		}

		opts := &pdf.RenderOpts{
			ImageFormat:  imageFormat,
			ImageQuality: imageQuality,
		}

		// Render file.
		outputPath, err = pdf.Render(inputPath, outputPath, password, pages, opts)
		if err != nil {
			printErr("Could not render input file: %s\n", err)
			return
		}

		fmt.Printf("File %s successfully rendered\n", inputPath)
		fmt.Printf("Output file saved to %s\n", outputPath)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must provide the input file")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringP("password", "p", "", "input file password")
	renderCmd.Flags().StringP("output-file", "o", "", "output file")
	renderCmd.Flags().StringP("pages", "P", "", "pages to render from the input file")
	renderCmd.Flags().StringP("image-format", "f", "jpeg", "format of the output images")
	renderCmd.Flags().IntP("image-quality", "q", 100, "quality of the output images")
}
