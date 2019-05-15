/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-5-15
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type HelloService struct{}

func (HelloService) hello(str string) (string, error) {
	return str, nil
}

type StringRequest struct {
	Str string `json:"str"`
}

type StringResponse struct {
	Hello string `json:"hello"`
}

func makeEndpoint(svc HelloService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(StringRequest)
		v, err := svc.hello(req.Str)
		return StringResponse{v}, err
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request StringRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func RunService() {
	svc := HelloService{}

	http.Handle("/string", httptransport.NewServer(
		makeEndpoint(svc),
		decodeRequest,
		encodeResponse))

	http.ListenAndServe(":8080", nil)
}

func RunReq() {

	time.Sleep(time.Duration(1) * time.Second)

	url := "http://127.0.0.1:8080/string"
	contentType := "application/json;charset=utf-8"

	req := StringRequest{}
	req.Str = "world!"
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body := bytes.NewBuffer(data)

	resp, err := http.Post(url, contentType, body)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(content))
}

func main() {
	go RunReq()

	RunService()
}
