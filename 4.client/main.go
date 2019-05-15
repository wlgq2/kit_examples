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
	svc = NewLogMiddleware(svc)

	setHandler := httptransport.NewServer(
		makeSetEndpoint(svc),
		decodeStructReq,
		encodeStructResp,
	)

	getHandler := httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeStructReq,
		encodeStructResp,
	)

	http.Handle("/set", setHandler)
	http.Handle("/get", getHandler)
	http.ListenAndServe(":8081", nil)
}

func RunClient() {
	time.Sleep(time.Duration(1) * time.Second)

	fmt.Println("set key1:abc")
	set := MakeSetClient(context.Background(), ":8081")
	set(context.Background(), KeyValueStruct{Key: "key1", Value: "abcd"})

	get := MakeGetClient(context.Background(), ":8081")
	response, err := get(context.Background(), KeyValueStruct{Key: "key1", Value: "0"})
	if err != nil {
		return
	}

	resp := response.(KeyValueStruct)

	fmt.Println("get : ", resp.Key, resp.Value)
}

func main() {
	go RunClient()

	RunService()

}
