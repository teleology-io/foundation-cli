package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/teleology-io/hermes"
)

type Api struct {
	hermes.Client
}

func Create(apiKey string) Api {
	baseUrl := os.Getenv("FOUNDATION_API_URL")
	if baseUrl == "" {
		baseUrl = "https://foundation-api.teleology.io"
	}

	return Api{
		hermes.Create(hermes.ClientConfiguration{
			BaseURL: baseUrl,
			Headers: hermes.Headers{
				"X-Api-Key": apiKey,
			},
			TransformResponse: func(res *hermes.Response, err error) (*hermes.Response, error) {
				if res.StatusCode != http.StatusOK {
					fmt.Println("ERROR:", string(res.Data))
					return nil, errors.New("request failed")
				}

				return res, err
			},
		}),
	}
}

func indent(data interface{}) (string, error) {
	out, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (api Api) GetEnvironment() {
	res, err := api.Send(hermes.Request{
		Url: "/v1/environment",
	})
	if err != nil {
		return
	}

	var env map[string]interface{}
	if err = json.Unmarshal(res.Data, &env); err != nil {
		return
	}

	for k, v := range env {
		fmt.Printf("%s=%v\n", k, v)
	}
}

func (api Api) GetConfiguration() {
	res, err := api.Send(hermes.Request{
		Url: "/v1/configuration",
	})
	if err != nil {
		return
	}

	var data = struct {
		Content  string `json:"content"`
		MimeType string `json:"mime_type"`
	}{}
	if err := json.Unmarshal(res.Data, &data); err != nil {
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

func (api Api) GetVariable(variableName string, uniqueID string) {
	data := map[string]string{
		"name": variableName,
	}
	if uniqueID != "" {
		data["uid"] = uniqueID
	}

	res, err := api.Send(hermes.Request{
		Method: hermes.POST,
		Url:    "/v1/variable",
		Headers: hermes.Headers{
			"Content-Type": "application/json",
		},
		Data: data,
	})
	if err != nil {
		return
	}

	var response = struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}{}
	if err := json.Unmarshal(res.Data, &response); err != nil {
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
