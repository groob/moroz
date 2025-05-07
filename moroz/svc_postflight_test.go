package moroz

import (
	"bytes"
	"compress/zlib"
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestDecodePostflightRequest(t *testing.T) {
	tests := []struct {
		name          string
		inputJSON     string
		expectedID    string
		expectedError bool
	}{
		{
			name:          "Valid JSON",
			inputJSON:     `{"rules_received":1,"rules_processed":1,"machine_id":"serial_number"}`,
			expectedID:    "serial_number",
			expectedError: false,
		},
		{
			name:          "Valid JSON with Strings",
			inputJSON:     `{"rules_received":"1","rules_processed":"1","machine_id":"serial_number"}`,
			expectedID:    "serial_number",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			zw := zlib.NewWriter(&buf)
			_, err := zw.Write([]byte(tt.inputJSON))
			if err != nil {
				t.Fatalf("failed to write compressed data: %v", err)
			}
			zw.Close()

			req, err := http.NewRequest("POST", "/v1/santa/postflight/serial_number", &buf)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			vars := map[string]string{"id": "serial_number"}
			req = mux.SetURLVars(req, vars)

			req.Header.Set("Content-Encoding", "deflate")

			result, err := decodePostflightRequest(context.Background(), req)
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			reqResult, ok := result.(postflightRequest)
			if !ok {
				t.Fatalf("expected postflightRequest, got %T", result)
			}

			if reqResult.MachineID != tt.expectedID {
				t.Errorf("expected MachineID %q, got %q", tt.expectedID, reqResult.MachineID)
			}
		})
	}
}
