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

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
)

func newHystrixBreaker(endpoint endpoint.Endpoint, qps int) endpoint.Endpoint {

	commandName := "my-endpoint"
	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                1000 * 30,
		ErrorPercentThreshold:  1,
		SleepWindow:            10000,
		MaxConcurrentRequests:  1000,
		RequestVolumeThreshold: 5,
	})
	middleware := circuitbreaker.Hystrix(commandName)
	rst := middleware(endpoint)
	rst = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(rst)
	return rst
}
