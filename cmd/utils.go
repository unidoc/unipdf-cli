/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

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
