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
	"fmt"
	"github.com/pkg/errors"
	"github.com/codeallergy/template/pkg/pb"
	"github.com/codeallergy/template/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (t *implUIGrpcServer) AdminPageScan(ctx context.Context, req *pb.AdminScanRequest) (resp *pb.AdminPageScanResponse, err error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "AdminPageScan", user.Username)
		}

	}()

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}
	limit := int(req.Limit)

	var total int
	var items []*pb.PageItem
	err = t.PageService.EnumPages(ctx, func(page *pb.PageEntity) bool {
		if offset > 0 {
			offset--
		} else if limit > 0 {
			items = append(items, &pb.PageItem{
				Position:     int32(total + 1),
				Name:         page.Name,
				Title:        page.Title,
				CreatedAt:    page.CreTimestamp,
			})
			limit--
		}
		total++
		return true
	})

	if err != nil {
		return nil, err
	}

	return &pb.AdminPageScanResponse{Items: items, Total: int32(total)}, nil

}

func (t *implUIGrpcServer) AdminCreatePage(ctx context.Context, req *pb.AdminPage) (resp *emptypb.Empty, err error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "AdminCreatePage", user.Username)
		}

	}()

	err = t.PageService.CreatePage(ctx, req)
	return &emptypb.Empty{}, err

}

func (t *implUIGrpcServer) parseContentType(ct string) (pb.ContentType, error) {
	contentType := pb.ContentType_MARKDOWN
	switch strings.ToUpper(strings.TrimSpace(ct)) {
	case "MARKDOWN":
		contentType = pb.ContentType_MARKDOWN
	case "HTML":
		contentType = pb.ContentType_HTML
	default:
		return 0, errors.Errorf("invalid content type '%s'", ct)
	}
	return contentType, nil
}

func (t *implUIGrpcServer) AdminGetPage(ctx context.Context, req *pb.PageName) (*pb.AdminPage, error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	page, err := t.PageService.GetPage(ctx, req.Name)
	if err == service.ErrPageNotFound {
		return nil, status.Errorf(codes.NotFound, "page not found")
	}
	if err != nil {
		return nil, t.wrapError(err, "AdminGetPage", user.Username)
	}

	return &pb.AdminPage{
		Name:        page.Name,
		Title:       page.Title,
		Content:     page.Content,
		ContentType: page.ContentType.String(),
	}, nil

}

func (t *implUIGrpcServer) AdminUpdatePage(ctx context.Context, req *pb.AdminPage) (resp *emptypb.Empty, err error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "AdminSavePage", user.Username)
		}

	}()

	err = t.PageService.UpdatePage(ctx, req)
	return &emptypb.Empty{}, err

}

func (t *implUIGrpcServer) AdminDeletePage(ctx context.Context, req *pb.PageName) (*emptypb.Empty, error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	err := t.PageService.RemovePage(ctx, req.Name)
	if err != nil {
		return nil, t.wrapError(err, "AdminDeletePage", user.Username)
	}

	return &emptypb.Empty{}, nil

}

func (t *implUIGrpcServer) AdminUserScan(ctx context.Context, req *pb.AdminScanRequest) (resp *pb.AdminUserScanResponse, err error) {

	user, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !user.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	defer func() {

		if err != nil {
			err = t.wrapError(err, "AdminUserScan", user.Username)
		}

	}()

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}
	limit := int(req.Limit)

	var total int
	var items []*pb.UserItem
	err = t.UserService.EnumUsers(ctx, func(user *pb.UserEntity) bool {
		if offset > 0 {
			offset--
		} else if limit > 0 {
			items = append(items, &pb.UserItem{
				Position:     int32(total + 1),
				Id:           user.UserId,
				Email:        user.Email,
				FullName:     getFullName(user),
				Role:         user.Role.String(),
				CreatedAt:    user.CreTimestamp,
			})
			limit--
		}
		total++
		return true
	})

	if err != nil {
		return nil, err
	}

	return &pb.AdminUserScanResponse{Items: items, Total: int32(total)}, nil

}

func (t *implUIGrpcServer) AdminGetUser(ctx context.Context, req *pb.UserId) (*pb.AdminUser, error) {

	admin, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !admin.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	user, err := t.UserService.GetUser(ctx, req.Id)
	if err == service.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, t.wrapError(err, "AdminGetUser", req.Id)
	}

	return &pb.AdminUser{
		Id:        user.UserId,
		Email:     user.Email,
		FullName:  getFullName(user),
		CreatedAt: user.CreTimestamp,
		Role: user.Role.String(),
	}, nil

}

func (t *implUIGrpcServer) AdminUpdateUser(ctx context.Context, req *pb.AdminUser) (resp *emptypb.Empty, err error) {

	resp = &emptypb.Empty{}

	admin, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !admin.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	if admin.Username == req.Id {
		return nil, status.Errorf(codes.Unauthenticated, "self role change not permitted")
	}

	var pbRole pb.UserRole
	role := strings.ToUpper(strings.TrimSpace(req.Role))
	switch role {
	case "USER":
		pbRole = pb.UserRole_USER
	case "ADMIN":
		pbRole = pb.UserRole_ADMIN
	default :
		return nil, status.Errorf(codes.InvalidArgument, "unknown role '%s'", role)
	}

	err = t.UserService.DoWithUser(ctx, req.Id, func(user *pb.UserEntity) error {
		user.Role = pbRole
		return nil
	})
	if err != nil {
		err = t.wrapError(err, "AdminUpdateUser", req.Id)
	}

	return
}

func (t *implUIGrpcServer) AdminDeleteUser(ctx context.Context, req *pb.UserId) (resp *emptypb.Empty, err error) {

	resp = &emptypb.Empty{}

	admin, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !admin.Roles["WEB_ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role WEB_ADMIN is required")
	}

	err = t.UserService.RemoveUser(ctx, req.Id)
	if err != nil {
		err = t.wrapError(err, "AdminDeleteUser", req.Id)
	} else {
		err = t.UserService.DropUserContent(context.Background(), req.Id)
		if err != nil {
			err = t.wrapError(err, "DropUserContent", req.Id)
		}
	}

	return
}

func (t *implUIGrpcServer) AdminRun(ctx context.Context, req *pb.Command)  (*pb.CommandResult, error) {

	admin, ok := t.AuthorizationMiddleware.GetUser(ctx)
	if !ok || !admin.Roles["ADMIN"] {
		return nil, status.Errorf(codes.Unauthenticated, "role ADMIN is required")
	}

	switch req.Command {
	case "add":
		return t.setUserRole(ctx, req, pb.UserRole_ADMIN)
	case "remove":
		return t.setUserRole(ctx, req, pb.UserRole_USER)
	case "list":
		var out strings.Builder
		err := t.UserService.EnumUsers(ctx, func(user *pb.UserEntity) bool {
			if user.Role == pb.UserRole_ADMIN {
				out.WriteString(fmt.Sprintf("%s, %s, ADMIN\n", user.Email, getFullName(user)))
			}
			return true
		})
		if err != nil {
			return nil, t.wrapError(err, "AdminRun", admin.Username)
		}
		return &pb.CommandResult{Content: out.String()}, err
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown command '%s', allowed commands 'add,remove,list'", req.Command)
	}

}

func (t *implUIGrpcServer) setUserRole(ctx context.Context, req *pb.Command, role pb.UserRole) (*pb.CommandResult, error) {
	if len(req.Args) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "command needs email argument")
	}
	email := req.Args[0]

	userId, err := t.UserService.GetUserIdByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = t.UserService.DoWithUser(ctx, userId, func(user *pb.UserEntity) error {
		user.Role = role
		return nil
	})
	if err == service.ErrUserNotFound {
		return nil, status.Errorf(codes.NotFound, "user '%s' not found", email)
	}
	if err != nil {
		return nil, t.wrapError(err, "setUserRole", email)
	}

	return &pb.CommandResult{Content: "OK"}, nil
}
