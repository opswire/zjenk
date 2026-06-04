package tool

import (
	"sort"

	"mcp-jenkins/internal/dto"
)

func paginateText(text string, page, charsPerPage, maxCharsPerPage int64) (string, dto.Pagination) {
	if charsPerPage <= 0 || charsPerPage > maxCharsPerPage {
		charsPerPage = maxCharsPerPage
	}

	total := int64(len(text))
	totalPages := (total + charsPerPage - 1) / charsPerPage
	if totalPages == 0 {
		totalPages = 1
	}
	if page <= 0 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * charsPerPage
	end := start + charsPerPage
	if end > total {
		end = total
	}

	return text[start:end], dto.Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		CharsPerPage: charsPerPage,
	}
}

// applyContext expands matched line indices by contextLines in each direction
// and returns a deduplicated, sorted slice of SearchMatch entries.
// Match=true marks lines that directly matched the search pattern.
func applyContext(lines []string, matchIndices []int, contextLines, maxContextLines int64) []dto.SearchMatch {
	if contextLines > maxContextLines {
		contextLines = maxContextLines
	}
	if contextLines < 0 {
		contextLines = 0
	}

	total := int64(len(lines))
	matchSet := make(map[int]bool, len(matchIndices))
	for _, idx := range matchIndices {
		matchSet[idx] = true
	}

	included := make(map[int]bool)
	for idx := range matchSet {
		lo := int64(idx) - contextLines
		if lo < 0 {
			lo = 0
		}
		hi := int64(idx) + contextLines
		if hi >= total {
			hi = total - 1
		}
		for i := lo; i <= hi; i++ {
			included[int(i)] = true
		}
	}

	indices := make([]int, 0, len(included))
	for i := range included {
		indices = append(indices, i)
	}
	sort.Ints(indices)

	result := make([]dto.SearchMatch, 0, len(indices))
	for _, i := range indices {
		result = append(result, dto.SearchMatch{
			Line:    i + 1,
			Content: lines[i],
			Match:   matchSet[i],
		})
	}

	return result
}

func paginateMatches(matches []dto.SearchMatch, page, pageSize int64) ([]dto.SearchMatch, dto.Pagination) {
	total := int64(len(matches))
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	if page <= 0 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	return matches[start:end], dto.Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		CharsPerPage: pageSize,
	}
}
