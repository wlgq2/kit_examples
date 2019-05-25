/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-5-15
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
)

func RunService(port string) {
	var svc KeyValueService
	svc = NewMapService()
	svc = NewLogMiddleware(svc)

	setHandler := httptransport.NewServer(
		makeSetEndpoint(svc),
		decodeReq,
		encodeResp,
	)

	getHandler := httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeReq,
		encodeResp,
	)

	mux := http.NewServeMux()
	mux.Handle("/set", setHandler)
	mux.Handle("/get", getHandler)
	http.ListenAndServe(port, mux)
}

func RunClient(port string) {
	time.Sleep(time.Duration(1) * time.Second)

	set := MakeSetClient(context.Background(), port)
	get := MakeGetClient(context.Background(), port)

	fmt.Println("set key1:abc")
	set(context.Background(), KeyValueStruct{Key: "key1", Value: "abc"})

	fmt.Println("set key2:123")
	set(context.Background(), KeyValueStruct{Key: "key2", Value: "123"})

	response, err := get(context.Background(), KeyValueStruct{Key: "key2", Value: "0"})
	if err != nil {
		return
	}
	resp := response.(KeyValueStruct)
	fmt.Println("get : ", resp.Key, resp.Value)

	response, err = get(context.Background(), KeyValueStruct{Key: "key1", Value: "0"})
	if err != nil {
		return
	}
	resp = response.(KeyValueStruct)
	fmt.Println("get : ", resp.Key, resp.Value)
}

func RunProxy(port string, ports []string) {
	var service KeyValueService
	service = &(MapService{})
	var proxyUrls []string
	for _, port := range ports {
		url := "http://localhost" + port
		proxyUrls = append(proxyUrls, url)
	}

	service = NewProxy(context.Background(), proxyUrls)(service)

	setHandler := httptransport.NewServer(
		makeSetEndpoint(service),
		decodeReq,
		encodeResp,
	)

	getHandler := httptransport.NewServer(
		makeGetEndpoint(service),
		decodeReq,
		encodeResp,
	)

	mux := http.NewServeMux()
	mux.Handle("/set", setHandler)
	mux.Handle("/get", getHandler)
	http.ListenAndServe(port, mux)
}
func main() {

	go RunClient(":8000")
	ports := []string{
		":8080",
		":8081",
		":8082",
	}
	for _, port := range ports {
		go RunService(port)
	}

	RunProxy(":8000", ports)
}
