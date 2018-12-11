/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmd

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

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
			return nil, errors.New("Invalid page range")
		}
		if lenIndices == 2 {
			start, err := strconv.Atoi(indices[0])
			if err != nil {
				return nil, errors.New("Invalid page number")
			}
			if start < 1 {
				return nil, errors.New("Page range start must be greater than 0")
			}

			end, err := strconv.Atoi(indices[1])
			if err != nil {
				return nil, errors.New("Invalid page number")
			}
			if end < 1 {
				return nil, errors.New("Page range end must be greater than 0")
			}

			if start > end {
				return nil, errors.New("Page range end must be greater than the start")
			}

			for page := start; page <= end; page++ {
				pages = append(pages, page)
			}

			continue
		}

		page, err := strconv.Atoi(indices[0])
		if err != nil {
			return nil, errors.New("Invalid page number")
		}

		pages = append(pages, page)
	}

	pages = UniqueIntSlice(pages)
	sort.Ints(pages)

	return pages, nil
}

func removeSpaces(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func UniqueIntSlice(items []int) []int {
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
