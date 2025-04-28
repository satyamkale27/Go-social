package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"min=1,max=20"`
	Offset int      `schema:"offset" validate:"min=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// Parse limit with default value
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	} else {
		fq.Limit = 20 // Default value
	}

	// Parse offset with default value
	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	} else {
		fq.Offset = 0 // Default value
	}

	// Parse sort with default value
	if sort := qs.Get("sort"); sort != "" {
		fq.Sort = sort
	} else {
		fq.Sort = "desc" // Default value
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = ParseTime(since)
	}
	until := qs.Get("until")
	if since != "" {
		fq.Until = ParseTime(until)
	}

	return fq, nil
}

func ParseTime(s string) string {

	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)

}
