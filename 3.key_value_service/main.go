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

func RunService() {
	var svc KeyValueService
	svc = NewMapService()

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

	http.Handle("/set", setHandler)
	http.Handle("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}

func RunClient() {
	time.Sleep(time.Duration(1) * time.Second)

	set := MakeSetClient(context.Background(), ":8080")
	get := MakeGetClient(context.Background(), ":8080")

	fmt.Println("set key1:abc")
	set(context.Background(), KeyValueStruct{Key: "key1", Value: "abc"})

	fmt.Println("set key2:123")
	set(context.Background(), KeyValueStruct{Key: "key2", Value: "123"})

	response, err := get(context.Background(), KeyValueStruct{Key: "key1", Value: "0"})
	if err != nil {
		return
	}
	resp := response.(KeyValueStruct)
	fmt.Println("get : ", resp.Key, resp.Value)

	response, err = get(context.Background(), KeyValueStruct{Key: "key2", Value: "0"})
	if err != nil {
		return
	}
	resp = response.(KeyValueStruct)
	fmt.Println("get : ", resp.Key, resp.Value)
}

func main() {
	go RunClient()

	RunService()

}
