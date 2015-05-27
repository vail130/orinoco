package httputils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func GetDataFromUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	return data, err
}

func PostDataToUrl(url string, contentType string, data []byte) ([]byte, error) {
	response, err := http.Post(url, contentType, bytes.NewBuffer(data))
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	return responseData, err
}

func Delete(url string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return make([]byte, 0), err
	}
	response, err := http.DefaultClient.Do(req)
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	return responseData, err
}
