/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import unipdf "github.com/unidoc/unidoc/pdf/model"

func Decrypt(inputPath, outputPath, password string) error {
	// Read input file.
	r, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Copy input file contents.
	w := unipdf.NewPdfWriter()
	if err := readerToWriter(r, &w, nil); err != nil {
		return err
	}

	// Save output file.
	if err = writePDF(outputPath, &w, false); err != nil {
		return err
	}

	return nil
}
