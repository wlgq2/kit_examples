package main

import (
	"context"
	"encoding/json"
	"errors"
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
		return service.Get(req.Key)
	}
}

func decodeStruct(_ context.Context, req *http.Request) (interface{}, error) {
	var request KeyValueStruct
	err := json.NewDecoder(req.Body).Decode(&request)
	return request, err
}

func encodeStruct(_ context.Context, writer http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(writer).Encode(resp)
}
