/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"errors"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	unicore "github.com/unidoc/unidoc/pdf/core"
	unicreator "github.com/unidoc/unidoc/pdf/creator"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func readPDF(filepath, password string) (*unipdf.PdfReader, int, bool, error) {
	// Open input file.
	f, err := os.Open(filepath)
	if err != nil {
		return nil, 0, false, err
	}
	defer f.Close()

	// Read input file.
	r, err := unipdf.NewPdfReader(f)
	if err != nil {
		return nil, 0, false, err
	}

	// Check if file is encrypted.
	encrypted, err := r.IsEncrypted()
	if err != nil {
		return nil, 0, false, err
	}

	// Decrypt using the specified password, if necessary.
	if encrypted {
		auth, err := r.Decrypt([]byte(password))
		if err != nil {
			return nil, 0, false, err
		}
		if !auth {
			return nil, 0, false, errors.New("Unable to decrypt the file with the specified password")
		}
	}

	// Get number of pages.
	pages, err := r.GetNumPages()
	if err != nil {
		return nil, 0, false, err
	}

	return r, pages, encrypted, nil
}

func readerToWriter(r *unipdf.PdfReader, w *unipdf.PdfWriter) error {
	if r == nil {
		return errors.New("Source PDF cannot be null")
	}
	if w == nil {
		return errors.New("Destination PDF cannot be null")
	}

	// Add pages.
	numPages, err := r.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		page, err := r.GetPage(i + 1)
		if err != nil {
			return err
		}

		if err = w.AddPage(page); err != nil {
			return err
		}
	}

	// Add forms.
	if r.AcroForm != nil {
		w.SetForms(r.AcroForm)
	}

	return nil
}

func readerToCreator(r *unipdf.PdfReader, w *unicreator.Creator) error {
	if r == nil {
		return errors.New("Source PDF cannot be null")
	}
	if w == nil {
		return errors.New("Destination PDF cannot be null")
	}

	// Add pages.
	numPages, err := r.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		page, err := r.GetPage(i + 1)
		if err != nil {
			return err
		}

		if err = w.AddPage(page); err != nil {
			return err
		}
	}

	// Add forms.
	if r.AcroForm != nil {
		w.SetForms(r.AcroForm)
	}

	return nil
}

func getDict(obj unicore.PdfObject) *unicore.PdfObjectDictionary {
	if obj == nil {
		return nil
	}

	obj = unicore.TraceToDirectObject(obj)
	dict, ok := obj.(*unicore.PdfObjectDictionary)
	if !ok {
		unicommon.Log.Debug("Error type check error (got %T)", obj)
		return nil
	}

	return dict
}
