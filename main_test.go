package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSmoke(t *testing.T) {
	t.Parallel()

	router := NewRouter()

	// Test parameters
	var url = "/"

	req := httptest.NewRequest("GET", url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var testResp TestResponse

	err := json.Unmarshal([]byte(rec.Body.String()), &testResp)
	if err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, testResp.Return, "hoge")
	}
}
