package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []error
}

// CircuitBreakerConfig defines circuit breaker behavior
type CircuitBreakerConfig struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
	MaxRequests      int
}

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	CircuitClosed CircuitBreakerState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config       CircuitBreakerConfig
	state        CircuitBreakerState
	failures     int
	lastFailTime time.Time
	requests     int
	mutex        sync.RWMutex
}

// EnhancedClient provides HTTP client with retry and circuit breaker
type EnhancedClient struct {
	httpClient     *http.Client
	retryConfig    RetryConfig
	circuitBreaker *CircuitBreaker
}

// NewEnhancedClient creates a new enhanced HTTP client
func NewEnhancedClient(timeout time.Duration) *EnhancedClient {
	return &EnhancedClient{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		retryConfig: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  100 * time.Millisecond,
			MaxDelay:      5 * time.Second,
			BackoffFactor: 2.0,
		},
		circuitBreaker: &CircuitBreaker{
			config: CircuitBreakerConfig{
				FailureThreshold: 5,
				RecoveryTimeout:  30 * time.Second,
				MaxRequests:      3,
			},
			state: CircuitClosed,
		},
	}
}

// Get performs a GET request with retry and circuit breaker
func (c *EnhancedClient) Get(url string) (*http.Response, error) {
	return c.DoWithRetry(func() (*http.Response, error) {
		return c.httpClient.Get(url)
	})
}

// GetWithContext performs a GET request with context
func (c *EnhancedClient) GetWithContext(ctx context.Context, url string) (*http.Response, error) {
	return c.DoWithRetry(func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		return c.httpClient.Do(req)
	})
}

// DoWithRetry executes a request with retry logic and circuit breaker
func (c *EnhancedClient) DoWithRetry(requestFunc func() (*http.Response, error)) (*http.Response, error) {
	// Check circuit breaker
	if !c.circuitBreaker.CanRequest() {
		return nil, fmt.Errorf("circuit breaker is open")
	}

	var lastErr error
	delay := c.retryConfig.InitialDelay

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * c.retryConfig.BackoffFactor)
			if delay > c.retryConfig.MaxDelay {
				delay = c.retryConfig.MaxDelay
			}
		}

		resp, err := requestFunc()
		if err == nil && resp.StatusCode < 500 {
			// Success or client error (don't retry client errors)
			c.circuitBreaker.RecordSuccess()
			return resp, nil
		}

		lastErr = err
		if err != nil {
			c.circuitBreaker.RecordFailure()
		} else {
			// Server error
			resp.Body.Close()
			lastErr = fmt.Errorf("server error: %d %s", resp.StatusCode, resp.Status)
			c.circuitBreaker.RecordFailure()
		}

		// Don't retry on last attempt
		if attempt == c.retryConfig.MaxRetries {
			break
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryConfig.MaxRetries+1, lastErr)
}

// CanRequest checks if the circuit breaker allows requests
func (cb *CircuitBreaker) CanRequest() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailTime) > cb.config.RecoveryTimeout {
			cb.state = CircuitHalfOpen
			cb.requests = 0
			return true
		}
		return false
	case CircuitHalfOpen:
		return cb.requests < cb.config.MaxRequests
	default:
		return false
	}
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures = 0
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitOpen
	} else if cb.failures >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
	}
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// SafeReadBody safely reads and closes an HTTP response body
func SafeReadBody(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()

	// Limit body size to prevent memory exhaustion
	const maxBodySize = 10 * 1024 * 1024 // 10MB
	limitedReader := io.LimitReader(resp.Body, maxBodySize)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// IsNetworkError checks if an error is network-related
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common network errors
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}

	// Check for DNS errors
	if dnsErr, ok := err.(*net.DNSError); ok {
		return dnsErr.Temporary()
	}

	return false
}

// IsRetryableError checks if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Network errors are generally retryable
	if IsNetworkError(err) {
		return true
	}

	// Context timeout/cancellation should not be retried
	if err == context.DeadlineExceeded || err == context.Canceled {
		return false
	}

	return false
}
