package git

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseUrl = "https://api.github.com"
	accept  = "application/vnd.github+json"
)

func (g *git) get(basePath string, path string, qs url.Values) (*http.Response, error) {
	uStr := g.cfg.BaseURL
	if uStr == "" {
		uStr = baseUrl
	}
	if len(basePath) > 0 {
		uStr = fmt.Sprintf("%s/%s", uStr, basePath)
	}
	if len(path) > 0 {
		uStr = fmt.Sprintf("%s/%s", uStr, path)
	}
	u, err := url.Parse(uStr)
	if err != nil {
		err = fmt.Errorf("failed to url parse %s: %v", uStr, err)
		return nil, err
	}
	if qs != nil {
		u.RawQuery = qs.Encode()
	}
	client := &http.Client{}
	uStr = u.String()

	// fmt.Println(uStr)
	req, err := http.NewRequest(http.MethodGet, uStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *git) post(basePath string, path string, qs url.Values, reqBody []byte) (*http.Response, error) {
	uStr := g.cfg.BaseURL
	if uStr == "" {
		uStr = baseUrl
	}
	if len(basePath) > 0 {
		uStr = uStr + "/" + basePath
	}
	if len(path) > 0 {
		uStr = uStr + "/" + path
	}
	u, err := url.Parse(uStr)
	if err != nil {
		err = fmt.Errorf("failed to url parse %s: %v", uStr, err)
		return nil, err
	}
	if qs != nil {
		u.RawQuery = qs.Encode()
	}
	client := &http.Client{}
	uStr = u.String()
	req, err := http.NewRequest(http.MethodPost, uStr, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *git) put(basePath string, path string, qs url.Values, reqBody []byte) (*http.Response, error) {
	uStr := g.cfg.BaseURL
	if uStr == "" {
		uStr = baseUrl
	}
	if len(basePath) > 0 {
		uStr = uStr + "/" + basePath
	}
	if len(path) > 0 {
		uStr = uStr + "/" + path
	}
	u, err := url.Parse(uStr)
	if err != nil {
		err = fmt.Errorf("failed to url parse %s: %v", uStr, err)
		return nil, err
	}
	if qs != nil {
		u.RawQuery = qs.Encode()
	}
	client := &http.Client{}
	uStr = u.String()
	// fmt.Println(uStr)
	req, err := http.NewRequest(http.MethodPut, uStr, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *git) patch(basePath string, path string, qs url.Values, reqBody []byte) (*http.Response, error) {
	uStr := g.cfg.BaseURL
	if uStr == "" {
		uStr = baseUrl
	}
	if len(basePath) > 0 {
		uStr = uStr + "/" + basePath
	}
	if len(path) > 0 {
		uStr = uStr + "/" + path
	}
	u, err := url.Parse(uStr)
	if err != nil {
		err = fmt.Errorf("failed to url parse %s: %v", uStr, err)
		return nil, err
	}
	if qs != nil {
		u.RawQuery = qs.Encode()
	}
	client := &http.Client{}
	uStr = u.String()
	req, err := http.NewRequest(http.MethodPatch, uStr, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
