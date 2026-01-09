package utils

import (
	"strconv"
	"strings"

	"snack-store-api/internal/model"
)

func ParsePagination(
	pageParam string,
	pageSizeParam string,
	defaultPage int,
	defaultPageSize int,
) (int, int, error) {
	page := defaultPage
	pageSize := defaultPageSize

	if strings.TrimSpace(pageParam) != "" {
		parsed, err := strconv.Atoi(strings.TrimSpace(pageParam))
		if err != nil {
			return 0, 0, err
		}
		page = parsed
	}

	if strings.TrimSpace(pageSizeParam) != "" {
		parsed, err := strconv.Atoi(strings.TrimSpace(pageSizeParam))
		if err != nil {
			return 0, 0, err
		}
		pageSize = parsed
	}

	return page, pageSize, nil
}

func BuildPageMetadata(page int, pageSize int, totalItem int64) model.PageMetadata {
	totalPage := int64(0)
	if pageSize > 0 && totalItem > 0 {
		totalPage = (totalItem + int64(pageSize) - 1) / int64(pageSize)
	}

	hasNext := totalPage > 0 && int64(page) < totalPage
	hasPrevious := totalPage > 0 && page > 1

	return model.PageMetadata{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItem:   totalItem,
		TotalPage:   totalPage,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}
}
