/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unipdf "github.com/unidoc/unipdf/v3/model"
)

// Split extracts the provided page list from PDF file specified by the
// inputPath parameter and saves the resulting file at the location
// specified by the outputPath parameter. A password can be passed in for
// encrypted input files.
// If the pages parameter is nil or an empty slice, all the pages of the input
// file are copied to the output file.
func Split(inputPath, outputPath, password string, pages []int) error {
	// Read input file.
	r, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Add selected pages to the writer.
	w := unipdf.NewPdfWriter()
	if err = readerToWriter(r, &w, pages); err != nil {
		return err
	}

	// Write output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}
