package domain

import (
	"github.com/online-store/internal/domain"
	"strconv"
)

const (
	DEFAULT_PAGESIZE = 10
	MAX_PAGESIZE     = 100
	DEFAULT_PAGE     = 1
)

func PageAndPageSizeValidation(pageSizeQuery string, pageQuery string) (pageSize int, page int, err error) {
	parse, err := strconv.Atoi(pageSizeQuery)
	if err == nil {
		pageSize = parse
		if parse == 0 {
			pageSize = DEFAULT_PAGESIZE
		}
		// dafault value maximum pageSize = 100
		if pageSize > MAX_PAGESIZE {
			pageSize = MAX_PAGESIZE
		}
		if parse < 0 {
			return 0, 0, domain.ErrInvalidUrlQueryParam
		}
	}
	if err != nil && pageSizeQuery == "" { // bypassing no query param
		pageSize = DEFAULT_PAGESIZE
	} else if err != nil {
		return 0, 0, domain.ErrInvalidUrlQueryParam
	}
	parse, err = strconv.Atoi(pageQuery)
	if err == nil {
		page = parse
	}
	if parse < 0 {
		return 0, 0, domain.ErrInvalidUrlQueryParam
	}
	if err != nil && pageQuery == "" { // bypassing no query param
		page = DEFAULT_PAGE
	} else if err != nil {
		return 0, 0, domain.ErrInvalidUrlQueryParam
	}
	err = nil
	return
}
