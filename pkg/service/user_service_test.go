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

package service_test

import (
	"context"
	"fmt"
	"github.com/codeallergy/badgerstore"
	"github.com/codeallergy/store"
	"github.com/codeallergy/glue"
	"github.com/codeallergy/template/pkg/api"
	"github.com/codeallergy/template/pkg/pb"
	"github.com/codeallergy/template/pkg/service"
	"github.com/codeallergy/sprintframework/pkg/core"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"strings"
	"testing"
)

func TestUserCRUID(t *testing.T) {

	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	configDir, err := os.MkdirTemp(os.TempDir(), "config-storage-test")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configStore, err := badgerstore.New("config-storage", configDir)
	require.NoError(t, err)
	defer configStore.Destroy()

	hostDir, err := os.MkdirTemp(os.TempDir(), "host-storage-test")
	require.NoError(t, err)
	defer os.RemoveAll(hostDir)

	hostStore, err := badgerstore.New("host-storage", hostDir)
	require.NoError(t, err)
	defer hostStore.Destroy()

	userService := service.UserService()

	ctx, err := glue.New(log, configStore, core.ConfigRepository(1000), hostStore, userService)
	require.NoError(t, err)
	defer ctx.Close()

	verifyUserCRUID(t, userService)
	verifyUserTransactional(t, userService, hostStore)

}

func verifyUserCRUID(t *testing.T, userService api.UserService) {

	ctx := context.Background()

	user, err := userService.CreateUser(ctx, &pb.RegisterRequest{
		FirstName: "Test",
		LastName: "T",
		Email: "test@test.com",
		Password: "test",
	})
	require.NoError(t, err)

	userId, err := userService.GetUserIdByEmail(ctx, "test@test.com")
	require.NoError(t, err)
	require.Equal(t, user.UserId, userId)

	err = userService.DoWithUser(ctx, userId, func(user *pb.UserEntity) error {
		user.LastName = "TT"
		return nil
	})
	require.NoError(t, err)

	user, err = userService.GetUser(ctx, userId)
	require.NoError(t, err)
	require.Equal(t, userId, user.UserId)
	require.Equal(t, "Test", user.FirstName)
	require.Equal(t, "TT", user.LastName)
	require.Equal(t, "test@test.com", user.Email)
	require.NotNil(t, user.PasswordHash)

	user.LastName = "TTT"
	err = userService.SaveUser(ctx, user)
	require.NoError(t, err)

	user, err = userService.AuthenticateUser(ctx, userId, "test")
	require.NoError(t, err)
	require.Equal(t, "TTT", user.LastName)

	var enumUser *pb.UserEntity
	err = userService.EnumUsers(ctx, func(user *pb.UserEntity) bool {
		enumUser = user
		return true
	})
	require.NoError(t, err)
	require.NotNil(t, enumUser)
	require.Equal(t, userId, enumUser.UserId)

	err = userService.RemoveUser(ctx, userId)
	require.NoError(t, err)

	_, err = userService.AuthenticateUser(ctx, userId, "test")
	require.Equal(t, service.ErrUserNotFound, err)

	err = userService.DropUserContent(ctx, userId)
	require.NoError(t, err)

}

func verifyUserTransactional(t *testing.T, userService api.UserService, transactionalManager store.TransactionalManager) {

	ctx := context.Background()

	user, err := userService.CreateUser(ctx, &pb.RegisterRequest{
		FirstName: "Test",
		LastName: "T",
		Email: "test@test.com",
		Password: "test",
	})
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		ctx = transactionalManager.BeginTransaction(context.Background(), false)
		err = userService.RemoveUser(ctx, user.UserId)
		if err == nil {
			err = userService.DropUserContent(ctx, user.UserId)
		}
		err = transactionalManager.EndTransaction(ctx, err)
		if err == nil || !strings.Contains(err.Error(), "concurrent transaction") {
			break
		}
		fmt.Printf("concurrent transaction %d\n", i)
	}

	require.NoError(t, err)

}