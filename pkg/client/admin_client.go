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


package client

import (
	"context"
	"github.com/codeallergy/template/pkg/api"
	"github.com/codeallergy/template/pkg/pb"
	"google.golang.org/grpc"
	"sync"
)

type implAdminClient struct {
	GrpcConn   *grpc.ClientConn                `inject`
	client     pb.AdminServiceClient
	closeOnce  sync.Once
}

func AdminClient() api.AdminClient {
	return &implAdminClient{}
}

func (t *implAdminClient) PostConstruct() error {
	t.client = pb.NewAdminServiceClient(t.GrpcConn)
	return nil
}

func (t *implAdminClient) AdminCommand(command string, args []string) (string, error) {

	req := &pb.Command {
		Command: command,
		Args: args,
	}

	if resp, err := t.client.AdminRun(context.Background(), req); err != nil {
		return "", err
	} else {
		return resp.Content, nil
	}
}

func (t *implAdminClient) Destroy() (err error) {
	t.closeOnce.Do(func() {
		if t.GrpcConn != nil {
			err = t.GrpcConn.Close()
		}
	})
	return
}


