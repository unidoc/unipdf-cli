/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/core"
	unipdf "github.com/unidoc/unidoc/pdf/model"
)

func MergePdfs(inputPaths []string, outputPath string) error {
	pdfWriter := unipdf.NewPdfWriter()

	var forms *unipdf.PdfAcroForm

	for docIdx, inputPath := range inputPaths {
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

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}

		// Handle forms.
		if pdfReader.AcroForm != nil {
			if forms == nil {
				forms = pdfReader.AcroForm
			} else {
				forms, err = MergeForms(forms, pdfReader.AcroForm, docIdx+1)
				if err != nil {
					return err
				}
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	// Set the merged forms object.
	if forms != nil {
		pdfWriter.SetForms(forms)
	}

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}

func MergeResources(r, r2 *unipdf.PdfPageResources) (*unipdf.PdfPageResources, error) {
	// Merge XObject resources.
	if r.XObject == nil {
		r.XObject = r2.XObject
	} else {
		xobjs := GetDict(r.XObject)
		if r2.XObject != nil {
			xobjs2 := GetDict(r2.XObject)
			for _, key := range xobjs2.Keys() {
				val := xobjs2.Get(key)
				xobjs.Set(key, val)
			}
		}
	}

	// Merge Colorspace resources.
	if r.ColorSpace == nil {
		r.ColorSpace = r2.ColorSpace
	} else {
		if r2.ColorSpace != nil {
			for key, val := range r2.ColorSpace.Colorspaces {
				// Add the r2 colorspaces to r.
				// Overwrite if duplicate.
				// Ensure only present once in Names.
				if _, has := r.ColorSpace.Colorspaces[key]; !has {
					r.ColorSpace.Names = append(r.ColorSpace.Names, key)
				}
				r.ColorSpace.Colorspaces[key] = val
			}
		}
	}

	// Merge ExtGState resources.
	if r.ExtGState == nil {
		r.ExtGState = r2.ExtGState
	} else {
		extgstates := GetDict(r.ExtGState)

		if r2.ExtGState != nil {
			extgstates2 := GetDict(r2.ExtGState)
			for _, key := range extgstates2.Keys() {
				val := extgstates2.Get(key)
				extgstates.Set(key, val)
			}
		}
	}

	if r.Shading == nil {
		r.Shading = r2.Shading
	} else {
		shadings := GetDict(r.Shading)
		if r2.Shading != nil {
			shadings2 := GetDict(r2.Shading)
			for _, key := range shadings2.Keys() {
				val := shadings2.Get(key)
				shadings.Set(key, val)
			}
		}
	}

	if r.Pattern == nil {
		r.Pattern = r2.Pattern
	} else {
		shadings := GetDict(r.Pattern)
		if r2.Pattern != nil {
			patterns2 := GetDict(r2.Pattern)
			for _, key := range patterns2.Keys() {
				val := patterns2.Get(key)
				shadings.Set(key, val)
			}
		}
	}

	if r.Font == nil {
		r.Font = r2.Font
	} else {
		fonts := GetDict(r.Font)
		if r2.Font != nil {
			fonts2 := GetDict(r2.Font)
			for _, key := range fonts2.Keys() {
				val := fonts2.Get(key)
				fonts.Set(key, val)
			}
		}
	}

	if r.ProcSet == nil {
		r.ProcSet = r2.ProcSet
	} else {
		procsets := GetDict(r.ProcSet)
		if r2.ProcSet != nil {
			procsets2 := GetDict(r2.ProcSet)
			for _, key := range procsets2.Keys() {
				val := procsets2.Get(key)
				procsets.Set(key, val)
			}
		}
	}

	if r.Properties == nil {
		r.Properties = r2.Properties
	} else {
		props := GetDict(r.Properties)
		if r2.Properties != nil {
			props2 := GetDict(r2.Properties)
			for _, key := range props2.Keys() {
				val := props2.Get(key)
				props.Set(key, val)
			}
		}
	}

	return r, nil
}

// Merge two interactive forms.
func MergeForms(form, form2 *unipdf.PdfAcroForm, docNum int) (*unipdf.PdfAcroForm, error) {
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
		dr, err := MergeResources(form.DR, form2.DR)
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
		field.T = core.MakeString(fmt.Sprintf("doc%d", docNum))
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
