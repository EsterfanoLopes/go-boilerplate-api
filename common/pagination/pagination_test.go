package pagination

import (
	"fmt"
	"go-boilerplate/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name          string
		from          int
		size          int
		expectedError string
	}{
		{
			name: "valid pagination",
			from: 0,
			size: 10,
		},
		{
			name:          "invalid size",
			from:          0,
			size:          0,
			expectedError: "size: cannot be blank.",
		},
		{
			name:          "negative size",
			from:          0,
			size:          -1,
			expectedError: "size: must be no less than 1.",
		},
		{
			name:          "too big size",
			from:          0,
			size:          50,
			expectedError: "size: must be no greater than 30.",
		},

		{
			name:          "invalid from",
			from:          -1,
			size:          1,
			expectedError: "from: must be no less than 0.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.from, tc.size)
			test.AssertError(t, err, tc.expectedError)
		})
	}
}

func TestPaginateQuery(t *testing.T) {
	p, _ := New(90, 30)
	query := "select * from xablau where id = $1"
	paginatedQuery := p.PaginateQuery(query)
	if paginatedQuery != "select * from xablau where id = $1 LIMIT 30 OFFSET 90" {
		t.Errorf("error on generated paginated query %s", paginatedQuery)
		return
	}
}

func TestFromRequest(t *testing.T) {
	testCases := []struct {
		name          string
		from          string
		size          string
		expected      Pagination
		expectedError string
	}{
		{
			name: "valid request",
			from: "10",
			size: "20",
			expected: Pagination{
				from: 10,
				size: 20,
			},
		},
		{
			name: "valid request with default values",
			expected: Pagination{
				from: 0,
				size: 10,
			},
		},
		{
			name:          "invalid size",
			from:          "0",
			size:          "abc",
			expectedError: `strconv.Atoi: parsing "abc": invalid syntax`,
		},
		{
			name:          "invalid from",
			from:          "abc",
			size:          "1",
			expectedError: `strconv.Atoi: parsing "abc": invalid syntax`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/xablau?from=%s&size=%s", tc.from, tc.size), nil)
			p, err := FromRequest(request)
			test.AssertError(t, err, tc.expectedError)

			if p.size != tc.expected.size || p.from != tc.expected.from {
				t.Errorf("unexptected pagination %+v", p)
			}
		})
	}
}

func TestGetPageNumber(t *testing.T) {
	testCases := []struct {
		name         string
		from         int
		size         int
		expectedPage int
	}{
		{
			name:         "when from is < than size",
			from:         2,
			size:         10,
			expectedPage: 0,
		},
		{
			name:         "when from is > than size",
			from:         20,
			size:         10,
			expectedPage: 2,
		},
		{
			name:         "when from is == size",
			from:         10,
			size:         10,
			expectedPage: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pagination, err := New(tc.from, tc.size)
			if err != nil {
				t.Errorf("error when try to create a pagination %s", err)
				return
			}
			page := pagination.GetPageNumber()
			if page != tc.expectedPage {
				t.Errorf("unexpected page number %d", page)
				return
			}
		})
	}
}
