package response

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestError(t *testing.T) {
    recorder := httptest.NewRecorder()

    Error(recorder, http.StatusNotFound, "user not found")

    if recorder.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", recorder.Code)
    }

    if recorder.Body.String() == "" {
        t.Fatal("expected JSON response body")
    }
}