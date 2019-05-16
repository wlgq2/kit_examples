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
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type KeyValueService interface {
	Set(string, string) (string, error)
	Get(string) (string, error)
}

var ERROR_NOT_KEY = errors.New("can find this key.")

type MapService struct {
	datas map[string]string
}

func NewMapService() *MapService {
	rst := MapService{
		datas: make(map[string]string),
	}
	return &rst
}

func (service *MapService) Set(key string, value string) (string, error) {
	service.datas[key] = value
	return key, nil
}

func (service *MapService) Get(key string) (string, error) {
	value, ok := service.datas[key]
	if ok {
		return value, nil
	} else {
		return value, ERROR_NOT_KEY
	}

}

type KeyValueStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func makeSetEndpoint(service KeyValueService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(KeyValueStruct)
		_, err := service.Set(req.Key, req.Value)
		return req, err
	}
}

func makeGetEndpoint(service KeyValueService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(KeyValueStruct)
		value, err := service.Get(req.Key)
		return KeyValueStruct{req.Key, value}, err
	}
}

func decodeReq(_ context.Context, req *http.Request) (interface{}, error) {
	var request KeyValueStruct
	err := json.NewDecoder(req.Body).Decode(&request)
	return request, err
}

func encodeResp(_ context.Context, writer http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(writer).Encode(resp)
}

func encodeReq(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeResp(_ context.Context, r *http.Response) (interface{}, error) {
	var response KeyValueStruct
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
