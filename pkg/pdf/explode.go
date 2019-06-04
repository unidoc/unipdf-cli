/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	unipdf "github.com/unidoc/unipdf/v3/model"
)

// Explode splits the PDF file specified by the inputPath parameter into single
// page PDF files. The extracted collection of PDF files is saved as a ZIP
// archive at the location specified by the outputPath parameter.
// A password can be passed in, if the input file is encrypted.
// If the pages parameter is nil or an empty slice, all pages are extracted.
func Explode(inputPath, outputPath, password string, pages []int) (string, error) {
	dir, inputFile := filepath.Split(inputPath)
	// Use input file directory if no output path is specified.
	inputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	if outputPath == "" {
		outputPath = filepath.Join(dir, inputFile+".zip")
	}

	// Read input file.
	r, pageCount, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return "", err
	}

	// Prepare output archive.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	// Extract pages.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	zw := zip.NewWriter(outputFile)
	for _, numPage := range pages {
		w := unipdf.NewPdfWriter()
		if err := readerToWriter(r, &w, []int{numPage}); err != nil {
			return "", err
		}

		// Add page to zip file.
		file, err := zw.Create(fmt.Sprintf("%s_%d.pdf", inputFile, numPage))
		if err != nil {
			return "", err
		}

		if err = w.Write(file); err != nil {
			return "", err
		}
	}

	return outputPath, zw.Close()
}
