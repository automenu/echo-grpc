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

package echo

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	echov1 "github.com/automenu/echo-grpc/api/echo/v1"
	"github.com/automenu/echo-grpc/api/echo/v1/echov1connect"
	"go.uber.org/zap"
)

// echoAPIHandler implements the Echo API handlers.
type echoAPIHandler struct {
	echov1connect.UnimplementedEchoAPIHandler
}

// NewEchoAPIHandler creates a new Echo API handler.
func NewEchoAPIHandler() (string, http.Handler) {
	return echov1connect.NewEchoAPIHandler(
		&echoAPIHandler{},
		connect.WithReadMaxBytes(1<<20),     // 1MB
		connect.WithSendMaxBytes(1<<20),     // 1MB
		connect.WithCompressMinBytes(1<<10), // 1KB
	)
}

// Echo echos the message sent by the client.
func (s *echoAPIHandler) Echo(ctx context.Context, req *connect.Request[echov1.EchoRequest]) (*connect.Response[echov1.EchoResponse], error) {
	zap.L().Info("Request",
		zap.Any("rpc.headers", req.Header()),
		zap.Any("rpc.msg", req.Any()),
	)

	msg := req.Msg.GetMessage()
	res := connect.NewResponse(&echov1.EchoResponse{
		Reply: msg,
	})
	res.Header().Set("EchoAPI-Version", "v1")
	return res, nil
}
