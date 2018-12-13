/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"os"
	"sort"
)

type PDFInfo struct {
	Name    string
	Pages   int
	Size    int64
	Objects map[string]int
	Version string

	Encrypted      bool
	EncryptionAlgo string
}

func GetPDFInfo(inputPath string, password string) (*PDFInfo, error) {
	info := &PDFInfo{
		Name: inputPath,
	}

	// Get file stat.
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}
	info.Size = fileInfo.Size()

	// Read input file.
	r, pages, encrypted, err := readPDF(inputPath, password)
	if err != nil {
		return nil, err
	}

	info.Encrypted = encrypted
	if encrypted {
		info.EncryptionAlgo = r.GetEncryptionMethod()
	}

	info.Version = r.PdfVersion().String()
	info.Pages = pages

	// Read PDF objects.
	objTypes, err := r.Inspect()
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for key, _ := range objTypes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	objects := map[string]int{}
	for _, key := range keys {
		objects[key] = objTypes[key]
	}
	info.Objects = objects

	return info, nil
}
