/*
   Copyright 2019, orcaer@yeah.net  All rights reserved.
   Author: orcaer@yeah.net
   Last modified: 2019-5-15
   Description: https://github.com/wlgq2/kit_examples
*/

package main

import (
	"context"

	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeClient(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/string"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeResponse,
	).Endpoint()
}
