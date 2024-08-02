package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type RequestMethod string

const (
	GET     RequestMethod = "GET"
	HEAD    RequestMethod = "HEAD"
	POST    RequestMethod = "POST"
	PUT     RequestMethod = "PUT"
	DELETE  RequestMethod = "DELETE"
	CONNECT RequestMethod = "CONNECT"
	OPTIONS RequestMethod = "OPTIONS"
	TRACE   RequestMethod = "TRACE"
	PATCH   RequestMethod = "PATCH"
)

type Headers = map[string]string

type Request struct {
	Method  RequestMethod
	Headers *Headers
	Url     string
	Data    interface{}
}

type ApiClient struct {
	baseUrl string
	headers *Headers
	client  http.Client
}

func Create(apiKey string) ApiClient {
	baseUrl := os.Getenv("FOUNDATION_API_URL")
	if baseUrl == "" {
		baseUrl = "https://foundation-api.teleology.io"
	}

	return ApiClient{
		baseUrl: baseUrl,
		headers: &Headers{
			"X-Api-Key": apiKey,
		},
		client: http.Client{},
	}
}

func _newrequest(method string, url string, data interface{}) (*http.Request, error) {
	if data == nil {
		return http.NewRequest(method, url, nil)
	}

	out, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, bytes.NewBuffer(out))
}

func _request(c ApiClient, request Request) ([]byte, error) {
	base, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, err
	}
	ref, err := url.Parse(request.Url)
	if err != nil {
		return nil, err
	}

	url := base.ResolveReference(ref).String()

	req, err := _newrequest(string(request.Method), url, request.Data)
	if err != nil {
		return nil, err
	}

	if c.headers != nil {
		for k, v := range *c.headers {
			req.Header.Set(k, v)
		}
	}

	if request.Headers != nil {
		for k, v := range *request.Headers {
			req.Header.Set(k, v)
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Println("ERROR:", string(body))
		return nil, errors.New("request failed")
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func indent(data interface{}) (string, error) {
	out, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (client ApiClient) GetEnvironment() {
	res, err := _request(client, Request{
		Url:    "/v1/environment",
		Method: GET,
	})
	if err != nil {
		return
	}

	var env map[string]interface{}
	if err = json.Unmarshal(res, &env); err != nil {
		return
	}

	for k, v := range env {
		fmt.Printf("%s=%v\n", k, v)
	}
}

func (client ApiClient) GetConfiguration() {
	res, err := _request(client, Request{
		Url:    "/v1/configuration",
		Method: GET,
	})
	if err != nil {
		return
	}

	var data = struct {
		Content  string `json:"content"`
		MimeType string `json:"mime_type"`
	}{}
	if err := json.Unmarshal(res, &data); err != nil {
		return
	}

	output := data.Content

	if data.MimeType == "application/json" {
		var content map[string]interface{}
		if err = json.Unmarshal([]byte(data.Content), &content); err != nil {
			return
		}

		out, err := indent(&content)
		if err == nil {
			output = out
		}
	}

	fmt.Printf("%s\n", output)
}

func (client ApiClient) GetVariable(variableName string, uniqueID string) {
	data := map[string]string{
		"name": variableName,
	}
	if uniqueID != "" {
		data["uid"] = uniqueID
	}

	res, err := _request(client, Request{
		Url:    "/v1/variable",
		Method: POST,
		Data:   data,
		Headers: &Headers{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return
	}

	var response = struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}{}
	if err := json.Unmarshal(res, &response); err != nil {
		return
	}

	value := response.Value

	switch response.Value.(type) {
	case map[string]interface{}:
		{
			out, err := indent(&response.Value)
			if err == nil {
				value = out
			}
			break
		}
	default:
		{
			break
		}
	}

	fmt.Printf("%v\n", value)
}
