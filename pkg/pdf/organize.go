/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"github.com/unidoc/unipdf/v3/common"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

// Organize extracts the provided page list from PDF file specified by the
// inputPath parameter then merges the individual pages and saves the
// resulting file at the location specified by the outputPath parameter.
// A password can be passed in for encrypted input files.
func Organize(inputPath, outputPath, password string, pages []int) error {
	// Read input file.
	pdfReader, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Add selected pages to the writer.
	pdfWriter := unipdf.NewPdfWriter()

	for i := 0; i < len(pages); i++ {
		page, err := pdfReader.GetPage(pages[i])
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}

	// Copy PDF version.
	version := pdfReader.PdfVersion()
	pdfWriter.SetVersion(version.Major, version.Minor)

	// Copy PDF info.
	info, err := pdfReader.GetPdfInfo()
	if err != nil {
		common.Log.Debug("ERROR: %v", err)
	} else {
		pdfWriter.SetDocInfo(info)
	}

	// Copy Catalog Metadata.
	if meta, ok := pdfReader.GetCatalogMetadata(); ok {
		if err := pdfWriter.SetCatalogMetadata(meta); err != nil {
			return err
		}
	}

	// Copy catalog mark information.
	if markInfo, ok := pdfReader.GetCatalogMarkInfo(); ok {
		if err := pdfWriter.SetCatalogMarkInfo(markInfo); err != nil {
			return err
		}
	}

	// Copy AcroForm.
	err = pdfWriter.SetForms(pdfReader.AcroForm)
	if err != nil {
		common.Log.Debug("ERROR: %v", err)
		return err
	}

	// Copy viewer preferences.
	if pref, ok := pdfReader.GetCatalogViewerPreferences(); ok {
		if err := pdfWriter.SetCatalogViewerPreferences(pref); err != nil {
			return err
		}
	}

	// Copy language preferences.
	if lang, ok := pdfReader.GetCatalogLanguage(); ok {
		if err := pdfWriter.SetCatalogLanguage(lang); err != nil {
			return err
		}
	}

	// Copy document outlines.
	pdfWriter.AddOutlineTree(pdfReader.GetOutlineTree())

	// Copy OC Properties.
	props, err := pdfReader.GetOCProperties()
	if err != nil {
		common.Log.Debug("ERROR: %v", err)
	} else {
		err = pdfWriter.SetOCProperties(props)
		if err != nil {
			common.Log.Debug("ERROR: %v", err)
		}
	}

	// Copy page labels.
	labelObj, err := pdfReader.GetPageLabels()
	if err != nil {
		common.Log.Debug("ERROR: %v", err)
	} else {
		err = pdfWriter.SetPageLabels(labelObj)
		if err != nil {
			common.Log.Debug("ERROR: %v", err)
		}
	}

	// Copy named destinations.
	namedDest, err := pdfReader.GetNamedDestinations()
	if err != nil {
		common.Log.Debug("ERROR: %v", err)
	} else {
		err = pdfWriter.SetNamedDestinations(namedDest)
		if err != nil {
			common.Log.Debug("ERROR: %v", err)
		}
	}

	// Copy name dictionary.
	nameDict, err := pdfReader.GetNameDictionary()
	if err != nil {
		common.Log.Debug("ERROR: %v", err)
	} else {
		err = pdfWriter.SetNameDictionary(nameDict)
		if err != nil {
			common.Log.Debug("ERROR: %v", err)
		}
	}

        // Copy StructTreeRoot dictionary.
	structTreeRoot, found := pdfReader.GetCatalogStructTreeRoot()
	if found {
		err := pdfWriter.SetCatalogStructTreeRoot(structTreeRoot)
		if err != nil {
			common.Log.Debug("ERROR: %v", err)
		}
	}

	// Copy global page rotation.
	if pdfReader.Rotate != nil {
		if err := pdfWriter.SetRotation(*pdfReader.Rotate); err != nil {
			common.Log.Debug("ERROR: %v", err)
		}
	}

	// Write output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &pdfWriter, safe)
}
