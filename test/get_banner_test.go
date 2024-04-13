package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mBayzigitov/dynamic-content-service/internal/dto"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (suite *BannerHandlerSuite) TestUserUnathorized() {
	testCases := []struct {
		name           string
		tagID          int
		featureID      int
		token          string
		expectedStatus int
	}{
		{
			name:           "ExpectUserUnathorized",
			tagID:          1,
			featureID:      10,
			expectedStatus: http.StatusUnauthorized,
			token:          "uup_29132993",
		},
		{
			name:           "ExpectAdminUnathorized",
			tagID:          2,
			featureID:      15,
			expectedStatus: http.StatusUnauthorized,
			token:          "uap_319299139",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			url := fmt.Sprintf(
				"/api/v1/user_banner?tag_id=%d&feature_id=%d",
				tc.tagID,
				tc.featureID,
			)
			req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
			suite.NoError(err, "failed to create request")

			req.Header.Set("X-Access-Token", tc.token)

			rec := httptest.NewRecorder()
			suite.router.ServeHTTP(rec, req)

			// check response status code
			suite.Equal(tc.expectedStatus, rec.Code, "unexpected status code")
		})
	}
}

func (suite *BannerHandlerSuite) TestAccessRestricted() {
	testCases := []struct {
		name           string
		tagID          int
		featureID      int
		token          string
		expectedStatus int
	}{
		{
			name:           "ExpectAccessRestricted",
			tagID:          1,
			featureID:      12,
			expectedStatus: http.StatusForbidden,
			token:          "invalid_token",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			url := fmt.Sprintf(
				"/api/v1/user_banner?tag_id=%d&feature_id=%d",
				tc.tagID,
				tc.featureID,
			)
			req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
			suite.NoError(err, "failed to create request")

			req.Header.Set("X-Access-Token", tc.token)

			rec := httptest.NewRecorder()
			suite.router.ServeHTTP(rec, req)

			// check response status code
			suite.Equal(tc.expectedStatus, rec.Code, "unexpected status code")
		})
	}
}

func (suite *BannerHandlerSuite) TestBannerNotFound() {
	testCases := []struct {
		name           string
		tagID          int
		featureID      int
		expectedStatus int
	}{
		{
			name:           "BannerNotFound1",
			tagID:          150,
			featureID:      123,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "BannerNotFound2",
			tagID:          200,
			featureID:      456,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "BannerNotFound3",
			tagID:          250,
			featureID:      789,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			url := fmt.Sprintf(
				"/api/v1/user_banner?tag_id=%d&feature_id=%d",
				tc.tagID,
				tc.featureID,
			)
			req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
			suite.NoError(err, "failed to create request")

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Access-Token", "aup_3101020")

			rec := httptest.NewRecorder()
			suite.router.ServeHTTP(rec, req)

			suite.Equal(tc.expectedStatus, rec.Code, "unexpected status code")
		})
	}
}

func (suite *BannerHandlerSuite) TestInvalidData() {
	testCases := []struct {
		name               string
		query              string
		expectedErrMessage string
		expectedStatus     int
	}{
		{
			name:               "ExpectInvalidFeatureId",
			query:              "/api/v1/user_banner?tag_id=2&feature_id=blablabla",
			expectedErrMessage: "Некорректное значение feature_id",
			expectedStatus:     400,
		},
		{
			name:               "ExpectInvalidTagId",
			query:              "/api/v1/user_banner?tag_id=blablabla&feature_id=11",
			expectedErrMessage: "Некорректное значение tag_id",
			expectedStatus:     400,
		},
		{
			name:               "ExpectInvalidUseLastRevision",
			query:              "/api/v1/user_banner?tag_id=2&feature_id=3&use_last_revision=blablabla",
			expectedErrMessage: "Некорректное значение use_last_revision",
			expectedStatus:     400,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req, err := http.NewRequest("GET", tc.query, bytes.NewBuffer([]byte("")))
			suite.NoError(err, "failed to create request")

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Access-Token", "aup_3101020")

			rec := httptest.NewRecorder()
			suite.router.ServeHTTP(rec, req)

			suite.Equal(tc.expectedStatus, rec.Code, "unexpected status code")

			var responseBody serverr.ApiError
			err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
			suite.NoError(err, "failed to unmarshal response")

			suite.Equal(tc.expectedErrMessage, responseBody.ErrType, "unexpected error message")
		})
	}
}

func (suite *BannerHandlerSuite) TestOkRequest() {
	testCases := []struct {
		name           string
		query          string
		expectedBanner string
	}{
		{
			name:       "OkRequest1",
			query: "/api/v1/user_banner?tag_id=1&feature_id=2&use_last_revision=false",
			expectedBanner: `{"title":"some_title 2","description":"Description of Banner 2"}`,
		},
		{
			name:       "OkRequest2",
			query: "/api/v1/user_banner?tag_id=1&feature_id=4&use_last_revision=true",
			expectedBanner: `{"title":"some_title 4","description":"Description of Banner 4"}`,
		},
		{
			name:       "OkRequest3",
			query: "/api/v1/user_banner?tag_id=1&feature_id=7&use_last_revision=false",
			expectedBanner: `{"title":"some_title 7","description":"Description of Banner 7"}`,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req, err := http.NewRequest("GET", tc.query, bytes.NewBuffer([]byte("")))
			suite.NoError(err, "failed to create request")

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Access-Token", "aap123123")

			rec := httptest.NewRecorder()
			suite.router.ServeHTTP(rec, req)

			suite.Equal(200, rec.Code, "unexpected status code")

			var responseBody dto.GetBannerResponseDto
			err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
			suite.NoError(err, "failed to unmarshal response")

			suite.Equal(tc.expectedBanner, string(responseBody.Content), "unexpected banner body")
		})
	}
}

func TestBannerHandlerSuite(t *testing.T) {
	suite.Run(t, new(BannerHandlerSuite))
}
