/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unicreator "github.com/unidoc/unidoc/pdf/creator"
)

func Watermark(inputPath, outputPath, watermarkPath, password string, pages []int) error {
	// Read input file.
	r, pageCount, _, err := readPDF(inputPath, password)
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
		numPages, err := r.GetNumPages()
		if err != nil {
			return err
		}

		pages = createPageRange(numPages)
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
	return writeCreatorPDF(outputPath, c, false)
}
