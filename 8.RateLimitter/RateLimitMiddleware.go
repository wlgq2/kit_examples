/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-8-30
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ERROR_REQ_LIMIT = errors.New("req limited.")

func MakeLimitter(cnt int, ms int) endpoint.Middleware {
	limitter := rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(ms)), cnt)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if limitter.Allow() {
				return next(ctx, request)
			} else {
				return nil, ERROR_REQ_LIMIT
			}
		}
	}
}
