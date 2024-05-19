package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"snippetbox.doichevkostia.dev/internal/assert"
	"snippetbox.doichevkostia.dev/internal/models/mocks"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	// start a test server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "pong")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  fmt.Sprintf("/snippet/view/%s", mocks.SnippetID),
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  fmt.Sprintf("/snippet/view/%s", uuid.New()),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
