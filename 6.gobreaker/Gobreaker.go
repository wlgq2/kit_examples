/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-5-25
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"time"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
)

func newGobreaker(endpoint endpoint.Endpoint, qps int) endpoint.Endpoint {

	middleware := circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "my gobreaker",
		MaxRequests: 100,             //半开状态
		Interval:    time.Second,     //计数周期
		Timeout:     time.Second * 2, //进入半开状态周期
		//ReadyToTrip   func(counts Counts) bool   打开条件，默认5次出错
		//OnStateChange func(name string, from State, to State)   //状态切换触发
	}))
	rst := middleware(endpoint)
	rst = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(rst)
	return rst
}
