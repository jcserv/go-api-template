package httputil

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type QueryParams struct {
	PaginateBy *OffsetPagination `json:"paginate"`
	FilterBy   []Filter          `json:"filters"`
	SortBy     []Sort            `json:"sorters"`
}

// ParseQueryParams parses query parameters for pagination, filtering, and sorting
// TODO: Provide set of valid filters, sorters
func ParseQueryParams(r *http.Request) (*QueryParams, error) {
	if len(r.URL.Query()) == 0 {
		return nil, nil
	}

	paginateBy, err := parsePagination(r)
	if err != nil {
		return nil, err
	}

	filters, err := parseFilters(r)
	if err != nil {
		return nil, err
	}

	sorters, err := parseSorting(r)
	if err != nil {
		return nil, err
	}

	return &QueryParams{
		PaginateBy: paginateBy,
		FilterBy:   filters,
		SortBy:     sorters,
	}, nil
}

const (
	MAX_LIMIT = 250
)

type OffsetPagination struct {
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	Count  bool `json:"count"`
}

func parsePagination(r *http.Request) (*OffsetPagination, error) {
	query := r.URL.Query()
	limitVal := query.Get("limit")
	limit, _ := strconv.Atoi(limitVal)
	if limit < 0 {
		return nil, fmt.Errorf("limit must be greater than or equal to 0")
	}
	if limit > MAX_LIMIT {
		return nil, fmt.Errorf("limit must be less than or equal to %d", MAX_LIMIT)
	}

	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		return nil, fmt.Errorf("offset must be greater than or equal to 0")
	}

	// If no limit is provided, use the max limit
	if limitVal == "" {
		limit = MAX_LIMIT
	}

	count := query.Get("count") == "true"
	return &OffsetPagination{
		Limit:  limit,
		Offset: offset,
		Count:  count,
	}, nil
}

// source: https://bookstack.soffid.com/books/scim/page/scim-query-syntax

type FilterOp string

const (
	Eq FilterOp = "eq"
)

type Filter struct {
	Field string   `json:"field"`
	Op    FilterOp `json:"op"`
	Value any      `json:"value"`
}

func parseFilters(r *http.Request) ([]Filter, error) {
	query := r.URL.Query()
	filterString := query.Get("filter")
	if filterString == "" {
		return nil, nil
	}

	var filters []Filter
	filterParts := strings.Split(filterString, "and")

	for _, part := range filterParts {
		fields := strings.Fields(part)
		if len(fields) < 3 {
			return nil, fmt.Errorf("invalid filter format: %s", part)
		}

		field := fields[0]
		op := strings.ToLower(fields[1])
		if op == "" {
			return nil, fmt.Errorf("missing operator in filter: %s", part)
		}
		if op != string(Eq) {
			return nil, fmt.Errorf("unsupported operator: %s", op)
		}
		value := parseFilterValue(strings.Join(fields[2:], " "))

		filters = append(filters, Filter{
			Field: field,
			Op:    FilterOp(op),
			Value: value,
		})
	}

	return filters, nil
}

func parseFilterValue(value string) any {
	if value == "true" {
		return true
	}

	if value == "false" {
		return false
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strings.Trim(value, "\"")
	}

	return value
}

type Sort struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending"` // ?sort=-type
}

func parseSorting(r *http.Request) ([]Sort, error) {
	query := r.URL.Query()
	sortString := query.Get("sort")
	if sortString == "" {
		return nil, nil
	}

	sortFields := strings.Split(sortString, ",")
	var sorters []Sort

	for _, field := range sortFields {
		ascending := true
		if field == "" {
			continue
		}
		if strings.HasPrefix(field, "-") {
			ascending = false
			field = strings.TrimPrefix(field, "-")
		}
		if field == "" {
			return nil, fmt.Errorf("empty sort field")
		}
		sorters = append(sorters, Sort{
			Field:     field,
			Ascending: ascending,
		})
	}

	return sorters, nil
}
