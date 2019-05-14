package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
)

func RunService() {
	var svc KeyValueService
	svc = NewMapService()

	setHandler := httptransport.NewServer(
		makeSetEndpoint(svc),
		decodeStruct,
		encodeStruct,
	)

	getHandler := httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeStruct,
		encodeStruct,
	)

	http.Handle("/set", setHandler)
	http.Handle("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}

func ReqTest(key string, value string, url string) (string, error) {

	contentType := "application/json;charset=utf-8"
	req := KeyValueStruct{
		Key:   key,
		Value: value,
	}

	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return key, err
	}

	body := bytes.NewBuffer(data)

	resp, err := http.Post(url, contentType, body)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return key, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return key, err
	}

	return string(content), nil
}

func ReqSet(key string, value string) error {

	_, err := ReqTest(key, value, "http://127.0.0.1:8080/set")
	return err
}

func ReqGet(key string) (string, error) {
	return ReqTest(key, "", "http://127.0.0.1:8080/get")
}

func RunReq() {
	time.Sleep(time.Duration(1) * time.Second)

	fmt.Println("set key1:123")
	ReqSet("key1", "123")

	fmt.Println("set key2:456")
	ReqSet("key2", "456")

	value, _ := ReqGet("key1")
	fmt.Println("get key1:", value)

	value, _ = ReqGet("key2")
	fmt.Println("get key2:", value)
}

func main() {
	go RunReq()

	RunService()

}
