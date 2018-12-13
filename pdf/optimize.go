/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"os"

	unipdf "github.com/unidoc/unidoc/pdf/model"
	unioptimize "github.com/unidoc/unidoc/pdf/model/optimize"
)

type OptimizeOpts struct {
	ImageQuality int
}

func OptimizePdf(inputPath, outputPath, password string, opts *OptimizeOpts) error {
	// Read input file.
	r, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Copy input file contents to the output file.
	w := unipdf.NewPdfWriter()
	if err = readerToWriter(r, &w); err != nil {
		return err
	}

	// Add optimizer.
	if opts == nil {
		opts = &OptimizeOpts{
			ImageQuality: 100,
		}
	}

	w.SetOptimizer(unioptimize.New(unioptimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    opts.ImageQuality,
	}))

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
