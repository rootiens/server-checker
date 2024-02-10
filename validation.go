package main

import (
	"errors"
	"net/url"
)

func checkHTTPUrl(u string) (string, string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", "", errors.Join(errors.New("URL Parse Error: "), err)
	}

	if parsedURL.Scheme == "" {
		u = "https://" + u
	}

	parsedURL, _ = url.Parse(u)

	_, err = url.ParseRequestURI(u)
	if err != nil {
		return "", "", errors.Join(errors.New("URL Parse Request Error: "), err)
	}

	return parsedURL.Scheme, parsedURL.Host, nil
}

func checkPort(port int) error {
	if port < 0 || port > 65353 {
		return errors.New("Invalid PORT")
	}

	return nil
}
