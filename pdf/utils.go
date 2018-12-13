/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"errors"
	"os"

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

func writePDF(filepath string, w *unipdf.PdfWriter, safe bool) error {
	// Create output file.
	of, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer of.Close()

	// Write output file.
	err = w.Write(of)
	if err != nil {
		return err
	}

	return nil
}

func writeCreatorPDF(filepath string, c *unicreator.Creator, safe bool) error {
	// Create output file.
	of, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer of.Close()

	// Write output file.
	err = c.Write(of)
	if err != nil {
		return err
	}

	return nil
}

func readerToWriter(r *unipdf.PdfReader, w *unipdf.PdfWriter, pages []int) error {
	if r == nil {
		return errors.New("Source PDF cannot be null")
	}
	if w == nil {
		return errors.New("Destination PDF cannot be null")
	}

	// Add pages.
	if len(pages) == 0 {
		numPages, err := r.GetNumPages()
		if err != nil {
			return err
		}

		pages = createPageRange(numPages)
	}

	for _, pageNum := range pages {
		page, err := r.GetPage(pageNum)
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

func readerToCreator(r *unipdf.PdfReader, w *unicreator.Creator, pages []int) error {
	if r == nil {
		return errors.New("Source PDF cannot be null")
	}
	if w == nil {
		return errors.New("Destination PDF cannot be null")
	}

	// Add pages.
	if len(pages) == 0 {
		numPages, err := r.GetNumPages()
		if err != nil {
			return err
		}

		pages = createPageRange(numPages)
	}

	for _, pageNum := range pages {
		page, err := r.GetPage(pageNum)
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

func createPageRange(count int) []int {
	if count <= 0 {
		return []int{}
	}

	var pages []int
	for i := 0; i < count; i++ {
		pages = append(pages, i+1)
	}

	return pages
}
