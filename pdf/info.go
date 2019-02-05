/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"os"
	"sort"
)

// FileStat contains basic information about a file.
type FileStat struct {
	// Name represents the name of the file.
	Name string

	// Size specifies the size in bytes of the file.
	Size int64
}

// FileInfo contains information about a PDF file.
type FileInfo struct {
	FileStat

	// Pages represents the number of pages the PDF file has.
	Pages int

	// Objects contains the types of objects the PDF file contains, along
	// with the count for each object type.
	Objects map[string]int

	// Version specifies the PDF version of the file.
	Version string

	// Encrypted specifies if the file is encrypted.
	Encrypted bool

	// EncryptionAlgo contains the name of the encryption algorithm used
	// to encrypt the PDF file. The field is empty for non-encrypted files.
	EncryptionAlgo string
}

// Info returns information about the PDF file specified by the inputPath
// parameter. A password can be passed in for encrypted input files.
func Info(inputPath string, password string) (*FileInfo, error) {
	info := &FileInfo{}
	info.Name = inputPath

	// Get file stat.
	fi, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}
	info.Size = fi.Size()

	// Read input file.
	r, pages, encrypted, _, err := readPDF(inputPath, password)
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
	for key := range objTypes {
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
