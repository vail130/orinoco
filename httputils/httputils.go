package httputils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Delete(url string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return make([]byte, 0), err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func GetDataFromUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func PostDataToUrl(url string, contentType string, data []byte) ([]byte, error) {
	response, err := http.Post(url, contentType, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func PutDataToUrl(url string, contentType string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return make([]byte, 0), err
	}
	req.Header.Add("Content-Type", contentType)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
