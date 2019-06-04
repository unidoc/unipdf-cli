/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/fdf"
	"github.com/unidoc/unipdf/v3/fjson"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

// FormExport exports all form field values from the PDF file specified
// by the inputPath parameters, as JSON.
func FormExport(inputPath string) (string, error) {
	fieldData, err := fjson.LoadFromPDFFile(inputPath)
	if err != nil {
		return "", err
	}
	if fieldData == nil {
		return "", nil
	}

	return fieldData.JSON()
}

// FormFillJSON fills the form field values from the PDF file specified by the
// inputPath parameter, using the values from the JSON file specified by the
// jsonPath parameter. The output PDF file is saved at the location specified
// by the outputPath parameter. The output file form annotations can be
// flattened by using the flatten parameter.
// A password can be specified for encrypted input files.
func FormFillJSON(inputPath, jsonPath, outputPath, password string, flatten bool) error {
	// Read JSON field data.
	fieldData, err := fjson.LoadFromJSONFile(jsonPath)
	if err != nil {
		return err
	}

	return formFill(inputPath, fieldData, outputPath, password, flatten)
}

// FormFillFDF fills the form field values from the PDF file specified by the
// inputPath parameter, using the values from the FDF file specified by the
// fdfPath parameter. The output PDF file is saved at the location specified
// by the outputPath parameter. The output file form annotations can be
// flattened by using the flatten parameter.
// A password can be specified for encrypted input files.
func FormFillFDF(inputPath, fdfPath, outputPath, password string, flatten bool) error {
	// Read field data.
	fieldData, err := fdf.LoadFromPath(fdfPath)
	if err != nil {
		return err
	}

	return formFill(inputPath, fieldData, outputPath, password, flatten)
}

// FormFlatten flattens all the form annotation from the PDF file specified by
// the inputPath parameter. The output PDF file is saved at the location
// specified by the outputPath parameter.
// A password can be specified for encrypted input files.
func FormFlatten(inputPath, outputPath, password string) error {
	// Read input file.
	r, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Flatten form.
	fieldAppearance := annotator.FieldAppearance{
		OnlyIfMissing: true,
	}

	if err = r.FlattenFields(true, fieldAppearance); err != nil {
		return err
	}
	r.AcroForm = nil

	// Copy input file contents.
	w := unipdf.NewPdfWriter()
	if err := readerToWriter(r, &w, nil); err != nil {
		return err
	}

	// Save output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}

func formFill(inputPath string, provider unipdf.FieldValueProvider, outputPath, password string, flatten bool) error {
	// Read input file.
	r, _, _, _, err := readPDF(inputPath, password)
	if err != nil {
		return err
	}

	// Populate the form data.
	if err = r.AcroForm.Fill(provider); err != nil {
		return err
	}

	// Flatten form.
	if flatten {
		fieldAppearance := annotator.FieldAppearance{
			OnlyIfMissing:        true,
			RegenerateTextFields: true,
		}

		if err = r.FlattenFields(true, fieldAppearance); err != nil {
			return err
		}
		r.AcroForm = nil
	}

	// Copy input file contents.
	w := unipdf.NewPdfWriter()
	if err := readerToWriter(r, &w, nil); err != nil {
		return err
	}

	// Save output file.
	safe := inputPath == outputPath
	return writePDF(outputPath, &w, safe)
}
