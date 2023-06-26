package util

import (
	"os"
	"strconv"
)

type Pagination struct {
	CurrentPage  int         `json:"currentPage"`
	PreviousPage interface{} `json:"previousPage"`
	NextPage     interface{} `json:"nextPage"`
	Total        int         `json:"total"`
	Limit        int         `json:"limit"`
	Data         interface{} `json:"data"`
}

type Paging struct {
	Page   int
	Limit  int
	Offset int
}

func NewPaging(qPage string, qLimit string) *Paging {
	page, _ := strconv.Atoi(qPage)
	limit, _ := strconv.Atoi(qLimit)

	// set page
	if page < 1 {
		page = 1
	}

	// set limit
	if limit < 1 {
		limit, _ = strconv.Atoi(os.Getenv("QUERY_LIMIT_DEFAULT"))
	}

	offset := (page - 1) * limit

	return &Paging{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func Paginate(paging *Paging, data interface{}, total int) *Pagination {

	// calculate previous page
	previousPage := calculatePreviousPage(paging.Page)

	var resPreviousPage interface{}
	if previousPage == -1 {
		resPreviousPage = nil
	} else {
		resPreviousPage = previousPage
	}

	// calculate next page
	nextPage := calculateNextPage(paging.Page, total, paging.Limit)

	var resNextPage interface{}
	if nextPage == -1 {
		resNextPage = nil
	} else {
		resNextPage = nextPage
	}

	return &Pagination{
		CurrentPage:  paging.Page,
		PreviousPage: resPreviousPage,
		NextPage:     resNextPage,
		Total:        total,
		Limit:        paging.Limit,
		Data:         data,
	}
}

func calculatePreviousPage(page int) int {
	previousPage := page - 1

	if previousPage < 1 {
		return -1
	}

	return previousPage
}

func calculateNextPage(page int, total int, limit int) int {
	leftover := float64(total) / float64(page*limit)

	if leftover <= 1 {
		return -1
	}

	return page + 1
}
