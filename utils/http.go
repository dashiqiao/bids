package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//通用发请求的

func DoRequest(method string, path string, data []byte, requestId string) (*http.Response, []byte, error) {
	buffer := bytes.Buffer{}
	var req *http.Request
	var err error
	var resp *http.Response
	// var client http.Client
	client := http.Client{Timeout: 5 * time.Second}

	if data != nil {
		buffer.Write(data)
		req, err = http.NewRequest(method, path, &buffer)
		if err != nil {
			return nil, nil, err
		}
	} else {
		req, err = http.NewRequest(method, path, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	//req, err := http.NewRequest(method, path, &buffer)
	if method != "DELETE" {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("X-REQUEST-ID", requestId)

	if method == "GET" {
		resp, err = client.Get(path)
		if err != nil {
			return nil, nil, err
		}
	} else {
		resp, err = client.Do(req)
		if err != nil {
			return nil, nil, err
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		// type ErrBody struct {
		// 	Error       string `json:"error"`
		// 	Description string `json:"error_description"`
		// }
		//
		// errBody := ErrBody{}
		// err := json.Unmarshal(body, &errBody)
		// if err == nil {
		// 	body = []byte(errBody.Description)
		// }
		log.Println(resp.StatusCode)
		return resp, nil, errors.New("请求失败！")
	}

	return resp, body, nil
}
