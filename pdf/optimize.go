/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"os"
	"time"

	unipdf "github.com/unidoc/unipdf/v3/model"
	unioptimize "github.com/unidoc/unipdf/v3/model/optimize"
)

// OptimizeOpts represents the options used for optimizing PDF files.
type OptimizeOpts struct {
	// ImageQuality specifies the quality of the optimized images.
	ImageQuality int

	// ImagePPI specifies the maximum pixels per inch of the optimized images.
	ImagePPI float64
}

// OptimizeResult contains information about the optimization process.
type OptimizeResult struct {
	// Original contains information about the original file.
	Original FileStat

	// Optimized contains information about the optimized file.
	Optimized FileStat

	// Duration specifies the optimization processing time in nanoseconds.
	Duration time.Duration
}

// Optimize optimizes the PDF file specified by the inputPath parameter, using
// the provided options and saves the result at the location specified by the
// outputPath parameter. A password can be specified for encrypted input files.
func Optimize(inputPath, outputPath, password string, opts *OptimizeOpts) (*OptimizeResult, error) {
	// Initialize starting time.
	start := time.Now()

	// Get input file stat.
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	// Read input file.
	r, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return nil, err
	}

	// Copy input file contents to the output file.
	w := unipdf.NewPdfWriter()
	if err = readerToWriter(r, &w, nil); err != nil {
		return nil, err
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
		ImageUpperPPI:                   opts.ImagePPI,
	}))

	// Write output file.
	safe := inputPath == outputPath
	if err = writePDF(outputPath, &w, safe); err != nil {
		return nil, err
	}

	// Get output file stat.
	outputFileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, err
	}

	return &OptimizeResult{
		Original: FileStat{
			Name: inputPath,
			Size: inputFileInfo.Size(),
		},
		Optimized: FileStat{
			Name: outputPath,
			Size: outputFileInfo.Size(),
		},
		Duration: time.Since(start),
	}, nil
}
