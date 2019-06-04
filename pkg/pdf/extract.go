/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	uniextractor "github.com/unidoc/unipdf/v3/extractor"
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
// In addition, the image extraction process can be controlled by using the
// options parameter. If the options parameter is nil, the default image
// extraction options are used.
func ExtractImages(inputPath, outputPath, password string, pages []int,
	options *uniextractor.ImageExtractOptions) (string, int, error) {
	// Use input file directory if no output path is specified.
	if outputPath == "" {
		dir, name := filepath.Split(inputPath)
		name = strings.TrimSuffix(name, filepath.Ext(name)) + ".zip"
		outputPath = filepath.Join(dir, name)
	}

	// Read input file.
	r, pageCount, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return "", 0, err
	}

	// Extract images.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	// Create zip file.
	zipBuffer := bytes.NewBuffer(nil)
	w := zip.NewWriter(zipBuffer)
	now := time.Now()
	var countImages int

	for _, numPage := range pages {
		// Get page.
		page, err := r.GetPage(numPage)
		if err != nil {
			return "", 0, err
		}

		// Extract page images.
		extractor, err := uniextractor.New(page)
		if err != nil {
			return "", 0, err
		}

		pageImages, err := extractor.ExtractPageImages(options)
		if err != nil {
			return "", 0, err
		}

		// Add images to zip file.
		images := pageImages.Images
		countImages += len(images)

		for i, pageImage := range images {
			img, err := pageImage.Image.ToGoImage()
			if err != nil {
				return "", 0, err
			}

			filename, err := w.CreateHeader(&zip.FileHeader{
				Name:     (fmt.Sprintf("p%d_%d.jpg", numPage, i)),
				Modified: now,
			})
			if err != nil {
				return "", 0, err
			}

			err = jpeg.Encode(filename, img, &jpeg.Options{Quality: 100})
			if err != nil {
				return "", 0, err
			}
		}
	}

	if err := w.Close(); err != nil {
		return "", 0, nil
	}

	if countImages == 0 {
		return "", 0, nil
	}

	// Write output file.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", 0, err
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, zipBuffer); err != nil {
		return "", 0, err
	}

	return outputPath, countImages, nil
}
