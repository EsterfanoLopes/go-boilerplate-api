package repository_test

import (
	"go-boilerplate/repository"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetOrderByStatusClause(t *testing.T) {
	testCases := []struct {
		name     string
		alias    string
		statuses [][]string
		expected string
	}{
		{
			name:  "status order by clause",
			alias: "credit_analysis",
			statuses: [][]string{
				{
					"NEW",
				},
				{
					"OLD",
				},
				{
					"ANY",
					"LIKE_ANY",
				},
			},
			expected: "case credit_analysis.last_status when 'NEW' then 0 when 'OLD' then 1 when 'ANY' then 2 when 'LIKE_ANY' then 2 else 99 end",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := repository.GetOrderByStatusClause(tc.alias, tc.statuses)
			if diff := cmp.Diff(result, tc.expected); diff != "" {
				t.Errorf("unexpected order by clause %s", diff)
				return
			}
		})
	}
}
