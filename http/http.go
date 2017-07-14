package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func FetchImage(url *url.URL, filename string) error {
	body, err := GetBody(url)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetBody(url *url.URL) ([]byte, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
