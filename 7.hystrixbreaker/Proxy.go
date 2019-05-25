/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-5-25
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

type KeyValueMiddleware func(KeyValueService) KeyValueService

type KeyValueProxy struct {
	ctx  context.Context
	next KeyValueService
	set  endpoint.Endpoint
	get  endpoint.Endpoint
}

func NewProxy(ctx context.Context, urls []string) KeyValueMiddleware {

	var maxAttempts = 3
	var maxTime = 250 * time.Millisecond

	var setEndpointer sd.FixedEndpointer
	var getEndpointer sd.FixedEndpointer
	qps := 100

	for _, url := range urls {
		setEndpoint := makeSetProxy(ctx, url)
		setEndpoint = newHystrixBreaker(setEndpoint, qps)
		setEndpointer = append(setEndpointer, setEndpoint)
		getEndpoint := makeGetProxy(ctx, url)
		getEndpointer = append(getEndpointer, getEndpoint)

	}

	balancerSet := lb.NewRoundRobin(setEndpointer)
	balancerGet := lb.NewRoundRobin(getEndpointer)
	retrySet := lb.Retry(maxAttempts, maxTime, balancerSet)
	retryGet := lb.Retry(maxAttempts, maxTime, balancerGet)

	return func(next KeyValueService) KeyValueService {
		return &(KeyValueProxy{ctx, next, retrySet, retryGet})
	}
}

func (service *KeyValueProxy) Set(key string, value string) (string, error) {
	response, err := service.set(service.ctx, KeyValueStruct{Key: key, Value: value})
	if err != nil {
		return "", err
	}

	resp := response.(KeyValueStruct)

	return resp.Key, nil
}

func (service *KeyValueProxy) Get(key string) (string, error) {
	response, err := service.get(service.ctx, KeyValueStruct{Key: key, Value: ""})
	if err != nil {
		return "", err
	}

	resp := response.(KeyValueStruct)

	return resp.Value, nil

}

func makeSetProxy(ctx context.Context, addr string) endpoint.Endpoint {
	u, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	u.Path = "/set"

	return httptransport.NewClient(
		"GET",
		u,
		encodeReq,
		decodeResp,
	).Endpoint()
}

func makeGetProxy(ctx context.Context, addr string) endpoint.Endpoint {
	u, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	u.Path = "/get"

	return httptransport.NewClient(
		"GET",
		u,
		encodeReq,
		decodeResp,
	).Endpoint()
}
