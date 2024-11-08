package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func SendHTTPRequest(
	httpMethod string,
	timeout time.Duration,
	headers map[string]string,
	url string,
) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("SendHTTPRequest -> NewRequestWithContext")
		return nil, http.StatusInternalServerError, errors.New("Error creating request")
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("SendHTTPRequest -> client.Do")
		if err == context.DeadlineExceeded {
			return nil, http.StatusRequestTimeout, errors.New("Request timed out")
		} else {
			return nil, http.StatusInternalServerError, errors.New("Error sending request")
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("SendHTTPRequest -> io.ReadAll")
		return nil, http.StatusInternalServerError, errors.New("Error reading response body")
	}

	return body, http.StatusOK, nil
}

func GetJSON(body []byte, v interface{}) error {
	err := json.Unmarshal(body, v)
	return err
}
