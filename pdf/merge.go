/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"fmt"

	unicommon "github.com/unidoc/unipdf/v3/common"
	unicore "github.com/unidoc/unipdf/v3/core"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

// Merge merges all the PDF files specified by the inputPaths parameter and
// saves the result at the location specified by the outputPath parameter.
func Merge(inputPaths []string, outputPath string) error {
	w := unipdf.NewPdfWriter()

	var forms *unipdf.PdfAcroForm
	for index, inputPath := range inputPaths {
		// Read file.
		r, pages, _, _, err := readPDF(inputPath, "")
		if err != nil {
			return err
		}

		// Add pages.
		for i := 0; i < pages; i++ {
			page, err := r.GetPage(i + 1)
			if err != nil {
				return err
			}

			err = w.AddPage(page)
			if err != nil {
				return err
			}
		}

		// Handle forms.
		if r.AcroForm != nil {
			if forms == nil {
				forms = r.AcroForm
			} else {
				forms, err = mergeForms(forms, r.AcroForm, index+1)
				if err != nil {
					return err
				}
			}
		}
	}

	// Set the merged forms object.
	if forms != nil {
		w.SetForms(forms)
	}

	// Write output file.
	return writePDF(outputPath, &w, false)
}

func mergeResources(r, r2 *unipdf.PdfPageResources) (*unipdf.PdfPageResources, error) {
	// Merge XObject resources.
	if r.XObject == nil {
		r.XObject = r2.XObject
	} else {
		xobjs := getDict(r.XObject)
		if r2.XObject != nil {
			xobjs2 := getDict(r2.XObject)
			for _, key := range xobjs2.Keys() {
				val := xobjs2.Get(key)
				xobjs.Set(key, val)
			}
		}
	}

	// Merge Colorspace resources.
	colorspaces, err := r.GetColorspaces()
	if err != nil {
		return nil, err
	}
	colorspaces2, err := r2.GetColorspaces()
	if err != nil {
		return nil, err
	}

	if colorspaces == nil {
		r.SetColorSpace(colorspaces2)
	} else {
		if colorspaces2 != nil {
			for key, val := range colorspaces2.Colorspaces {
				// Add the r2 colorspaces to r. Overwrite if duplicate.
				// Ensure only present once in Names.
				if _, has := colorspaces.Colorspaces[key]; !has {
					colorspaces.Names = append(colorspaces.Names, key)
				}
				r.SetColorspaceByName(unicore.PdfObjectName(key), val)
			}
		}
	}

	// Merge ExtGState resources.
	if r.ExtGState == nil {
		r.ExtGState = r2.ExtGState
	} else {
		extgstates := getDict(r.ExtGState)

		if r2.ExtGState != nil {
			extgstates2 := getDict(r2.ExtGState)
			for _, key := range extgstates2.Keys() {
				val := extgstates2.Get(key)
				extgstates.Set(key, val)
			}
		}
	}

	if r.Shading == nil {
		r.Shading = r2.Shading
	} else {
		shadings := getDict(r.Shading)
		if r2.Shading != nil {
			shadings2 := getDict(r2.Shading)
			for _, key := range shadings2.Keys() {
				val := shadings2.Get(key)
				shadings.Set(key, val)
			}
		}
	}

	if r.Pattern == nil {
		r.Pattern = r2.Pattern
	} else {
		shadings := getDict(r.Pattern)
		if r2.Pattern != nil {
			patterns2 := getDict(r2.Pattern)
			for _, key := range patterns2.Keys() {
				val := patterns2.Get(key)
				shadings.Set(key, val)
			}
		}
	}

	if r.Font == nil {
		r.Font = r2.Font
	} else {
		fonts := getDict(r.Font)
		if r2.Font != nil {
			fonts2 := getDict(r2.Font)
			for _, key := range fonts2.Keys() {
				val := fonts2.Get(key)
				fonts.Set(key, val)
			}
		}
	}

	if r.ProcSet == nil {
		r.ProcSet = r2.ProcSet
	} else {
		procsets := getDict(r.ProcSet)
		if r2.ProcSet != nil {
			procsets2 := getDict(r2.ProcSet)
			for _, key := range procsets2.Keys() {
				val := procsets2.Get(key)
				procsets.Set(key, val)
			}
		}
	}

	if r.Properties == nil {
		r.Properties = r2.Properties
	} else {
		props := getDict(r.Properties)
		if r2.Properties != nil {
			props2 := getDict(r2.Properties)
			for _, key := range props2.Keys() {
				val := props2.Get(key)
				props.Set(key, val)
			}
		}
	}

	return r, nil
}

// mergeForms merges two interactive forms.
func mergeForms(form, form2 *unipdf.PdfAcroForm, docNum int) (*unipdf.PdfAcroForm, error) {
	if form.NeedAppearances == nil {
		form.NeedAppearances = form2.NeedAppearances
	}

	if form.SigFlags == nil {
		form.SigFlags = form2.SigFlags
	}

	if form.CO == nil {
		form.CO = form2.CO
	}

	if form.DR == nil {
		form.DR = form2.DR
	} else if form2.DR != nil {
		dr, err := mergeResources(form.DR, form2.DR)
		if err != nil {
			return nil, err
		}
		form.DR = dr
	}

	if form.DA == nil {
		form.DA = form2.DA
	}

	if form.Q == nil {
		form.Q = form2.Q
	}

	if form.XFA == nil {
		form.XFA = form2.XFA
	} else {
		if form2.XFA != nil {
			unicommon.Log.Debug("TODO: Handle XFA merging - Currently just using first one that is encountered")
		}
	}

	// Fields.
	if form.Fields == nil {
		form.Fields = form2.Fields
	} else {
		field := unipdf.NewPdfField()
		field.T = unicore.MakeString(fmt.Sprintf("doc%d", docNum))
		field.Kids = []*unipdf.PdfField{}
		if form2.Fields != nil {
			for _, subfield := range *form2.Fields {
				// Update parent.
				subfield.Parent = field
				field.Kids = append(field.Kids, subfield)
			}

		}
		*form.Fields = append(*form.Fields, field)
	}

	return form, nil
}

func getDict(obj unicore.PdfObject) *unicore.PdfObjectDictionary {
	if obj == nil {
		return nil
	}

	obj = unicore.TraceToDirectObject(obj)
	dict, ok := obj.(*unicore.PdfObjectDictionary)
	if !ok {
		unicommon.Log.Debug("Error type check error (got %T)", obj)
		return nil
	}

	return dict
}
