package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)


type jwt struct {
	Token string
}

func getToken() (string, error) {
	link := "http://localhost:8056/login"
	//link := "https://tmp.mama.sh/api/login"
	var jsonStr = []byte(`{"username":"martin", "password":"T3stpass!"}`)
	req, err := http.NewRequest("POST", link, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var token jwt
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
