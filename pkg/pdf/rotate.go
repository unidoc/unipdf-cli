/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	unicreator "github.com/unidoc/unipdf/v3/creator"
)

// Rotate rotates the pages of the PDF file specified by the inputPath
// by the angle specified by the angle parameter. The rotated PDF file is saved
// at the location specified by the outputPath parameter.
// A password can be passed in, if the input file is encrypted.
// If the pages parameter is nil or an empty slice, all pages are rotated.
func Rotate(inputPath, outputPath string, angle int, password string, pages []int) (string, error) {
	if angle%90 != 0 {
		return "", errors.New("rotation angle must be a multiple of 90 degrees")
	}

	// Generate output path from the input path, if no output path is specified.
	dir, inputFile := filepath.Split(inputPath)

	inputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	if outputPath == "" {
		outputPath = filepath.Join(dir, fmt.Sprintf("%s_rotated.pdf", inputFile))
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

	// Rotate pages.
	if len(pages) == 0 {
		pages = createPageRange(pageCount)
	}

	selectedPages := map[int]bool{}
	for _, page := range pages {
		selectedPages[page] = true
	}

	c := unicreator.New()
	for i := 0; i < pageCount; i++ {
		numPage := i + 1

		page, err := r.GetPage(numPage)
		if err != nil {
			return "", err
		}

		if err = c.AddPage(page); err != nil {
			return "", err
		}

		rotate, _ := selectedPages[numPage]
		if !rotate || angle == 0 {
			continue
		}

		if err = c.RotateDeg(int64(angle)); err != nil {
			return "", err
		}
	}

	// Add forms.
	if r.AcroForm != nil {
		c.SetForms(r.AcroForm)
	}

	// Write output file.
	safe := inputPath == outputPath
	return outputPath, writeCreatorPDF(outputPath, c, safe)
}
