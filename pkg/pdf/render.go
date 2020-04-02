/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"archive/zip"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/render"
)

// RenderOpts represents the options used for rendering PDF pages to images.
type RenderOpts struct {
	// ImageFormat specifies the file format of the rendered images.
	// Supported formats: jpeg, png.
	ImageFormat string

	// ImageQuality specifies the quality of the rendered images.
	// Only applies to rendered JPEG images.
	ImageQuality int
}

// Render renders the pages of the PDF file specified by the inputPath parameter
// to image targets. The rendered images are saved as a ZIP archive at the
// location specified by the outputPath parameter.
// A password can be passed in, if the input file is encrypted.
// If the pages parameter is nil or an empty slice, all pages are rendered.
func Render(inputPath, outputPath, password string, pages []int, opts *RenderOpts) (string, error) {
	// Use input file directory if no output path is specified.
	dir, inputFile := filepath.Split(inputPath)

	inputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	if outputPath == "" {
		outputPath = filepath.Join(dir, inputFile+".zip")
	}

	// Read input file.
	r, pageCount, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return "", err
	}

	// Extract pages.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	// Create render options, if none are specified.
	if opts == nil {
		opts = &RenderOpts{ImageFormat: "jpeg", ImageQuality: 100}
	}
	if opts.ImageQuality < 0 || opts.ImageQuality > 100 {
		opts.ImageQuality = 100
	}

	// Create image encode function.
	var encodeFunc func(w io.Writer, img image.Image) error
	imgExt := "jpg"

	switch opts.ImageFormat {
	case "jpeg":
		encodeFunc = func(w io.Writer, img image.Image) error {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: opts.ImageQuality})
		}
	case "png":
		imgExt = "png"
		encodeFunc = func(w io.Writer, img image.Image) error {
			return png.Encode(w, img)
		}
	default:
		return "", fmt.Errorf("unsupported image format: %s", opts.ImageFormat)
	}

	// Prepare output archive.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	zw := zip.NewWriter(outputFile)

	// Render pages.
	device := render.NewImageDevice()
	for _, numPage := range pages {
		// Get page.
		page, err := r.GetPage(numPage)
		if err != nil {
			return "", err
		}

		// Render page to image.
		img, err := device.Render(page)
		if err != nil {
			return "", err
		}

		// Add rendered image to zip file.
		file, err := zw.Create(fmt.Sprintf("%s_%d.%s", inputFile, numPage, imgExt))
		if err != nil {
			return "", err
		}
		if err := encodeFunc(file, img); err != nil {
			return "", err
		}
	}

	return outputPath, zw.Close()
}
