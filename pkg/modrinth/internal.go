package modrinth

import (
	"io"
	"net/http"
)

const baseUrl string = "https://api.modrinth.com/v2"

func makeRequest(url string) ([]byte, error) {
	resp, getErr := http.Get(url)

	if getErr != nil {
		return []byte{}, getErr
	}

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		return []byte{}, readErr
	}

	return body, nil
}
