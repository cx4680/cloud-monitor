package tools

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGet(path string) (string, error) {
	resp, err := http.Get(path)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// HttpPost params 格式 a=123&b=234
func HttpPost(path string, params string) (string, error) {
	resp, err := http.Post(path,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HttpPostForm(path string, params map[string][]string) (string, error) {
	resp, err := http.PostForm(path, params)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}

func HttpPostJson(path string, params interface{}, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", path, bytes.NewBuffer([]byte(ToString(params))))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HttpHeaderPostJson(path string, header http.Header,params interface{}) (string, error) {
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer([]byte(ToString(params))))
	if err != nil {
		return "", err
	}
	req.Header=header
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}