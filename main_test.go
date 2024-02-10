package main

import (
	"fmt"
	"testing"
)

func TestHTTPChecker(t *testing.T) {
	websites := map[string]int{
		"https://google.com":           200,
		"https://facebook.com":         408,
		"https://lfe.org":              200,
		"https://random.websitedomain": 113,
	}

	for website, expectedResCode := range websites {
		req := HTTPCheckRequest{
			Website:           website,
			Port:              80,
			Method:            "GET",
			TimeoutResilience: 10,
		}

		res, err := HTTPChecker(req)

		if err != nil {
			t.Errorf("error for %s: %v", website, err)
			continue
		}

		if res.StatusCode == expectedResCode {
			fmt.Println("test passed")

		} else {
			t.Errorf("test error for %s: expected %d but got %d", website, expectedResCode, res.StatusCode)
		}
	}
}
