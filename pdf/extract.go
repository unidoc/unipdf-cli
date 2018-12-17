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
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func ExtractText(inputPath, password string) (string, error) {
	// Read input file.
	r, pages, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return "", err
	}

	// Extract text.
	var text string
	for i := 0; i < pages; i++ {
		// Get page.
		page, err := r.GetPage(i + 1)
		if err != nil {
			return "", err
		}

		// Get page streams.
		streams, err := page.GetContentStreams()
		if err != nil {
			return "", err
		}

		var pageContent string
		for _, stream := range streams {
			pageContent += stream
		}

		// Extract page text.
		parser := unicontent.NewContentStreamParser(pageContent)

		pageText, err := parser.ExtractText()
		if err != nil {
			return "", err
		}

		text += pageText
	}

	return text, nil
}

func ExtractImages(inputPath, outputPath, password string) error {
	if outputPath == "" {
		dir, name := filepath.Split(inputPath)
		name = strings.TrimSuffix(name, filepath.Ext(name)) + ".zip"
		outputPath = filepath.Join(dir, name)
	}

	// Read input file.
	r, pages, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Prepare output archive.
	zipf, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipf.Close()

	// Extract images.
	zipw := zip.NewWriter(zipf)
	for i := 0; i < pages; i++ {
		// Get page.
		page, err := r.GetPage(i + 1)
		if err != nil {
			return err
		}

		// List images on the page.
		rgbImages, err := extractImagesOnPage(page)
		if err != nil {
			return err
		}

		// Add images to zip.
		for idx, img := range rgbImages {
			gimg, err := img.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := zipw.Create(fmt.Sprintf("p%d_%d.jpg", i+1, idx))
			if err != nil {
				return err
			}

			err = jpeg.Encode(imgf, gimg, &jpeg.Options{Quality: 100})
			if err != nil {
				return err
			}
		}
	}

	return zipw.Close()
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
			// BI: Inline image.
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
				// Default if not specified?
				cs = unipdf.NewPdfColorspaceDeviceGray()
			}
			fmt.Printf("Cs: %T\n", cs)

			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				return nil, err
			}

			rgbImages = append(rgbImages, &rgbImg)
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// Do: XObject.
			name := op.Params[0].(*unicore.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == unipdf.XObjectTypeImage {
				fmt.Printf(" XObject Image: %s\n", *name)

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

				// Process the content stream in the Form object too:
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
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