package main

import "fmt"

func main() {
	req := HTTPCheckRequest{
		Website:           "https://crstrip.com",
		Port:              443,
		Method:            "GET",
		TimeoutResilience: 10,
	}

	res, err := HTTPChecker(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Host: ", req.Website)
	fmt.Println("Code: ", res.StatusCode)
	fmt.Println("Title: ", res.StatusTitle)
	fmt.Println("IP: ", res.IP)
	fmt.Println("TLS: ", res.TLS)
	fmt.Println("Response Time: ", res.ResponseTime)
}
