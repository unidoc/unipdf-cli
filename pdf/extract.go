/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	uniextractor "github.com/unidoc/unidoc/pdf/extractor"
)

// ExtractText returns all text content from the PDF file specified by the
// inputPath parameter. A password can be specified for encrypted PDF files.
// Also, a list of pages from which to extract text can be passed in.
// If the pages parameter is nil or an empty slice, the text is extracted from
// all the pages of the file.
func ExtractText(inputPath, password string, pages []int) (string, error) {
	// Read input file.
	r, pageCount, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return "", err
	}

	// Extract text.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	var text string
	for _, numPage := range pages {
		// Get page.
		page, err := r.GetPage(numPage)
		if err != nil {
			return "", err
		}

		// Extract page text.
		extractor, err := uniextractor.New(page)
		if err != nil {
			return "", err
		}

		pageText, err := extractor.ExtractText()
		if err != nil {
			return "", err
		}

		text += pageText
	}

	return text, nil
}

// ExtractImages extracts all image content from the PDF file specified by the
// inputPath parameter. The extracted collection of images is saved as a ZIP
// archive at the location specified by the outputPath parameter.
// A password can be passed in, if the input file is encrypted.
// Also, a list of pages from which to extract images can be passed in.
// If the pages parameter is nil or an empty slice, the images are extracted
// from all the pages of the file.
func ExtractImages(inputPath, outputPath, password string, pages []int) (string, error) {
	// Use input file directory if no output path is specified.
	if outputPath == "" {
		dir, name := filepath.Split(inputPath)
		name = strings.TrimSuffix(name, filepath.Ext(name)) + ".zip"
		outputPath = filepath.Join(dir, name)
	}

	// Read input file.
	r, pageCount, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return "", err
	}

	// Prepare output archive.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	// Extract images.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	w := zip.NewWriter(outputFile)
	for _, numPage := range pages {
		// Get page.
		page, err := r.GetPage(numPage)
		if err != nil {
			return "", err
		}

		// Extract page images.
		extractor, err := uniextractor.New(page)
		if err != nil {
			return "", err
		}

		pageImages, err := extractor.ExtractPageImages()
		if err != nil {
			return "", err
		}

		// Add images to zip file.
		for i, pageImage := range pageImages.Images {
			img, err := pageImage.Image.ToGoImage()
			if err != nil {
				return "", err
			}

			filename, err := w.Create(fmt.Sprintf("p%d_%d.jpg", numPage, i))
			if err != nil {
				return "", err
			}

			err = jpeg.Encode(filename, img, &jpeg.Options{Quality: 100})
			if err != nil {
				return "", err
			}
		}
	}

	return outputPath, w.Close()
}
