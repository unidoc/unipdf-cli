/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"errors"
	"os"
	"sort"

	unipdf "github.com/unidoc/unidoc/pdf/model"
)

type PDFInfo struct {
	Name    string
	Pages   int
	Size    int64
	Objects map[string]int

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

	// Open and read input file.
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := unipdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	// Check if encrypted and try to decrypt using the specified password.
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}
	info.Encrypted = isEncrypted

	if isEncrypted {
		info.EncryptionAlgo = pdfReader.GetEncryptionMethod()

		auth, err := pdfReader.Decrypt([]byte(password))
		if err != nil {
			return nil, err
		}
		if !auth {
			return nil, errors.New("Unable to decrypt PDF with the specified password")
		}
	}

	// Get number of pages.
	info.Pages, err = pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	// Read PDF objects.
	objTypes, err := pdfReader.Inspect()
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
