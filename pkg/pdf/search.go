/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"strings"

	uniextractor "github.com/unidoc/unipdf/v3/extractor"
)

// SearchResult contains information about a found search term inside a PDF page.
type SearchResult struct {
	// The page the search term was found on.
	Page int

	// The number of occurrences of the search term inside the page.
	Occurrences int
}

// Search searches the provided text in the PDF file specified by the inputPath
// parameter. A password can be passed in for encrypted input files.
func Search(inputPath, text, password string) ([]*SearchResult, error) {
	// Read input file.
	r, pages, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return nil, err
	}

	// Search specified text.
	var results []*SearchResult
	for i := 0; i < pages; i++ {
		// Get page.
		numPage := i + 1

		page, err := r.GetPage(numPage)
		if err != nil {
			return nil, err
		}

		// Extract page text.
		extractor, err := uniextractor.New(page)
		if err != nil {
			return nil, err
		}

		pageText, err := extractor.ExtractText()
		if err != nil {
			return nil, err
		}

		occurrences := strings.Count(pageText, text)
		if occurrences == 0 {
			continue
		}

		results = append(results, &SearchResult{
			Page:        numPage,
			Occurrences: occurrences,
		})
	}

	return results, nil
}
