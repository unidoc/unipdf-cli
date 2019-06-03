/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unisecurity "github.com/unidoc/unipdf/v3/core/security"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

// EncryptOpts contains settings for encrypting a PDF file.
type EncryptOpts struct {
	// OwnerPassword represents the owner password used to encrypt the file.
	OwnerPassword string

	// UserPassword represents the user password used to encrypt the file.
	UserPassword string

	// Algorithm represents the encryption algorithm used to encrypt the file.
	Algorithm unipdf.EncryptionAlgorithm

	// Permissions specifies the operations the user can execute on
	// the encrypted PDF file.
	Permissions unisecurity.Permissions
}

// Encrypt encrypts the PDF file specified by the inputPath parameter,
// using the specified options and saves the result at the location
// specified by the outputPath parameter.
func Encrypt(inputPath, outputPath string, opts *EncryptOpts) error {
	// Read input file.
	r, _, _, _, err := readPDF(inputPath, "")
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
		Algorithm:   opts.Algorithm,
		Permissions: opts.Permissions,
	}

	err = w.Encrypt([]byte(opts.UserPassword), []byte(opts.OwnerPassword), encryptOpts)
	if err != nil {
		return err
	}

	// Save output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}
