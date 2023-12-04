package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

var (
	errStatusNotOK = errors.New("status not ok")
)

func GetJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func PostJSONWithHeaders(url string, headers map[string]string, payload interface{}, target interface{}) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(p))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return errStatusNotOK
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func PostJSON(url string, payload interface{}, target interface{}) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	r, err := myClient.Post(url, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return errStatusNotOK
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func GetJSONWithStatusCode(url string, target interface{}) (error, int) {
	r, err := myClient.Get(url)
	if r == nil {
		return err, 404
	}
	if err != nil {
		return err, r.StatusCode
	}
	defer r.Body.Close()

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return errStatusNotOK, r.StatusCode
	}

	return json.NewDecoder(r.Body).Decode(target), r.StatusCode
}

func GetJSONWithAgent(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"User-Agent": {req.UserAgent()},
	}

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
