package http_client

//go:generate mockgen -source=interface.go -destination=mocks/mock-http-client.go -package=mocks
import (
	"bytes"
	"io"
	"net/http"
)

// New instance of httpClient
func New() IHttpClient {
	return &httpClient{
		client: &http.Client{},
	}
}

// Get a http request to url with headers
//
// Parameters:
//   - url: string
//   - headers: map[string]string
//
// Returns:
//   - []byte: response body
//   - int: response status code
//   - error: error
func (hc *httpClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	if url == "" {
		return nil, 0, errInvalidUrl
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}

// Post a http request to url with headers
//
// Parameters:
//   - url: string
//   - postBody: []byte
//   - headers: map[string]string
//
// Returns:
//   - []byte: response body
//   - int: response status code
//   - error: error
func (hc *httpClient) Post(
	url string,
	postBody []byte,
	headers map[string]string,
) ([]byte, int, error) {
	if url == "" {
		return nil, 0, errInvalidUrl
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, 0, err
	}
	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}

func (hc *httpClient) Put(
	url string,
	postBody []byte,
	headers map[string]string,
) ([]byte, int, error) {
	if url == "" {
		return nil, 0, errInvalidUrl
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, 0, err
	}
	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}

func (hc *httpClient) Delete(
	url string,
	postBody []byte,
	headers map[string]string,
) ([]byte, int, error) {
	if url == "" {
		return nil, 0, errInvalidUrl
	}
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, 0, err
	}
	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}
