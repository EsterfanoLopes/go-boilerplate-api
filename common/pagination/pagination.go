// Package pagination generic functions to build pagination in responses
package pagination

import (
	"fmt"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const defaultSize = 10

// Pagination component
type Pagination struct {
	from int
	size int
}

// GetFrom returns the current from value
func (p Pagination) GetFrom() int {
	return p.from
}

// GetSize returns the current size value
func (p Pagination) GetSize() int {
	return p.size
}

// GetPageNumber returns the current page number
func (p Pagination) GetPageNumber() int {
	if p.from > 0 && p.size > 0 && p.from >= p.size {
		return p.from / p.size
	}
	return 0
}

// Response of pagination
type Response struct {
	Total   int         `json:"total"`
	Results interface{} `json:"results"`
}

// GetResponse returns pagination paginated response
func (p Pagination) GetResponse(total int, results interface{}) Response {
	return Response{
		Total:   total,
		Results: results,
	}
}

// New creates a pagination complex struct instance
func New(from, size int) (Pagination, error) {
	p := Pagination{
		from: from,
		size: size,
	}

	if err := p.validate(); err != nil {
		return Pagination{}, err
	}

	return p, nil
}

// FromRequest parse pagination data from HTTP request
func FromRequest(r *http.Request) (Pagination, error) {
	fromParameter := r.URL.Query().Get("from")
	sizeParameter := r.URL.Query().Get("size")

	from := 0
	if fromParameter != "" {
		f, err := strconv.Atoi(fromParameter)
		if err != nil {
			return Pagination{}, err
		}
		from = f
	}

	size := defaultSize
	if sizeParameter != "" {
		s, err := strconv.Atoi(sizeParameter)
		if err != nil {
			return Pagination{}, err
		}
		size = s
	}

	p := Pagination{
		from: from,
		size: size,
	}
	if err := p.validate(); err != nil {
		return Pagination{}, err
	}

	return p, nil
}

// PaginateQuery return a query with paginate string
func (p Pagination) PaginateQuery(sqlQuery string) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", sqlQuery, p.size, p.from)
}

func (p Pagination) validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.from, validation.Min(0)),
		validation.Field(&p.size, validation.Required, validation.Min(1), validation.Max(30)),
	)
}
