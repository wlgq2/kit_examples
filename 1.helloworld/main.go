/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-5-15
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type HelloService struct{}

func (HelloService) hello(str string) (string, error) {
	return "hello world!", nil
}

type StringStruct struct {
	Str string `json:"str"`
}

func makeEndpoint(svc HelloService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(StringStruct)
		v, err := svc.hello(req.Str)
		return StringStruct{v}, err
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request StringStruct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {

	svc := HelloService{}

	http.Handle("/string", httptransport.NewServer(
		makeEndpoint(svc),
		decodeRequest,
		encodeResponse))

	http.ListenAndServe(":8080", nil)
}
