package adapter

import (
	"bytes"
	"net/http"
)

type DestinationAdapter struct {
	destinationHost string
}

func NewDestinationAdapter(destinationHost string) *DestinationAdapter {
	return &DestinationAdapter{
		destinationHost: destinationHost,
	}
}

func (d *DestinationAdapter) Do(path string, method string, queries map[string]string, headers map[string][]string, body []byte) ([]byte, error) {
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
		return nil, err
	}

	// set request headers
	for key, value := range headers {
		for _, v := range value {
			request.Header.Add(key, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body = buf.Bytes()

	return body, nil
}
