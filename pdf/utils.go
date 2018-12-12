/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"errors"

	unicommon "github.com/unidoc/unidoc/common"
	unicore "github.com/unidoc/unidoc/pdf/core"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

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
