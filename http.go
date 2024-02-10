package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type HTTPCheckRequest struct {
	Website           string
	Port              int
	Method            string
	TimeoutResilience uint8
}

type HTTPChecKResponse struct {
	StatusCode   int
	StatusTitle  string
	IP           string
	ResponseTime string
	TLS          bool
}

func newHTTPClient(timeout uint8) *http.Client {
	httpTransport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: time.Duration(timeout) * time.Second,
	}

	return &http.Client{
		Transport: httpTransport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
}

func HTTPChecker(request HTTPCheckRequest) (*HTTPChecKResponse, error) {
	scheme, hostName, err := checkHTTPUrl(request.Website)
	if err != nil {
		return nil, err
	}

	ipAddr, err := net.ResolveIPAddr("ip", hostName)
	if err != nil {
		return &HTTPChecKResponse{
			StatusCode:   113,
			StatusTitle:  "No Route To Host",
			IP:           "",
			TLS:          false,
			ResponseTime: "0s",
		}, nil
	}

	err = checkPort(request.Port)
	if err != nil {
		request.Port = 80
	}

	switch scheme {
	case "http":
		request.Port = 80
	case "https":
		request.Port = 443
	default:
	}

	request.Website = scheme + "://" + hostName

	if request.TimeoutResilience < 1 || request.TimeoutResilience > 60 {
		request.TimeoutResilience = 20
	}

	client := newHTTPClient(request.TimeoutResilience)

	req, err := http.NewRequest(request.Method, request.Website+":"+fmt.Sprint(request.Port), nil)
	if err != nil {
		return nil, errors.Join(errors.New("Making new http request error: "), err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return &HTTPChecKResponse{
				StatusCode:   408,
				StatusTitle:  "Connection timed out",
				IP:           "",
				TLS:          false,
				ResponseTime: fmt.Sprintf("%ds", request.TimeoutResilience),
			}, nil
		} else {
			return nil, errors.Join(errors.New("Error while doing http request: "), err)
		}
	}

	defer resp.Body.Close()

	elapsed := time.Since(start)

	statusTitle := strings.Trim(resp.Status, fmt.Sprintf("%d ", resp.StatusCode))

	return &HTTPChecKResponse{
		StatusCode:   resp.StatusCode,
		StatusTitle:  statusTitle,
		IP:           ipAddr.IP.String(),
		TLS:          resp.TLS.HandshakeComplete,
		ResponseTime: fmt.Sprint(elapsed),
	}, nil
}
