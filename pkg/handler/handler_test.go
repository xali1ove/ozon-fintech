package handler

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	ozon_fintech "ozon-fintech"
	mock_service "ozon-fintech/pkg/service/mocks"
	"testing"
)

func TestHandler_createShortURL(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices)

	testTable := []struct {
		name               string
		input              string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:  "OK",
			input: `{"base_url":"https://ozon.ru/sas"}`,
			mockBehavior: func(s *mock_service.MockServices) {
				s.EXPECT().CreateShortURL(&ozon_fintech.Link{BaseURL: "https://ozon.ru/sas"}).Return("token", nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "error_validation",
			input:              `{"base_url":"https://ozon"}`,
			mockBehavior:       func(s *mock_service.MockServices) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "error_internal",
			input: `{"base_url":"https://ozon.ru/sas"}`,
			mockBehavior: func(repos *mock_service.MockServices) {
				repos.EXPECT().CreateShortURL(&ozon_fintech.Link{BaseURL: "https://ozon.ru/sas"}).Return("", fmt.Errorf("some error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	c := gomock.NewController(t)
	defer c.Finish()

	service := mock_service.NewMockServices(c)
	handlers := NewHandler(service)

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(service)

			r := echo.New()
			handlers.InitRotes(r)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/tokens",
				bytes.NewBufferString(tc.input),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
}

func TestHandler_getBaseURL(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices)

	testTable := []struct {
		name               string
		input              string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:  "OK",
			input: "zxcl012_yz",
			mockBehavior: func(s *mock_service.MockServices) {
				s.EXPECT().
					GetBaseURL(&ozon_fintech.Link{Token: "zxcl012_yz"}).Return("https://ozon.ru", nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "error_validation",
			input:              "z!&l012_yz",
			mockBehavior:       func(s *mock_service.MockServices) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "error_internal",
			input: "zxcl012_yz",
			mockBehavior: func(s *mock_service.MockServices) {
				s.EXPECT().
					GetBaseURL(&ozon_fintech.Link{Token: "zxcl012_yz"}).Return("", fmt.Errorf("some error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:  "error_not_such_baseURL",
			input: "zxcl012_yz",
			mockBehavior: func(s *mock_service.MockServices) {
				s.EXPECT().
					GetBaseURL(&ozon_fintech.Link{Token: "zxcl012_yz"}).Return("", nil)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	c := gomock.NewController(t)
	defer c.Finish()

	service := mock_service.NewMockServices(c)
	handlers := NewHandler(service)

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(service)

			r := echo.New()
			handlers.InitRotes(r)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/tokens/"+tc.input,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
}
