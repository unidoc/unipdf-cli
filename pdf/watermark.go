/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unicreator "github.com/unidoc/unipdf/v3/creator"
)

// Watermark adds the watermark image specified by the watermarkPath parameter
// to the pages of the PDF file specified by the inputPath parameter.
// A password can be passed in for encrypted input files.
// The resulting file is saved at the location specified by the outputPath
// parameter.
// Also, a list of pages to add watermark to can be passed in. Every page that
// is not included in the pages slice is left intact.
// If the pages parameter is nil or an empty slice, all the pages of the input
// file are watermarked.
func Watermark(inputPath, outputPath, watermarkPath, password string, pages []int) error {
	// Read input file.
	r, pageCount, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Open watermark image.
	c := unicreator.New()

	watermark, err := c.NewImageFromFile(watermarkPath)
	if err != nil {
		return err
	}

	// Add pages.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	for i := 0; i < pageCount; i++ {
		numPage := i + 1

		page, err := r.GetPage(numPage)
		if err != nil {
			return err
		}

		var hasWatermark bool
		for _, page := range pages {
			if page == numPage {
				hasWatermark = true
				break
			}
		}

		if err = c.AddPage(page); err != nil {
			return err
		}

		if !hasWatermark {
			continue
		}

		watermark.ScaleToWidth(c.Context().PageWidth)
		watermark.SetPos(0, (c.Context().PageHeight-watermark.Height())/2)
		watermark.SetOpacity(0.5)

		if err = c.Draw(watermark); err != nil {
			return err
		}
	}

	// Add forms.
	if r.AcroForm != nil {
		c.SetForms(r.AcroForm)
	}

	// Write output file.
	safe := inputPath == outputPath
	return writeCreatorPDF(outputPath, c, safe)
}
