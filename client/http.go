package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AuthedRequest sends an authorized JSON request to the Charm and Glow HTTP servers.
func (cc *Client) AuthedJSONRequest(method string, path string, reqBody interface{}, respBody interface{}) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		return err
	}
	resp, err := cc.AuthedRequest(method, path, "application/json", buf)
	if err != nil {
		return err
	}
	if respBody != nil {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		return dec.Decode(respBody)
	}
	return nil
}

// AuthedRequest sends an authorized request to the Charm and Glow HTTP servers.
func (cc *Client) AuthedRequest(method string, path string, contentType string, reqBody io.Reader) (*http.Response, error) {
	client := &http.Client{}
	cfg := cc.Config
	auth, err := cc.Auth()
	if err != nil {
		return nil, err
	}
	jwt := auth.JWT
	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s:%d%s", cc.httpScheme, cfg.Host, cfg.HTTPPort, path), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", jwt))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("server error: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return resp, nil
}

// AuthedRawRequest sends an authorized request with no request body to the Charm and Glow HTTP servers.
func (cc *Client) AuthedRawRequest(method string, path string) (*http.Response, error) {
	return cc.AuthedRequest(method, path, "", nil)
}