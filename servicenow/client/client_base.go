package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DefaultRequestTimeout is the default HTTP timeout for a single ServiceNow request.
const DefaultRequestTimeout = 60 * time.Second

// maxRequestAttempts is the maximum number of HTTP attempts (1 original + 3 retries).
const maxRequestAttempts = 4

// Client is the client used to interact with ServiceNow API.
type Client struct {
	BaseURL    string
	Auth       string
	HTTPClient *http.Client
	UserAgent  string
}

// ServiceNowClient defines possible methods to call on the ServiceNowClient.
type ServiceNowClient interface {
	GetObject(ctx context.Context, endpoint string, id string, responseObjectOut Record) error
	GetObjectByName(ctx context.Context, endpoint string, name string, responseObjectOut Record) error
	GetObjectByTitle(ctx context.Context, endpoint string, title string, responseObjectOut Record) error
	GetObjectByQuery(ctx context.Context, endpoint string, query string, responseObjectOut Record) error
	CreateObject(ctx context.Context, endpoint string, record Record) error
	UpdateObject(ctx context.Context, endpoint string, record Record) error
	DeleteObject(ctx context.Context, endpoint string, id string) error

	// Generic untyped record operations used by servicenow_record. The endpoint
	// is a bare table name (e.g. "incident") — these methods append ".do" and
	// JSONv2 query parameters internally.
	CreateRecord(ctx context.Context, table, scope string, fields map[string]string) (map[string]string, error)
	GetRecord(ctx context.Context, table, sysID string) (map[string]string, error)
	GetRecordByQuery(ctx context.Context, table, query string) (map[string]string, error)
	UpdateRecord(ctx context.Context, table, sysID string, fields map[string]string) (map[string]string, error)
}

// BaseResult is representing the default properties of all results.
type BaseResult struct {
	ID               string       `json:"sys_id,omitempty"`
	ProtectionPolicy string       `json:"sys_policy,omitempty"`
	Scope            string       `json:"sys_scope,omitempty"`
	Status           string       `json:"__status,omitempty"`
	Error            *ErrorDetail `json:"__error,omitempty"`
}

// Record is the interface for any BaseResult.
type Record interface {
	GetID() string
	GetScope() string
	GetStatus() string
	GetError() *ErrorDetail
}

// BaseResultList represents the response from the API. Records are always returned inside an array.
type BaseResultList struct {
	Records []json.RawMessage
}

// ErrorDetail is the details of an error. Should be included in the json if status is not success.
type ErrorDetail struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// NotFoundError is returned when a record cannot be located in ServiceNow (zero records or HTTP 404).
type NotFoundError struct {
	Reason string
}

func (e *NotFoundError) Error() string {
	if e.Reason == "" {
		return "record not found"
	}
	return e.Reason
}

// IsNotFound reports whether err (or anything it wraps) is a NotFoundError.
func IsNotFound(err error) bool {
	var nf *NotFoundError
	return errors.As(err, &nf)
}

// NewClient is a factory method used to return a new ServiceNowClient.
func NewClient(baseURL string, username string, password string) *Client {
	credentials := username + ":" + password
	return &Client{
		BaseURL:    baseURL,
		Auth:       "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials)),
		HTTPClient: &http.Client{Timeout: DefaultRequestTimeout},
		UserAgent:  "terraform-provider-servicenow",
	}
}

// GetID returns the ID of a BaseRecord.
func (record BaseResult) GetID() string {
	return record.ID
}

// GetStatus returns the Status of a BaseRecord.
func (record BaseResult) GetStatus() string {
	return record.Status
}

// GetError returns the Error of a BaseRecord, if any.
func (record BaseResult) GetError() *ErrorDetail {
	return record.Error
}

// GetScope returns the Scope of a BaseRecord.
func (record BaseResult) GetScope() string {
	return record.Scope
}

// validateOnlyOneResultReceived returns a NotFoundError when the result list is empty.
func validateOnlyOneResultReceived(results BaseResultList) error {
	if len(results.Records) == 0 {
		return &NotFoundError{Reason: "no records found"}
	}
	if len(results.Records) > 1 {
		return fmt.Errorf("more than one record received")
	}
	return nil
}

// GetObject retrieves an object via a specific endpoint with a GET method and a specified sys_id.
func (client *Client) GetObject(ctx context.Context, endpoint string, id string, responseObjectOut Record) error {
	jsonResponse, err := client.requestJSON(ctx, "GET", endpoint+"?JSONv2&sysparm_query=sys_id="+id, nil)
	if err != nil {
		return err
	}
	return parseResponseToRecord(jsonResponse, responseObjectOut)
}

// GetObjectByName retrieves an object via its name attribute.
func (client *Client) GetObjectByName(ctx context.Context, endpoint string, name string, responseObjectOut Record) error {
	jsonResponse, err := client.requestJSON(ctx, "GET", endpoint+"?JSONv2&sysparm_query=name="+url.QueryEscape(name), nil)
	if err != nil {
		return err
	}
	return parseResponseToRecord(jsonResponse, responseObjectOut)
}

// GetObjectByTitle retrieves an object via its title attribute.
func (client *Client) GetObjectByTitle(ctx context.Context, endpoint string, title string, responseObjectOut Record) error {
	jsonResponse, err := client.requestJSON(ctx, "GET", endpoint+"?JSONv2&sysparm_query=title="+url.QueryEscape(title), nil)
	if err != nil {
		return err
	}
	return parseResponseToRecord(jsonResponse, responseObjectOut)
}

// GetObjectByQuery retrieves an object via a raw encoded query string (e.g. "name=sys_user^operation=read").
func (client *Client) GetObjectByQuery(ctx context.Context, endpoint string, query string, responseObjectOut Record) error {
	jsonResponse, err := client.requestJSON(ctx, "GET", endpoint+"?JSONv2&sysparm_query="+url.QueryEscape(query), nil)
	if err != nil {
		return err
	}
	return parseResponseToRecord(jsonResponse, responseObjectOut)
}

// CreateObject creates a new object in ServiceNow.
func (client *Client) CreateObject(ctx context.Context, endpoint string, objectToCreate Record) error {
	u := endpoint + "?JSONv2&sysparm_action=insert"
	if objectToCreate.GetScope() != "" {
		u += "&sysparm_record_scope=" + objectToCreate.GetScope()
	}

	jsonResponse, err := client.requestJSON(ctx, "POST", u, objectToCreate)
	if err != nil {
		return err
	}

	return parseResponseToRecord(jsonResponse, objectToCreate)
}

// UpdateObject updates an object using a specific endpoint, sys_id and object data.
func (client *Client) UpdateObject(ctx context.Context, endpoint string, object Record) error {
	_, err := client.requestJSON(ctx, "POST", endpoint+"?JSONv2&sysparm_action=update&sysparm_query=sys_id="+object.GetID(), object)
	return err
}

// DeleteObject deletes an object using a specific endpoint and sys_id.
func (client *Client) DeleteObject(ctx context.Context, endpoint string, id string) error {
	_, err := client.requestJSON(ctx, "POST", endpoint+"?JSONv2&sysparm_action=deleteRecord&sysparm_sys_id="+id, nil)
	return err
}

// requestJSON executes an HTTP request and returns the raw response data.
// It retries up to 3 additional times (4 attempts total) with exponential backoff
// on transient errors: HTTP 429, 502, 503, 504, and network failures. The
// Retry-After header is honored when present on 429 responses. Context
// cancellation during backoff aborts the retry loop immediately.
func (client *Client) requestJSON(ctx context.Context, method string, path string, jsonData interface{}) ([]byte, error) {
	var marshalled []byte
	if jsonData != nil {
		m, err := json.Marshal(jsonData)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		marshalled = m
	}

	httpClient := client.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultRequestTimeout}
	}

	tflog.Debug(ctx, "ServiceNow API request", map[string]interface{}{
		"method": method,
		"path":   path,
	})

	var lastErr error
	var retryAfter time.Duration
	for attempt := 0; attempt < maxRequestAttempts; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			if retryAfter > backoff {
				backoff = retryAfter
			}

			tflog.Debug(ctx, "ServiceNow API retry", map[string]interface{}{
				"method":  method,
				"path":    path,
				"attempt": attempt,
				"backoff": backoff.String(),
			})

			timer := time.NewTimer(backoff)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			}
		}

		responseData, statusCode, ra, err := client.doSingleRequest(ctx, httpClient, method, path, marshalled)
		if err != nil {
			// NotFoundError: do not retry, return immediately.
			if IsNotFound(err) {
				return nil, err
			}
			// Context errors: do not retry, return immediately.
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			// Distinguish retryable vs non-retryable HTTP responses by status code.
			if statusCode != 0 && !isRetryableStatus(statusCode) {
				return nil, err
			}
			// Retryable: transport error (statusCode == 0) or retryable status code.
			lastErr = err
			retryAfter = ra
			continue
		}
		return responseData, nil
	}
	return nil, lastErr
}

// doSingleRequest performs a single HTTP request attempt. It returns the response
// body on success, otherwise returns an error along with the response status code
// (0 for transport errors) and any Retry-After duration parsed from the response.
func (client *Client) doSingleRequest(ctx context.Context, httpClient *http.Client, method, path string, marshalled []byte) ([]byte, int, time.Duration, error) {
	body := bytes.NewBuffer(marshalled)
	request, err := http.NewRequestWithContext(ctx, method, client.BaseURL+"/"+path, body)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("build request: %w", err)
	}

	request.Header.Set("Authorization", client.Auth)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	if client.UserAgent != "" {
		request.Header.Set("User-Agent", client.UserAgent)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		tflog.Debug(ctx, "ServiceNow API request failed", map[string]interface{}{
			"method": method,
			"path":   path,
			"error":  err.Error(),
		})
		// Honor context errors (cancellation / deadline) verbatim so callers can detect them.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, 0, 0, ctxErr
		}
		return nil, 0, 0, err
	}
	defer func() { _ = response.Body.Close() }()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		tflog.Debug(ctx, "ServiceNow API response body read failed", map[string]interface{}{
			"method": method,
			"path":   path,
			"error":  err.Error(),
		})
		// Treat body read failure as a transport-style error (no status code).
		return nil, 0, 0, fmt.Errorf("read response body: %w", err)
	}

	if response.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, "ServiceNow API response not found", map[string]interface{}{
			"method": method,
			"path":   path,
			"status": response.StatusCode,
		})
		return nil, response.StatusCode, 0, &NotFoundError{Reason: fmt.Sprintf("HTTP 404 from %s %s", method, path)}
	}
	if response.StatusCode >= 300 || response.StatusCode < 200 {
		tflog.Debug(ctx, "ServiceNow API response error status", map[string]interface{}{
			"method": method,
			"path":   path,
			"status": response.StatusCode,
		})
		ra := parseRetryAfter(response.Header.Get("Retry-After"))
		return nil, response.StatusCode, ra, fmt.Errorf("HTTP response status %s, %s", response.Status, responseData)
	}

	tflog.Trace(ctx, "ServiceNow API response", map[string]interface{}{
		"status": response.StatusCode,
	})

	return responseData, response.StatusCode, 0, nil
}

// isRetryableStatus reports whether an HTTP status code is one we should retry.
func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

// parseRetryAfter parses a Retry-After header value as either an integer number
// of seconds or an HTTP-date. Returns 0 when the header is missing or malformed.
func parseRetryAfter(h string) time.Duration {
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func parseResponseToRecord(jsonResponse []byte, responseObjectOut Record) error {
	baseResultsList := BaseResultList{}
	if err := json.Unmarshal(jsonResponse, &baseResultsList); err != nil {
		return err
	}

	if err := validateOnlyOneResultReceived(baseResultsList); err != nil {
		return err
	}

	if err := json.Unmarshal(baseResultsList.Records[0], responseObjectOut); err != nil {
		return err
	}

	return checkStatus(responseObjectOut)
}

// checkStatus inspects the ServiceNow JSONv2 envelope status field.
func checkStatus(record Record) error {
	if record.GetStatus() != "success" {
		errDetail := record.GetError()
		if errDetail == nil {
			return fmt.Errorf("error from ServiceNow: unknown failure")
		}
		return fmt.Errorf("error from ServiceNow -> %s: %s", errDetail.Message, errDetail.Reason)
	}
	return nil
}
