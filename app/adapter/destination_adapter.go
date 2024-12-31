package adapter

import (
	"bytes"
	"encoding/json"
	"gemini-poc/app/dto"
	"net/http"

	"go.uber.org/zap"
)

type DestinationAdapter struct {
	destinationHost string

	log *zap.Logger
}

func NewDestinationAdapter(
	destinationHost string,
	log *zap.Logger,
) *DestinationAdapter {
	return &DestinationAdapter{
		destinationHost: destinationHost,
		log:             log,
	}
}

func (d *DestinationAdapter) Do(path string, method string, queries map[string]string, headers map[string][]string, body []byte) (*int, []byte, error) {
	client := &http.Client{}
	bodyReader := bytes.NewReader(body)

	// set query params
	if len(queries) > 0 {
		query := ""
		for key, value := range queries {
			query += key + "=" + value + "&"
		}
		path += "?" + query[:len(query)-1]
	}

	request, err := http.NewRequest(method, d.destinationHost+path, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	// set request headers
	for key, value := range headers {
		for _, v := range value {
			request.Header.Add(key, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body = buf.Bytes()

	return &response.StatusCode, body, nil
}

func (d *DestinationAdapter) Login(req *dto.AuthRequest) (*dto.AuthResponse, error) {
	url := d.destinationHost + "/api/v1/auth/authenticate"

	payload, err := json.Marshal(req)
	if err != nil {
		d.log.Error("Error marshalling request", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))

	if err != nil {
		d.log.Error("Error creating request", zap.Error(err))
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")

	rawRes, err := client.Do(httpReq)
	if err != nil {
		d.log.Error("Error sending request", zap.Error(err))
		return nil, err
	}
	defer rawRes.Body.Close()

	var res dto.AuthResponse
	if err := json.NewDecoder(rawRes.Body).Decode(&res); err != nil {
		d.log.Error("Error decoding response", zap.Error(err))
		return nil, err
	}

	return &res, nil
}

func (d *DestinationAdapter) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	url := d.destinationHost + "/api/v1/auth/refresh"

	payload, err := json.Marshal(req)
	if err != nil {
		d.log.Error("Error marshalling request", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))

	if err != nil {
		d.log.Error("Error creating request", zap.Error(err))
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")

	rawRes, err := client.Do(httpReq)
	if err != nil {
		d.log.Error("Error sending request", zap.Error(err))
		return nil, err
	}
	defer rawRes.Body.Close()

	var res dto.AuthResponse
	if err := json.NewDecoder(rawRes.Body).Decode(&res); err != nil {
		d.log.Error("Error decoding response", zap.Error(err))
		return nil, err
	}

	return &res, nil
}
