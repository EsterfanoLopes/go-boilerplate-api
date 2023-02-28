package api_test

import (
	"fmt"
	common "go-boilerplate/common"
	"go-boilerplate/test"
	"net/http"
	"testing"

	"go-boilerplate/repository"
)

func TestCommentRoutes(t *testing.T) {
	date, _ := common.ToTime("2021-01-06T20:35:00-03:00")
	test.FreezeTime(t, date)

	nextID := repository.GetNextID(t, "comment")

	testCases := []test.APITestCase{
		{
			Name:    "v1 post comment",
			Route:   "http://localhost:9000/v1/comment",
			Method:  http.MethodPost,
			Status:  http.StatusCreated,
			Payload: `{"accountId": "34178e2a-b9be-48ef-bfb4-3973747ae257","advertiserId": "77e04ae6-c3dc-4a60-8b52-d1fc35d42098","description": "Pessoa foi visitar e ninguém viu, dessa vez não vamos dar vacilo","listingId": "2323232323","owner": {"accountId": "1071a242-5d3f-45e5-9a7a-b64b9ab68e98","name": "José Silva","email":"jose.silva@mailinator.com"},"type": "SCHEDULE"}`,
			Body:    fmt.Sprintf(`{"id":%d}`, nextID),
		},
		{
			Name:    "v1 put comment",
			Route:   fmt.Sprintf("http://localhost:9000/v1/comment/%d", nextID),
			Method:  http.MethodPut,
			Status:  http.StatusOK,
			Payload: `{"accountId": "34178e2a-b9be-48ef-bfb4-3973747ae257","advertiserId": "77e04ae6-c3dc-4a60-8b52-d1fc35d42098","description": "A pessoa tentou realizar a visita, mas não obteve atendimento, estarei enviando um presente para ela.","listingId": "2323232323","owner": {"accountId": "1071a242-5d3f-45e5-9a7a-b64b9ab68e98","name": "José Silva","email":"jose.silva@mailinator.com"},"type": "SCHEDULE"}`,
			Body:    fmt.Sprintf(`{"id":%d}`, nextID),
		},
		{
			Name:   "v1 get comment",
			Route:  fmt.Sprintf("http://localhost:9000/v1/comment/%d", nextID),
			Method: http.MethodGet,
			Status: http.StatusOK,
			Body: fmt.Sprintf(
				`{"id":%d,"type":"SCHEDULE","description":"A pessoa tentou realizar a visita, mas não obteve atendimento, estarei enviando um presente para ela.","advertiserId":"77e04ae6-c3dc-4a60-8b52-d1fc35d42098","accountId":"34178e2a-b9be-48ef-bfb4-3973747ae257","listingId":"2323232323","updated":true,"owner":{"name":"José Silva","email":"jose.silva@mailinator.com","accountId":"1071a242-5d3f-45e5-9a7a-b64b9ab68e98"},"createdAt":"2021-01-06T20:35:00-03:00","updatedAt":"2021-01-06T20:35:00-03:00"}`,
				nextID,
			),
		},
		{
			Name:   "v1 get comments",
			Route:  "http://localhost:9000/v1/comment?advertiserId=77e04ae6-c3dc-4a60-8b52-d1fc35d42098&accountId=34178e2a-b9be-48ef-bfb4-3973747ae257&listingId=2323232323",
			Method: http.MethodGet,
			Status: http.StatusOK,
			Body: fmt.Sprintf(
				`{"total":1,"results":[{"id":%d,"type":"SCHEDULE","description":"A pessoa tentou realizar a visita, mas não obteve atendimento, estarei enviando um presente para ela.","advertiserId":"77e04ae6-c3dc-4a60-8b52-d1fc35d42098","accountId":"34178e2a-b9be-48ef-bfb4-3973747ae257","listingId":"2323232323","updated":true,"owner":{"name":"José Silva","email":"jose.silva@mailinator.com","accountId":"1071a242-5d3f-45e5-9a7a-b64b9ab68e98"},"createdAt":"2021-01-06T20:35:00-03:00","updatedAt":"2021-01-06T20:35:00-03:00"}]}`,
				nextID,
			),
		},
		{
			Name:   "v1 delete comment",
			Route:  fmt.Sprintf("http://localhost:9000/v1/comment/%d", nextID),
			Method: http.MethodDelete,
			Status: http.StatusNoContent,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, tc.Run)
	}
}
