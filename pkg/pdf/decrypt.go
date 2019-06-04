/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import unipdf "github.com/unidoc/unipdf/v3/model"

// Decrypt decrypts the PDF file specified by the inputPath parameter,
// using the specified password and saves the result to the destination
// specified by the outputPath parameter.
func Decrypt(inputPath, outputPath, password string) error {
	// Read input file.
	r, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Copy input file contents.
	w := unipdf.NewPdfWriter()
	if err := readerToWriter(r, &w, nil); err != nil {
		return err
	}

	// Save output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}
