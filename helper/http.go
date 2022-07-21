package helper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Http_Get(url string, headers map[string]string) (outresponse []byte, err error) {
	httpposturl := url

	request, error := http.NewRequest("GET", httpposturl, nil)
	if error != nil {
		log.Fatal(error)
		return nil, error
	}

	for k, v := range headers {
		fmt.Printf("key[%s] value[%s]\n", k, v)
		//request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		request.Header.Set(k, v)
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, error
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, error
	}
	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil, error
	}
	fmt.Printf("%s", body)
	return body, err
}

func Http_Post(httpposturl string, headers map[string]string, jsonpayload []byte) ([]byte, error) {
	fmt.Println(httpposturl)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonpayload))
	if error != nil {
		log.Fatal(error)
		return nil, error
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		log.Fatal(error)
		return nil, error
	}
	defer response.Body.Close()

	if error != nil {
		log.Fatal(error)
		return nil, error
	}
	body, _ := ioutil.ReadAll(response.Body)
	return body, error
}
