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

	unicontent "github.com/unidoc/unidoc/pdf/contentstream"
	unicore "github.com/unidoc/unidoc/pdf/core"
	uniextractor "github.com/unidoc/unidoc/pdf/extractor"
	unipdf "github.com/unidoc/unidoc/pdf/model"
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

		// List images on the page.
		rgbImages, err := extractImagesOnPage(page)
		if err != nil {
			return "", err
		}

		// Add images to zip file.
		for i, img := range rgbImages {
			img, err := img.ToGoImage()
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

func extractImagesOnPage(page *unipdf.PdfPage) ([]*unipdf.Image, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}

	return extractImagesInContentStream(contents, page.Resources)
}

func extractImagesInContentStream(contents string, resources *unipdf.PdfPageResources) ([]*unipdf.Image, error) {
	rgbImages := []*unipdf.Image{}
	cstreamParser := unicontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}

	// Range through all the content stream operations.
	processedXObjects := map[string]bool{}

	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			iimg, ok := op.Params[0].(*unicontent.ContentStreamInlineImage)
			if !ok {
				continue
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				return nil, err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return nil, err
			}
			if cs == nil {
				cs = unipdf.NewPdfColorspaceDeviceGray()
			}

			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				return nil, err
			}

			rgbImages = append(rgbImages, &rgbImg)
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			name := op.Params[0].(*unicore.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == unipdf.XObjectTypeImage {
				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return nil, err
				}

				img, err := ximg.ToImage()
				if err != nil {
					return nil, err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					return nil, err
				}
				rgbImages = append(rgbImages, &rgbImg)
			} else if xtype == unipdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return nil, err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return nil, err
				}

				// Process the content stream in the Form object too.
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too.
				formRgbImages, err := extractImagesInContentStream(string(formContent), formResources)
				if err != nil {
					return nil, err
				}
				rgbImages = append(rgbImages, formRgbImages...)
			}
		}
	}

	return rgbImages, nil
}
