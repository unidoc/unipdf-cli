/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unipdf "github.com/unidoc/unipdf/v3/model"
)

// Passwd changes the owner and user password of an encrypted PDF file.
// The resulting PDF file is saved at the location specified by the outputPath
// parameter.
func Passwd(inputPath, outputPath, ownerPassword, newOwnerPassword, newUserPassword string) error {
	// Read input file.
	r, _, _, perms, err := readPDF(inputPath, ownerPassword)
	if err != nil {
		return err
	}

	// Copy input file contents.
	w := unipdf.NewPdfWriter()
	if err := readerToWriter(r, &w, nil); err != nil {
		return err
	}

	// Encrypt output file.
	encryptOpts := &unipdf.EncryptOptions{
		Permissions: perms,
	}

	err = w.Encrypt([]byte(newUserPassword), []byte(newOwnerPassword), encryptOpts)
	if err != nil {
		return err
	}

	// Save output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}
