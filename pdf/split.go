/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"os"

	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func SplitPdf(inputPath, outputPath string, pageNums []int) error {
	pdfWriter := unipdf.NewPdfWriter()

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for _, pageNum := range pageNums {
		if pageNum < 0 || pageNum > numPages {
			continue
		}

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

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
