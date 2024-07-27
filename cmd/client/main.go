// Copyright 2024 The Echo gRPC Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	echov1 "github.com/automenu/echo-grpc/api/echo/v1"
	"github.com/automenu/echo-grpc/api/echo/v1/echov1connect"
)

func main() {
	// Parse server address from command line flag
	serverAddr := flag.String(
		"addr",
		"http://localhost:8080",
		"The server address in the format of scheme://host:port",
	)
	flag.Parse()

	// Initialize Echo API client
	client := echov1connect.NewEchoAPIClient(
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				MaxConnsPerHost:       100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 10 * time.Second,
		},
		*serverAddr,
		connect.WithGRPC(),
	)

	// Call Echo API
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.Echo(
		ctx,
		connect.NewRequest(&echov1.EchoRequest{
			Message: "Hello gRPC!",
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", res.Msg.Reply)
}
