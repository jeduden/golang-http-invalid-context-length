package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const responseBody = "helloworld"

// Helper function to run test cases and perform assertions
func runTest(t *testing.T, testCase string, statusCode int, contentLength string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contentLength != "" {
			w.Header().Set("Content-Length", contentLength)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	// Note: Error reading body is acceptable when Content-Length < actual data
	if err != nil {
		t.Logf("Error reading response body: %v", err)
	}

	// Print response details for visibility
	printResponseDetails(t, testCase, resp, body, contentLength)

	// ASSERTION 1: 204 No Content responses MUST have empty body
	if resp.StatusCode == http.StatusNoContent && len(body) > 0 {
		t.Errorf("FAIL: Status 204 must return empty body, but got %d bytes", len(body))
	}

	// ASSERTION 2: Body length must never exceed Content-Length header
	contentLengthHeader := resp.Header.Get("Content-Length")
	if contentLengthHeader != "" {
		expectedLen, err := strconv.Atoi(contentLengthHeader)
		if err == nil && len(body) > expectedLen {
			t.Errorf("FAIL: Body length (%d) exceeds Content-Length header (%d)", len(body), expectedLen)
		}
	}
}

// Helper function to print response details
func printResponseDetails(t *testing.T, testCase string, resp *http.Response, body []byte, requestedContentLength string) {
	t.Logf("\n========================================")
	t.Logf("Test Case: %s", testCase)
	t.Logf("========================================")
	t.Logf("Handler wrote: %q (%d bytes)", responseBody, len(responseBody))
	if requestedContentLength != "" {
		t.Logf("Handler set Content-Length: %s", requestedContentLength)
	} else {
		t.Logf("Handler set Content-Length: (not set)")
	}
	t.Logf("Status Code: %d", resp.StatusCode)
	t.Logf("Status: %s", resp.Status)
	t.Logf("Response Content-Length Header: %s", resp.Header.Get("Content-Length"))
	t.Logf("Response Transfer-Encoding: %v", resp.TransferEncoding)
	t.Logf("Client received body: %q (%d bytes)", string(body), len(body))

	// Check for Content-Length constraints
	contentLengthHeader := resp.Header.Get("Content-Length")
	if contentLengthHeader != "" {
		expectedLen, err := strconv.Atoi(contentLengthHeader)
		if err == nil {
			if len(body) > expectedLen {
				t.Logf("❌ VIOLATION: Body (%d bytes) > Content-Length (%d)", len(body), expectedLen)
			} else if len(body) <= expectedLen {
				t.Logf("✓ Body (%d bytes) <= Content-Length (%d)", len(body), expectedLen)
			}
		}
	}

	// Check for 204 No Content constraint
	if resp.StatusCode == http.StatusNoContent {
		if len(body) == 0 {
			t.Logf("✓ Status 204 correctly has empty body")
		} else {
			t.Logf("❌ VIOLATION: Status 204 has non-empty body (%d bytes)", len(body))
		}
	}

	t.Logf("========================================\n")
}

// Test case a) No content-length. Status 200
func TestNoContentLength200(t *testing.T) {
	runTest(t, "a) No content-length, Status 200", http.StatusOK, "")
}

// Test case b) No content-length. Status 204
func TestNoContentLength204(t *testing.T) {
	runTest(t, "b) No content-length, Status 204", http.StatusNoContent, "")
}

// Test case c) Content-Length: 2 Status 204
func TestContentLength2Status204(t *testing.T) {
	runTest(t, "c) Content-Length: 2, Status 204", http.StatusNoContent, "2")
}

// Test case d) Content-Length: 11 Status 204
func TestContentLength11Status204(t *testing.T) {
	runTest(t, "d) Content-Length: 11, Status 204", http.StatusNoContent, "11")
}

// Test case e) Content-Length: 2 Status 200
func TestContentLength2Status200(t *testing.T) {
	runTest(t, "e) Content-Length: 2, Status 200", http.StatusOK, "2")
}
