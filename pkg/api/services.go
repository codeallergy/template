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

package api

import (
	"context"
	"github.com/codeallergy/store"
	"github.com/codeallergy/glue"
	"github.com/codeallergy/template/pkg/pb"
	"reflect"
)


var UserServiceClass = reflect.TypeOf((*UserService)(nil)).Elem()

type UserService interface {
	glue.InitializingBean

	GenerateUserId(ctx context.Context) (string, error)

	CreateUser(ctx context.Context, req *pb.RegisterRequest) (*pb.UserEntity, error)

	ResetPassword(ctx context.Context, email string, newPassword string) (string, error)

	// could be userId or email
	AuthenticateUser(ctx context.Context, username, password string) (*pb.UserEntity, error)

	GetUser(ctx context.Context, userId string) (*pb.UserEntity, error)

	GetUserIdByEmail(ctx context.Context, email string) (string, error)

	SaveUser(ctx context.Context, user *pb.UserEntity) error

	RemoveUser(ctx context.Context, userId string) error

	DropUserContent(ctx context.Context, userId string) error

	DoWithUser(ctx context.Context, userId string, cb func(user *pb.UserEntity) error) error

	DumpUser(ctx context.Context, userId string, cb func(entry *store.RawEntry) bool) error

	EnumUsers(ctx context.Context, cb func(user *pb.UserEntity) bool) error

	SaveRecoverCode(ctx context.Context, email string, rc *pb.RecoverCodeEntity, ttlSeconds int) error

	ValidateRecoverCode(ctx context.Context, email string, code string) error
}

var SecurityLogServiceClass = reflect.TypeOf((*SecurityLogService)(nil)).Elem()

type SecurityLogService interface {

	LogEvent(ctx context.Context, userId, eventName, remoteIP, userAgent string) error

	EnumEvents(ctx context.Context, userId string, cb func(item *pb.SecurityLogEntity) bool) error

}

var PageServiceClass = reflect.TypeOf((*PageService)(nil)).Elem()

type PageService interface {

	// ErrPageNotFound on error
	GetPage(ctx context.Context, name string) (*pb.PageEntity, error)

	CreatePage(ctx context.Context, page *pb.AdminPage) error

	UpdatePage(ctx context.Context, page *pb.AdminPage) error

	RemovePage(ctx context.Context, name string) error

	EnumPages(ctx context.Context, cb func(page *pb.PageEntity) bool) error

}
