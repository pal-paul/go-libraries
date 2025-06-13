package slack

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &ErrRateLimit{time.Duration(retry) * time.Second}
	}
	// Slack seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (m *slack) postRequest(
	endpoint string,
	headers map[string]string,
	reqBody []byte,
) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(
		m.cfg.Context,
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return resp, fmt.Errorf("error post to slack: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.cfg.Token))
	resp, err = m.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	if err := checkStatusCode(resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (m *slack) postForm(
	endpoint string,
	headers map[string]string,
	values url.Values,
) (resp *http.Response, err error) {
	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequestWithContext(m.cfg.Context, http.MethodPost, endpoint, reqBody)
	if err != nil {
		return resp, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.cfg.Token))
	resp, err = m.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	if err := checkStatusCode(resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (m *slack) getRequest(
	endpoint string,
	headers map[string]string,
	values url.Values,
) (resp *http.Response, err error) {
	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequestWithContext(m.cfg.Context, http.MethodGet, endpoint, reqBody)
	if err != nil {
		return resp, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.cfg.Token))
	resp, err = m.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	if err := checkStatusCode(resp); err != nil {
		return resp, err
	}
	return resp, nil
}
