package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asiafrolova/Final_task/orkestrator_service/internal/handlers"
	"github.com/asiafrolova/Final_task/orkestrator_service/internal/logger"
	"github.com/asiafrolova/Final_task/orkestrator_service/pkg/orkestrator"
)

func TestAddExpressionGetIDHandlerOK(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{"1+2"},
		{"1+2+3"},
		{"((1+1)-(2-1))"},
		{"3.14-1"},
		{"1/1"},
	}
	for _, tc := range testCases {
		logger.Init()
		postBody := map[string]interface{}{
			"expression": tc.input,
		}
		body, _ := json.Marshal(postBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handlers.AddExpressionsHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected error with expression %v, statuscode: %v", tc.input, resp.StatusCode)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Unexpected error with expression %v: %v", tc.input, err)
		}
		var respBody *map[string]string

		err = json.Unmarshal(body, &respBody)
		if err != nil {
			t.Fatalf("Unexpected error with expression %v : %v", tc.input, err)
		}
		id := ""
		for key, val := range *respBody {
			if key == "id" {
				id = fmt.Sprint(val)
				break
			}
		}
		if id == "" {
			t.Fatalf("Invalid id with expression %v", tc.input)
		}
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/expressions?id=%s", id), nil)
		w = httptest.NewRecorder()
		handlers.GetExpressionByIDHandler(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected error with expression %v, statuscode: %v", tc.input, resp.StatusCode)
		}

	}

}
func TestAddExpressionGetIDHandlerError(t *testing.T) {
	testCases := []struct {
		input string
		want  error
	}{
		{"1+2)", orkestrator.ErrInvalidExpression},
		{"1+2**3", orkestrator.ErrInvalidExpression},
		{"((1+101))))-(2-1))", orkestrator.ErrInvalidExpression},
		{"3.14-1-~3", orkestrator.ErrInvalidExpression},
		{"1/1a", orkestrator.ErrInvalidExpression},
	}
	for _, tc := range testCases {
		logger.Init()
		postBody := map[string]interface{}{
			"expression": tc.input,
		}
		body, _ := json.Marshal(postBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handlers.AddExpressionsHandler(w, req)
		resp := w.Result()
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("Invalid status with expression %v, statuscode: %v", tc.input, resp.StatusCode)
		}

		id := ""

		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/expressions?id=%s", id), nil)
		w = httptest.NewRecorder()
		handlers.GetExpressionByIDHandler(w, req)
		resp = w.Result()
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("Invalid status with expression %v, statuscode: %v", tc.input, resp.StatusCode)
		}

	}

}
