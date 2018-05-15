package utils

import (
	"math"
	"net/url"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Pagination struct {
	Query         *gorm.DB
	TotalEntities int        `json:"total_entities" `
	PerPage       int        `json:"per_page" `
	Path          string     `json:"path" `
	Page          int        `json:"page" `
	UrlQuery      url.Values `json:"url_query" `
	TotalPages    int        `json:"total_pages" `
}


func (p *Pagination) Paginate(page interface{}) *gorm.DB {
	switch pageVal := page.(type) {
	case string:
		pageNo, err := strconv.ParseInt(pageVal, 10, 64)
		if err != nil {
			return p.Query
		}
		p.Page = int(pageNo)

		p.Query.Count(&p.TotalEntities)
		if p.TotalEntities == 0 {
			return p.Query
		}

		p.TotalPages = int(math.Ceil(float64(p.TotalEntities) / float64(p.PerPage)))

		if !(p.Page > 0 && p.Page <= p.TotalPages) {
			p.Page = 1
		}

		query := p.Query.Offset((p.Page - 1) * p.PerPage).Limit(p.PerPage)

		return query

	case int:
		p.Page = pageVal

		p.Query.Count(&p.TotalEntities)
		if p.TotalEntities == 0 {
			return p.Query
		}

		p.TotalPages = int(math.Ceil(float64(p.TotalEntities) / float64(p.PerPage)))

		if !(p.Page > 0 && p.Page <= p.TotalPages) {
			p.Page = 1
		}

		query := p.Query.Offset((p.Page - 1) * p.PerPage).Limit(p.PerPage)

		return query

	default:
		return p.Query
	}

}
