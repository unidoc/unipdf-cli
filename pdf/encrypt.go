/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unisecurity "github.com/unidoc/unidoc/pdf/core/security"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

type EncryptOpts struct {
	OwnerPassword string
	UserPassword  string
	Algorithm     unipdf.EncryptionAlgorithm
	Permissions   unisecurity.Permissions
}

func Encrypt(inputPath, outputPath string, opts *EncryptOpts) error {
	// Read input file.
	r, _, _, err := readPDF(inputPath, "")
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
	if err = writePDF(outputPath, &w, safe); err != nil {
		return err
	}

	return nil
}
