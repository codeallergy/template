/*
 * Copyright (c) 2022-2023 Zander Schwid & Co. LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */


package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/codeallergy/template/pkg/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func getCallerInfo(ctx context.Context) (string, string) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ""
	}

	return getRemoteAddress(md), getUserAgent(md)
}

func getRemoteAddress(md metadata.MD) string {

	headers, ok := md["x-forwarded-for"]
	if !ok {
		return ""
	}

	if len(headers) == 0 {
		return ""
	}

	return headers[0]
}

func getUserAgent(md metadata.MD) string {

	headers, ok := md["grpcgateway-user-agent"]
	if !ok {
		return ""
	}

	if len(headers) == 0 {
		return ""
	}

	return headers[0]
}

func getFullName(user *pb.UserEntity) string {
	var out strings.Builder
	if user.FirstName != "" {
		out.WriteString(user.FirstName)
	}
	if user.MiddleName != "" {
		if out.Len() > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(user.MiddleName)
	}
	if user.LastName != "" {
		if out.Len() > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(user.LastName)
	}
	return out.String()
}

const alpnProtoStrH2 = "h2"

func appendH2ToNextProtos(ps []string) []string {
	for _, p := range ps {
		if p == alpnProtoStrH2 {
			return ps
		}
	}
	ret := make([]string, 0, len(ps)+1)
	ret = append(ret, ps...)
	return append(ret, alpnProtoStrH2)
}

func (t *implUIGrpcServer) wrapError(err error, method, username string) error {
	if _, ok := status.FromError(err); ok {
		return err
	}
	issue := err.Error()
	if strings.HasPrefix(issue, "nowrap:") {
		issue = strings.TrimSpace(strings.TrimPrefix(issue, "nowrap:"))
		return errors.New(issue)
	}
	message := "internal error"
	if strings.Contains("concurrent transaction", issue) {
		message = "concurrent transaction"
	} else if strings.Contains("not found", issue) {
		message = "object not found"
	} else if strings.Contains("exist", issue) {
		message = "object already exist"
	}
	id := t.NodeService.Issue().String()
	t.Log.Error(method, zap.String("errorId", id), zap.Any("username", username), zap.Error(err))
	return status.Errorf(codes.Internal, "%s %s", message, id)
}
