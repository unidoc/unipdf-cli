/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

type fileMatcher func(string) bool

func pdfMatcher(inputPath string) bool {
	return strings.ToLower(filepath.Ext(inputPath)) == ".pdf"
}

// parsePageRange parses a string of page ranges separated by commas and
// returns a slice of integer page numbers.
// Example page range string: 1-3,4,6-7
// The returned slice of pages contains pages 1,2,3 (1-3), 4 and 6,7 (6-7),
// while page number 5 is skipped.
func parsePageRange(pageRange string) ([]int, error) {
	var pages []int

	rngs := strings.Split(removeSpaces(pageRange), ",")
	for _, rng := range rngs {
		if rng == "" {
			continue
		}

		indices := strings.Split(rng, "-")

		lenIndices := len(indices)
		if lenIndices > 2 {
			return nil, errors.New("invalid page range")
		}
		if lenIndices == 2 {
			start, err := strconv.Atoi(indices[0])
			if err != nil {
				return nil, errors.New("invalid page number")
			}
			if start < 1 {
				return nil, errors.New("page range start must be greater than 0")
			}

			end, err := strconv.Atoi(indices[1])
			if err != nil {
				return nil, errors.New("invalid page number")
			}
			if end < 1 {
				return nil, errors.New("page range end must be greater than 0")
			}

			if start > end {
				return nil, errors.New("page range end must be greater than the start")
			}

			for page := start; page <= end; page++ {
				pages = append(pages, page)
			}

			continue
		}

		page, err := strconv.Atoi(indices[0])
		if err != nil {
			return nil, errors.New("invalid page number")
		}

		pages = append(pages, page)
	}

	pages = uniqueIntSlice(pages)
	sort.Ints(pages)

	return pages, nil
}

func parseInputPaths(inputPaths []string, recursive bool, matcher fileMatcher) ([]string, error) {
	var err error
	var files []string
	acc := map[string]bool{}

	for _, inputPath := range inputPaths {
		// Convert relative paths to absolute ones.
		if !filepath.IsAbs(inputPath) {
			inputPath, err = filepath.Abs(inputPath)
			if err != nil {
				return nil, err
			}
		}

		// Add visited file to the accumulator.
		if _, ok := acc[inputPath]; ok {
			continue
		}
		acc[inputPath] = true

		// Get file info.
		inputFile, err := os.Lstat(inputPath)
		if err != nil {
			return nil, err
		}

		// Check file type.
		switch mode := inputFile.Mode(); {
		case mode.IsRegular():
			if matcher == nil || matcher(inputPath) {
				files = append(files, inputPath)
			}
		case mode.IsDir():
			dirFiles, err := parseInputDir(inputPath, recursive, acc, matcher)
			if err != nil {
				return nil, err
			}
			files = append(files, dirFiles...)
		}
	}

	return files, nil
}

func parseInputDir(dir string, recursive bool, acc map[string]bool, matcher fileMatcher) ([]string, error) {
	inputPaths, err := dirFiles(dir)
	if err != nil {
		return nil, err
	}

	if acc == nil {
		acc = map[string]bool{}
	}

	var files []string
	for _, inputPath := range inputPaths {
		// Convert relative paths to absolute ones.
		inputPath = filepath.Join(dir, inputPath)
		if !filepath.IsAbs(inputPath) {
			inputPath, err = filepath.Abs(inputPath)
			if err != nil {
				return nil, err
			}
		}

		// Add visited file to the accumulator.
		if _, ok := acc[inputPath]; ok {
			continue
		}
		acc[inputPath] = true

		// Get file info.
		inputFile, err := os.Lstat(inputPath)
		if err != nil {
			return nil, err
		}

		// Check file type.
		switch mode := inputFile.Mode(); {
		case mode.IsRegular():
			if matcher == nil || matcher(inputPath) {
				files = append(files, inputPath)
			}
		case mode.IsDir():
			if !recursive {
				continue
			}

			subdirFiles, err := parseInputDir(inputPath, recursive, acc, matcher)
			if err != nil {
				return nil, err
			}
			files = append(files, subdirFiles...)
		}
	}

	return files, nil
}

func dirFiles(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdirnames(-1)
}

func generateOutputPath(inputPath, outputDir, nameSuffix string, overwrite bool) string {
	if overwrite {
		return inputPath
	}

	dir, name := filepath.Split(inputPath)
	if outputDir != "" {
		return filepath.Join(outputDir, name)
	}

	name = strings.TrimSuffix(name, filepath.Ext(name))
	return filepath.Join(dir, fmt.Sprintf("%s_%s.pdf", name, nameSuffix))
}

func clampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}

	return val
}

func removeSpaces(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func uniqueIntSlice(items []int) []int {
	uniq := make([]int, len(items))

	index := 0
	catalog := map[int]struct{}{}
	for _, item := range items {
		if _, ok := catalog[item]; ok {
			continue
		}

		catalog[item] = struct{}{}
		uniq[index] = item
		index++
	}

	return uniq[0:index]
}

func printErr(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func printUsageErr(cmd *cobra.Command, format string, a ...interface{}) {
	fmt.Printf("Error: "+format+"\n", a...)
	cmd.Help()
	os.Exit(1)
}
