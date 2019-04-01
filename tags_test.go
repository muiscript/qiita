package qiita

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"strings"
	"testing"
)

func TestClient_GetTag(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "tags", "GetTag")

	tests := []struct {
		desc       string
		inputTagID string

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod         string
		expectedRequestPath    string
		expectedRawQuery       string
		expectedErrString      string
		expectedID             string
		expectedIconURL        string
		expectedItemsCount     int
		expectedFollowersCount int
	}{
		{
			desc:       "success",
			inputTagID: "react",

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:         http.MethodGet,
			expectedRequestPath:    "/tags/react",
			expectedID:             "React",
			expectedIconURL:        "https://s3-ap-northeast-1.amazonaws.com/qiita-tag-image/c4d0439277f132acce23de37f694617b95af5475/medium.jpg?1513495262",
			expectedItemsCount:     2693,
			expectedFollowersCount: 2403,
		},
		{
			desc:       "failure-not_exist",
			inputTagID: "nonexistent",

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/tags/nonexistent",
			expectedErrString:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			tag, err := cli.GetTag(context.Background(), tt.inputTagID)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedID, tag.ID)
				assert.Equal(t, tt.expectedIconURL, tag.IconURL)
				assert.Equal(t, tt.expectedItemsCount, tag.ItemsCount)
				assert.Equal(t, tt.expectedFollowersCount, tag.FollowersCount)
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}

		})
	}
}

func TestClient_GetTags(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "tags", "GetTags")

	tests := []struct {
		desc         string
		inputPage    int
		inputPerPage int
		inputSort    Sort

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
		expectedPage        int
		expectedPerPage     int
		expectedFirstPage   int
		expectedLastPage    int
		expectedTotalCount  int
		expectedTagsLen     int
	}{
		{
			desc:         "success",
			inputPage:    3,
			inputPerPage: 2,
			inputSort:    Count,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/tags",
			expectedRawQuery:    "page=3&per_page=2&sort=count",
			expectedPage:        3,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    100,
			expectedTotalCount:  85414,
			expectedTagsLen:     2,
		},
		{
			desc:         "failure-out_of_range",
			inputPage:    101,
			inputPerPage: 2,
			inputSort:    Count,

			expectedErrString: "page parameter should be",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			tagsResp, err := cli.GetTags(context.Background(), tt.inputPage, tt.inputPerPage, tt.inputSort)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, tagsResp.Page)
				assert.Equal(t, tt.expectedPerPage, tagsResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, tagsResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, tagsResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, tagsResp.TotalCount)
				assert.Equal(t, tt.expectedTagsLen, len(tagsResp.Tags))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}

func TestClient_GetTagItems(t *testing.T) {
	mockFilesBaseDir := path.Join("testdata", "responses", "tags", "GetTagItems")

	tests := []struct {
		desc         string
		inputTagID   string
		inputPage    int
		inputPerPage int

		mockResponseHeaderFile string
		mockResponseBodyFile   string

		expectedMethod      string
		expectedRequestPath string
		expectedRawQuery    string
		expectedErrString   string
		expectedPage        int
		expectedPerPage     int
		expectedFirstPage   int
		expectedLastPage    int
		expectedTotalCount  int
		expectedItemsLen    int
	}{
		{
			desc:         "success",
			inputTagID:   "react",
			inputPage:    3,
			inputPerPage: 2,

			mockResponseHeaderFile: "success-header",
			mockResponseBodyFile:   "success-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/tags/react/items",
			expectedRawQuery:    "page=3&per_page=2",
			expectedPage:        3,
			expectedPerPage:     2,
			expectedFirstPage:   1,
			expectedLastPage:    100,
			expectedTotalCount:  2694,
			expectedItemsLen:    2,
		},
		{
			desc:         "failure-not_exist",
			inputTagID:   "nonexistent",
			inputPage:    3,
			inputPerPage: 2,

			mockResponseHeaderFile: "not_exist-header",
			mockResponseBodyFile:   "not_exist-body",

			expectedMethod:      http.MethodGet,
			expectedRequestPath: "/tags/nonexistent/items",
			expectedRawQuery:    "page=3&per_page=2",
			expectedErrString:   "not found",
		},
		{
			desc:         "failure-out_of_range",
			inputTagID:   "react",
			inputPage:    101,
			inputPerPage: 2,

			expectedErrString: "page parameter should be",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, teardown := setup(t, mockFilesBaseDir, tt.mockResponseHeaderFile, tt.mockResponseBodyFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
			defer teardown()

			itemsResp, err := cli.GetTagItems(context.Background(), tt.inputTagID, tt.inputPage, tt.inputPerPage)
			if tt.expectedErrString == "" {
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.expectedPage, itemsResp.Page)
				assert.Equal(t, tt.expectedPerPage, itemsResp.PerPage)
				assert.Equal(t, tt.expectedFirstPage, itemsResp.FirstPage)
				assert.Equal(t, tt.expectedLastPage, itemsResp.LastPage)
				assert.Equal(t, tt.expectedTotalCount, itemsResp.TotalCount)
				assert.Equal(t, tt.expectedItemsLen, len(itemsResp.Items))
			} else {
				if !assert.NotNil(t, err) {
					t.FailNow()
				}

				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString), fmt.Sprintf("'%s' should contain '%s'", err.Error(), tt.expectedErrString))
			}
		})
	}
}