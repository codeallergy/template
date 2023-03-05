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

package service

import "github.com/pkg/errors"

var (

	ErrNotImplemented = errors.New("not implemented")

	ErrIntegrityDB = errors.New("db integrity")

	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
	ErrUserInvalidPassword = errors.New("wrong password")

	ErrInvalidRecoverCode = errors.New("invalid recover code")

	ErrPageNotFound = errors.New("page not found")
)


