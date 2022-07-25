package helper

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"io/ioutil"

// 	//"log"
// 	"net/http"

// 	"github.com/rs/zerolog/log"
// )

// func Http_Get(url string, headers map[string]string) (outresponse []byte, err error) {
// 	httpposturl := url

// 	request, err := http.NewRequest("GET", httpposturl, nil)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}

// 	for k, v := range headers {
// 		fmt.Printf("key[%s] value[%s]\n", k, v)
// 		//request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		request.Header.Set(k, v)
// 	}

// 	res, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}
// 	err = res.Body.Close()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}
// 	fmt.Printf("%s", body)
// 	return body, err
// }

// func Http_Post(httpposturl string, jsonpayload []byte) ([]byte, error) {
// 	log.Info().Msg(httpposturl)
// 	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonpayload))
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}

// 	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	if err != nil {
// 		log.Fatal().Err(err).Msg("")
// 		return nil, err
// 	}
// 	body, _ := ioutil.ReadAll(response.Body)
// 	return body, err
// }
