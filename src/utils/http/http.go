package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdhttp "net/http"
)

func Get(ctx context.Context, url string, headers map[string]string, resp interface{}) (error, []byte) {
	return doRequest(ctx, stdhttp.MethodGet, url, headers, nil, resp)
}

func Post(ctx context.Context, url string, headers map[string]string, body interface{}, resp interface{}) error {
	err, _ := doRequest(ctx, stdhttp.MethodPost, url, headers, body, resp)
	return err
}

func doRequest(ctx context.Context, method string, url string, headers map[string]string, body interface{}, resp interface{}) (error, []byte) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err, nil
		}
		reader = bytes.NewReader(payload)
	}

	//创建请求
	req, err := stdhttp.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err, nil
	}

	//遍历映射请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &stdhttp.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err, nil
	}

	if httpResp.StatusCode < stdhttp.StatusOK || httpResp.StatusCode >= stdhttp.StatusMultipleChoices {
		return fmt.Errorf("http %s %s returned status %d: %s", method, url, httpResp.StatusCode, string(respBody)), respBody
	}

	if resp == nil || len(respBody) == 0 {
		return nil, respBody
	}

	if err := json.Unmarshal(respBody, resp); err != nil {
		return err, respBody
	}

	return nil, respBody
}
