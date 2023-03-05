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
	"crypto/tls"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"
	"github.com/codeallergy/store"
	"github.com/codeallergy/glue"
	"github.com/codeallergy/template/pkg/api"
	"github.com/codeallergy/template/pkg/pb"
	"github.com/codeallergy/template/pkg/service"
	"github.com/codeallergy/sprint"
	"github.com/codeallergy/sprintframework/pkg/util"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
	"time"
)


type implUIGrpcServer struct {
	pb.UnimplementedAuthServiceServer
	pb.UnimplementedSiteServiceServer
	pb.UnimplementedAdminServiceServer

	WebappName string `value:"webapp.name,default=Light-Template"`

	GrpcServer       *grpc.Server   `inject`
	UIGatewayServer  *http.Server   `inject:"bean=control-gateway-server"`

	ControlTls     bool     `value:"control-grpc-server.tls,default=false"`
	GrpcAddress    string   `value:"control-grpc-server.listen-address,default="`

	Properties            glue.Properties    `inject`
	AuthorizationMiddleware sprint.AuthorizationMiddleware `inject`
	NodeService           sprint.NodeService  `inject`
	MailService           sprint.MailService  `inject`

	UserService           api.UserService   `inject`
	SecurityLogService    api.SecurityLogService  `inject`
	PageService           api.PageService   `inject`
	TransactionalManager  store.TransactionalManager  `inject:"bean=host-storage"`

	Log             *zap.Logger          `inject`

	loginCnt        atomic.Int64
	registerCnt     atomic.Int64
	restoreCnt      atomic.Int64

	AccessTokenMinutes   int   `value:"auth.access-token-minutes,default=20"`
	RefreshTokenHours    int   `value:"auth.refresh-token-hours,default=24"`
}

func UIGrpcServer() api.GRPCServer {
	return &implUIGrpcServer{}
}

func (t *implUIGrpcServer) BeanName() string {
	return "ui_grpc_server"
}

func (t *implUIGrpcServer) PostConstruct() (err error) {
	pb.RegisterAuthServiceServer(t.GrpcServer, t)
	pb.RegisterSiteServiceServer(t.GrpcServer, t)
	pb.RegisterAdminServiceServer(t.GrpcServer, t) // no gateway

	api, err := util.FindGatewayHandler(t.UIGatewayServer, "/api/")
	if err != nil {
		return err
	}

	// interceptors do not work yet
	//pb.RegisterAuthServiceHandlerServer(context.Background(), api, t)

	var opts []grpc.DialOption

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	tlsConfig.NextProtos = appendH2ToNextProtos(tlsConfig.NextProtos)

	tlsCredentials := credentials.NewTLS(tlsConfig)
	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))

	if t.GrpcAddress == "" {
		return errors.New("property 'control-grpc-server.listen-address' is empty")
	}

	pb.RegisterAuthServiceHandlerFromEndpoint(context.Background(), api, t.GrpcAddress, opts)
	pb.RegisterSiteServiceHandlerFromEndpoint(context.Background(), api, t.GrpcAddress, opts)

	return nil
}

func (t *implUIGrpcServer) GetStats(cb func(name, value string) bool) error {

	cb("login.cnt", strconv.FormatInt(t.loginCnt.Load(), 10))
	cb("register.cnt", strconv.FormatInt(t.registerCnt.Load(), 10))
	cb("restore.cnt", strconv.FormatInt(t.restoreCnt.Load(), 10))

	return nil
}

func (t *implUIGrpcServer) Page(ctx context.Context, req *pb.PageName) (*pb.PageContent, error) {

	page, err := t.PageService.GetPage(ctx, req.Name)
	if err == service.ErrPageNotFound {
		return &pb.PageContent{
			Title:   "Page Not Found",
			Content: fmt.Sprintf("Oops, requested page '%s' is not found.", req.Name),
		}, nil
	}

	defer func() {

		if err != nil {
			id := t.NodeService.Issue().String()
			t.Log.Error("Page", zap.String("errorId", id), zap.Error(err))
			err = status.Errorf(codes.Internal, "internal error %s", id)
		}

	}()

	if err != nil {
		return nil, err
	}

	var content string
	if page.ContentType == pb.ContentType_MARKDOWN {
		content = string(markdown.ToHTML([]byte(page.Content), nil, nil))
	} else {
		content = page.Content
	}

	return &pb.PageContent{Title: page.Title, Content: content }, nil
}


func (t *implUIGrpcServer) UserDelete(ctx context.Context, req *pb.UserId) (resp *emptypb.Empty, err error) {

	resp = &emptypb.Empty{}

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_USER"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_USER is required")
	}

	if user.Username != req.Id {
		return nil, status.Errorf(codes.Unauthenticated, "logged in user '%s' can not delete user id '%s'", user.Username, req.Id)
	}

	entity, err := t.UserService.GetUser(ctx, req.Id)
	if err == service.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	err = t.UserService.RemoveUser(ctx, req.Id)
	if err != nil {
		err = t.wrapError(err, "UserDelete", req.Id)
	} else {
		err = t.UserService.DropUserContent(context.Background(), req.Id)
		if err != nil {
			err = t.wrapError(err, "UserDelete", req.Id)
		}
	}

	subject := fmt.Sprintf("Goodbye %s.", entity.FirstName)
	sender := t.Properties.GetString("mail.sender", "noreply@localhost")

	mail := sprint.Mail{
		Sender:      sender,
		Recipients:   []string{entity.Email},
		Subject:      subject,
		TextTemplate: "resources:mail/deleted_user_text.tmpl",
		HtmlTemplate: "resources:mail/deleted_user_html.tmpl",
		Data:         map[string]interface{} {
			"FirstName": entity.FirstName,
			"Project": t.WebappName,
		},
	}

	go t.MailService.SendMail(&mail, time.Minute, false)

	return
}
