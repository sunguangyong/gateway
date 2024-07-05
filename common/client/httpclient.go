package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var cli *http.Client
var once sync.Once

func init() {
	once.Do(func() {
		cli = new(http.Client)
	})
}

func PostJson(req, body interface{}, host, path string) error {
	bs, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("....req...", string(bs))
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", host, path), bytes.NewBuffer(bs))
	httpReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil
	}
	resp, err := cli.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println("....resp...", string(bodyBytes))
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		return err
	}
	return nil
}

func Get(req map[string]string, body interface{}, host, path string) error {
	pathParams := make([]string, 0)
	for k, v := range req {
		pathParams = append(pathParams, fmt.Sprintf("%s=%s", k, v))
	}

	url := fmt.Sprintf("%s/%s?%s", host, path, strings.Join(pathParams, "&"))
	fmt.Println("....req...", url)
	httpReq, err := http.NewRequest("GET", url, nil)
	httpReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil
	}
	resp, err := cli.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println("....resp...", string(bodyBytes))
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		fmt.Println("err:------------------", err)
		return err
	}
	return nil
}
