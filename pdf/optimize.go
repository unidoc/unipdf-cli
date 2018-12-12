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
	// Open input file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read PDF file.
	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	// Decrypt file if necessary.
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(password))
		if err != nil {
			return err
		}
	}

	// Copy input file contents to the output file.
	pdfWriter := unipdf.NewPdfWriter()
	if err = readerToWriter(pdfReader, &pdfWriter); err != nil {
		return err
	}

	// Add optimizer.
	if opts == nil {
		opts = &OptimizeOpts{
			ImageQuality: 100,
		}
	}

	pdfWriter.SetOptimizer(unioptimize.New(unioptimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    opts.ImageQuality,
	}))

	// Write output file.
	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
