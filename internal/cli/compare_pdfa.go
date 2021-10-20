package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/unidoc/unipdf-cli/pkg/pdf"

	"github.com/unidoc/unipdf-cli/internal/verapdf"
)

const comparePdfACmdDesc = `Is a command that compare UniPdf PDF/A validation implementation with the 'VeraPDF' XML report.
The VeraPDF report could have different form depending on the way it was obtained. 
Define the 'report-type' flag with one of the following values:
	- report obtained from online tool - 'online' - https://demo.verapdf.org/
	- report obtained from verapdf CLI - 'cli' - verapdf --format xml
	- report obtained from verapdf-gui - 'gui' - verapdf-gui or verapdf
The command requires two arguments:
	- VeraPDF XML report
	- Pdf file for which given report was generated.
Currently only two VeraPDF profiles are supported:
	- PDF/A-1A
	- PDF/A-1B
`

var comparePdfAReportType verapdf.ReportType

var comparePdfACmdExample = strings.Join([]string{
	fmt.Sprintf("%s compare-pdfa --report-type online report.xml document.pdf", appName),
	fmt.Sprintf("%s compare-pdfa -r cli report.xml document.pdf", appName),
	fmt.Sprintf("%s compare-pdfa -r gui report.xml document.pdf", appName),
}, "\n",
)

// comparePdfACmd represents the decrypt command.
var comparePdfACmd = &cobra.Command{
	Use:                   "compare-pdfa [FLAG]... INPUT_VERAPDF_XML_REPORT_FILE INPUT_PDF_DOCUMENT_FILE",
	Short:                 "Compares UniPdf PDF/A validation with VeraPDF report",
	Long:                  comparePdfACmdDesc,
	Example:               comparePdfACmdExample,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input parameters.
		if len(args) != 2 {
			printUsageErr(cmd, "compare-pdfa requires exactly two arguments: verapdf XML report and matching Pdf document\n")
		}

		if !comparePdfAReportType.Valid() {
			printUsageErr(cmd, "report type undefined\n")
		}
		xmlReport := args[0]
		pdfDoc := args[1]

		rf, err := os.Open(xmlReport)
		if err != nil {
			printErr("opening xml report file failed: %v\n", err)
		}
		defer rf.Close()

		df, err := os.Open(pdfDoc)
		if err != nil {
			printErr("opening pdf document file failed: %v", err)
		}
		defer df.Close()

		result, err := pdf.ComparePdfARules(df, rf, comparePdfAReportType)
		if err != nil {
			printErr("verifying PDF/A failed: %v\n", err)
		}

		if result.Passed {
			fmt.Println("Results of the VeraPdf report matches with the UniPdf PDF/A validation.")
			os.Exit(1)
		}

		if len(result.MismatchedRules) != 0 {
			sb := strings.Builder{}
			sb.WriteString("Mismatched rules:\n")
			var i int
			for k, v := range result.MismatchedRules {
				sb.WriteRune('\t')
				sb.WriteString(k)
				sb.WriteRune(':')
				sb.WriteRune('\t')
				sb.WriteString(v)
				if i != len(result.MismatchedRules)-1 {
					sb.WriteRune(',')
				}
				sb.WriteRune('\n')
				i++
			}
			fmt.Print(sb.String())
		}
		if result.IsExpectedValid {
			fmt.Printf("VeraPdf marks document as valid\n")
		} else {
			fmt.Printf("VeraPdf marks document as invalid\n")
		}
		if result.IsUniPdfValid {
			fmt.Printf("UniPDF/A marks document as valid\n")
		} else {
			fmt.Printf("UniPDF/A marks document as invalid\n")
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("compare-pdfa requires exactly two arguments: verapdf XML report and matching Pdf document")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(comparePdfACmd)
	comparePdfACmd.Flags().VarP(&comparePdfAReportType, "report-type", "r", "verapdf type of the report")
	if err := comparePdfACmd.MarkFlagRequired("report-type"); err != nil {
		fmt.Printf("compare-pdfa marking flag required: 'report-type' failed: %v\n", err)
		os.Exit(1)
	}
}
