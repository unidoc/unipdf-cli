/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func Split(inputPath, outputPath, password string, pageNums []int) error {
	// Read input file.
	r, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Add selected pages to the writer.
	w := unipdf.NewPdfWriter()
	if err = readerToWriter(r, &w, pageNums); err != nil {
		return err
	}

	// Write output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}
