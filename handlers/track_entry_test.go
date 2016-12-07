package handlers

import (
	"github.com/ducktrack/wing/config"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var appConfig config.Config

func TestMain(m *testing.M) {
	appConfig = config.Config{
		Exporter: "file",
		FileExporter: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
	}

	os.Exit(m.Run())
}

func TestWhenRequestMethodOptions(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := &TrackEntryHandler{Config: &appConfig}

	req, _ := http.NewRequest("OPTIONS", "/", nil)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 200 {
		t.Errorf("handler returned wrong status code:\ngot:\n%v\n\nwant:\n%v\n", status, 200)
	}

	if rr.Body.String() != "" {
		t.Errorf(
			"handler returned unexpected body:\ngot:\n%v\n\nwant:\n%v\n",
			rr.Body.String(),
			"",
		)
	}
}

func TestWhenRequestMethodDifferentThanPost(t *testing.T) {
	for _, method := range []string{"GET", "PUT", "DELETE", "PATCH"} {
		rr := httptest.NewRecorder()
		handler := &TrackEntryHandler{Config: &appConfig}
		req, _ := http.NewRequest(method, "/", nil)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != 405 {
			t.Errorf("handler returned wrong status code:\ngot:\n%v\n\nwant:\n%v\n", status, 405)
		}

		expected := `{"message": "Method Not Allowed"}`
		if rr.Body.String() != expected {
			t.Errorf(
				"handler returned unexpected body:\ngot:\n%v\n\nwant:\n%v\n",
				rr.Body.String(),
				expected,
			)
		}
	}
}

func TestWhenJsonPayloadIsInvalid(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := &TrackEntryHandler{Config: &appConfig}

	req, _ := http.NewRequest("POST", "/", strings.NewReader("{")) // invalid
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 422 {
		t.Errorf("handler returned wrong status code:\ngot:\n%v\n\nwant:\n%v\n", status, 422)
	}

	expected := `{"message": "Invalid JSON payload"}`
	if rr.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body:\ngot:\n%v\n\nwant:\n%v\n",
			rr.Body.String(),
			expected,
		)
	}
}

func TestWhenBase64IsInvalid(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := &TrackEntryHandler{Config: &appConfig}

	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`{"created_at": 12345, "markup": "123"}`),
	)

	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 422 {
		t.Errorf("handler returned wrong status code:\ngot:\n%v\n\nwant:\n%v\n", status, 422)
	}

	expected := `{"message": "Invalid base64 payload"}`
	if rr.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body:\ngot:\n%v\n\nwant:\n%v\n",
			rr.Body.String(),
			expected,
		)
	}
}

func TestWhenItSavesTheRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := &TrackEntryHandler{Config: &appConfig}

	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`{"created_at": 1480979268, "markup": "PGh0bWw+PC9odG1sPg=="}`),
	)

	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 201 {
		t.Errorf("handler returned wrong status code:\ngot:\n%v\n\nwant:\n%v\n", status, 201)
	}

	if rr.Body.String() != "" {
		t.Errorf(
			"handler returned unexpected body:\ngot:\n%v\n\nwant:\n%v\n",
			rr.Body.String(),
			"",
		)
	}

	request := http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	_, err := request.Cookie(RECORD_ID_COOKIE_NAME)

	if err != nil {
		t.Errorf("expected handler to create '%s' cookie", RECORD_ID_COOKIE_NAME)
	}
}
