package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const responseBody = "helloworld"

// Test case a) No content-length. Status 200
func TestNoContentLength200(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	printResponseDetails(t, "a) No content-length, Status 200", resp, body)
}

// Test case b) No content-length. Status 204
func TestNoContentLength204(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
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
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	printResponseDetails(t, "b) No content-length, Status 204", resp, body)
}

// Test case c) Content-Length: 2 Status 204
func TestContentLength2Status204(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "2")
		w.WriteHeader(http.StatusNoContent)
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
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	printResponseDetails(t, "c) Content-Length: 2, Status 204", resp, body)
}

// Test case d) Content-Length: 11 Status 204
func TestContentLength11Status204(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "11")
		w.WriteHeader(http.StatusNoContent)
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
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	printResponseDetails(t, "d) Content-Length: 11, Status 204", resp, body)
}

// Test case e) Content-Length: 2 Status 200
func TestContentLength2Status200(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "2")
		w.WriteHeader(http.StatusOK)
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
	// Note: We expect an error here when Content-Length is less than actual body
	if err != nil {
		t.Logf("Error reading response body: %v", err)
	}

	printResponseDetails(t, "e) Content-Length: 2, Status 200", resp, body)
}

// Helper function to print response details
func printResponseDetails(t *testing.T, testCase string, resp *http.Response, body []byte) {
	t.Logf("\n========================================")
	t.Logf("Test Case: %s", testCase)
	t.Logf("========================================")
	t.Logf("Status Code: %d", resp.StatusCode)
	t.Logf("Status: %s", resp.Status)
	t.Logf("Content-Length Header: %s", resp.Header.Get("Content-Length"))
	t.Logf("Transfer-Encoding: %v", resp.TransferEncoding)
	t.Logf("Response Body Length: %d", len(body))
	t.Logf("Response Body: %q", string(body))
	t.Logf("Expected Body: %q", responseBody)
	t.Logf("Expected Body Length: %d", len(responseBody))
	
	// Check for mismatches
	contentLengthHeader := resp.Header.Get("Content-Length")
	if contentLengthHeader != "" {
		var clValue int
		fmt.Sscanf(contentLengthHeader, "%d", &clValue)
		if clValue != len(body) {
			t.Logf("⚠️  MISMATCH: Content-Length header (%d) != actual body length (%d)", clValue, len(body))
		}
		if clValue != len(responseBody) {
			t.Logf("⚠️  MISMATCH: Content-Length header (%d) != written body length (%d)", clValue, len(responseBody))
		}
	}
	
	if !bytes.Equal(body, []byte(responseBody)) {
		t.Logf("⚠️  Body content differs from what was written")
	}
	
	t.Logf("========================================\n")
}
