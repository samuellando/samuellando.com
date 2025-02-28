package search

import (
	"fmt"
	"sort"
	"strings"

	"samuellando.com/internal/store"
)

// Interface for stored items that are Searchable
type Searchable interface {
	Id() int
	Title() string
	ToString() string
}

type indexItem struct {
	Path string
	Type string
	Text string
	Item Searchable
}

type scoredItem struct {
	indexItem
	score      float64
	MatchIndex int
}

type index struct {
	items []indexItem
}

type indexFunc func() index

func GenerateIndex[T Searchable](typ, basePath string, s store.Store[T]) func() index {
	return func() index {
		all, err := s.GetAll()
		if err != nil {
			return index{}
		}
		items := make([]indexItem, 0, len(all))
		for _, elem := range all {
			path := fmt.Sprintf("%s/%d", basePath, elem.Id())
			items = append(items, indexItem{Type: typ, Path: path, Text: elem.ToString(), Item: elem})
		}
		return index{items}
	}
}

func searchIndexes(search string, indexes ...indexFunc) []scoredItem {
	all := make([]indexItem, 0)
	for _, index := range indexes {
		all = append(all, index().items...)
	}
	scored := fuzzyRank(all, search)
	results := make([]scoredItem, 0)
	for _, i := range scored {
		if i.score <= 1.0/3 {
			results = append(results, i)
		}
	}
	return results
}

func fuzzyRank(elems []indexItem, search string) []scoredItem {
	scoredItems := make([]scoredItem, len(elems))
	for i, elem := range elems {
		scoredItems[i] = fuzzyScore(elem, search)
	}

	// Sort scoredItems by score in ascending order
	sort.Slice(scoredItems, func(i, j int) bool {
		return scoredItems[i].score < scoredItems[j].score
	})

	return scoredItems
}

func fuzzyScore(elem indexItem, search string) scoredItem {
	se := scoredItem{elem, 1, -1}
	if len(search) == 0 {
		return se
	}
	b := strings.ToLower(elem.Text)
	search = strings.ToLower(search)
	for i := range len(b) - len(search) {
		ld := levenshtein(search, b[i:i+len(search)])
		score := float64(ld) / float64(len(search))
		if score < se.score {
			se.score = score
			se.MatchIndex = i
		}
	}
	return se
}

func levenshtein(a, b string) int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 {
				dp[i][j] = j
			} else if j == 0 {
				dp[i][j] = i
			} else if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}

	return dp[m][n]
}
