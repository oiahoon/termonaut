package unit

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/oiahoon/termonaut/internal/network"
)

func TestEnhancedClientRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Fail first 2 attempts
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Succeed on 3rd attempt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := network.NewEnhancedClient(5 * time.Second)
	
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestEnhancedClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with short timeout
	client := network.NewEnhancedClient(500 * time.Millisecond)
	
	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !network.IsNetworkError(err) {
		t.Errorf("Expected network error, got: %v", err)
	}
}

func TestEnhancedClientWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := network.NewEnhancedClient(5 * time.Second)
	
	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := client.GetWithContext(ctx, server.URL)
	if err == nil {
		t.Error("Expected context timeout error, got nil")
	}
}

func TestCircuitBreakerOpen(t *testing.T) {
	failureCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failureCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := network.NewEnhancedClient(1 * time.Second)
	
	// Make enough requests to open the circuit breaker
	for i := 0; i < 10; i++ {
		client.Get(server.URL)
	}

	// Circuit breaker should now be open, preventing further requests
	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("Expected circuit breaker error, got nil")
	}

	if !fmt.Sprintf("%v", err) != "circuit breaker is open" {
		t.Logf("Circuit breaker error: %v", err)
	}
}

func TestSafeReadBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response body"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	body, err := network.SafeReadBody(resp)
	if err != nil {
		t.Fatalf("SafeReadBody failed: %v", err)
	}

	expected := "test response body"
	if string(body) != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, string(body))
	}
}

func TestSafeReadBodyNilResponse(t *testing.T) {
	_, err := network.SafeReadBody(nil)
	if err == nil {
		t.Error("Expected error for nil response, got nil")
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"generic error", fmt.Errorf("generic error"), false},
		{"context deadline", context.DeadlineExceeded, false},
		{"context canceled", context.Canceled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := network.IsNetworkError(tt.err)
			if result != tt.expected {
				t.Errorf("IsNetworkError(%v) = %v, expected %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"context deadline", context.DeadlineExceeded, false},
		{"context canceled", context.Canceled, false},
		{"generic error", fmt.Errorf("generic error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := network.IsRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("IsRetryableError(%v) = %v, expected %v", tt.err, result, tt.expected)
			}
		})
	}
}
