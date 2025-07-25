package coverparse

import (
	"cmp"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func IsCover(p string) bool {
	p = filepath.Ext(p)
	p = strings.ToLower(p)
	_, ok := filetypePriorities[p]
	return ok
}

// Compare ranks two potential cover paths, suitable for [slices.SortFunc].
func Compare(a, b string) int {
	a, b = strings.ToLower(a), strings.ToLower(b)
	return cmp.Or(
		slices.Compare(posArtTypes(a), posArtTypes(b)),
		slices.Compare(posNumbers(a), posNumbers(b)),
		cmp.Compare(posFiletype(a), posFiletype(b)),
	)
}

// BestBetween updates the current best candidate if the new path is better.
func BestBetween(cover string, other string) string {
	if cover == "" {
		return other
	}
	if Compare(cover, other) > 0 {
		return other
	}
	return cover
}

var artTypePriorities = map[string]int{
	"front":    3,
	"cover":    3,
	"album":    3,
	"folder":   2,
	"albumart": 2,
	"scan":     1,
	"back":     0, // ignore
	"artist":   0, // ignore
}

var artTypeExpr *regexp.Regexp

func init() {
	var quoted []string
	for k := range artTypePriorities {
		quoted = append(quoted, regexp.QuoteMeta(k))
	}
	slices.SortFunc(quoted, func(a, b string) int {
		return cmp.Or(
			len(b)-len(a),
			cmp.Compare(a, b),
		)
	})
	quoteExpr := strings.Join(quoted, "|")
	artTypeExpr = regexp.MustCompile(quoteExpr)
}

func posArtTypes(path string) []int {
	matches := artTypeExpr.FindAllString(path, -1)
	if len(matches) == 0 {
		return []int{1} // a score of 1 is worse than any priority (which are <= 0)
	}
	r := make([]int, len(matches))
	for i, m := range matches {
		r[i] = -artTypePriorities[m]
	}
	return r
}

var numbersExpr = regexp.MustCompile(`\d+`)

func posNumbers(path string) []int {
	matches := numbersExpr.FindAllString(path, -1)
	r := make([]int, len(matches))
	for i, m := range matches {
		r[i], _ = strconv.Atoi(m)
	}
	return r
}

var filetypePriorities = map[string]int{
	".png":  2,
	".jpg":  1,
	".jpeg": 1,
	".bmp":  1,
	".gif":  1,
}

func posFiletype(path string) int {
	return -filetypePriorities[filepath.Ext(path)]
}
