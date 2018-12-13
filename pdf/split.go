/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"os"

	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func SplitPdf(inputPath, outputPath, password string, pageNums []int) error {
	// Read input file.
	r, pages, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Add selected pages to the writer.
	w := unipdf.NewPdfWriter()
	for _, pageNum := range pageNums {
		if pageNum < 0 || pageNum > pages {
			continue
		}

		page, err := r.GetPage(pageNum)
		if err != nil {
			return err
		}

		err = w.AddPage(page)
		if err != nil {
			return err
		}
	}

	// Add forms to the writer.
	if r.AcroForm != nil {
		w.SetForms(r.AcroForm)
	}

	// Create output file.
	of, err := os.Create(outputPath)
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
