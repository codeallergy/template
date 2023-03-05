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

import (
	"context"
	"github.com/pkg/errors"
	"github.com/codeallergy/store"
	"github.com/codeallergy/template/pkg/api"
	"github.com/codeallergy/template/pkg/pb"
	"github.com/codeallergy/template/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"time"
)


type implSecurityLogService struct {
	Log            *zap.Logger          `inject`
	HostStorage    store.DataStore      `inject:"bean=host-storage"`
	TransactionalManager  store.TransactionalManager  `inject:"bean=host-storage"`

	LogTtl   int   `value:"security-log.ttl,default=31536000"`  // one year ttl
}

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05.000"
)

func SecurityLogService() api.SecurityLogService {
	return &implSecurityLogService{}
}

func (t *implSecurityLogService) LogEvent(ctx context.Context, userId, eventName, remoteIP, userAgent string) (err error) {

	userId = utils.NormalizeUserId(userId)
	if userId == "" {
		return errors.New("userId is empty")
	}

	ctx = t.TransactionalManager.BeginTransaction(ctx, false)
	defer func() {
		err = t.TransactionalManager.EndTransaction(ctx, err)
	}()

	current := time.Now()
tryAgain:
	utc := current.UTC()

	var has bool
	if has, err = t.hasEvent(ctx, userId, utc); err != nil {
		return err
	} else if has {
		current = current.Add(time.Millisecond)
		goto tryAgain
	}

	event := &pb.SecurityLogEntity{
		EventName: eventName,
		EventTime: current.Unix(),
		RemoteIp:  remoteIP,
		UserAgent: userAgent,
	}

	err = t.HostStorage.Set(ctx).ByKey("%s:user:security-log:%s", userId, utc.Format(DDMMYYYYhhmmss)).WithTtl(t.LogTtl).Proto(event)
	return
}

func (t *implSecurityLogService) hasEvent(ctx context.Context, userId string, utc time.Time) (bool, error) {
	event := new(pb.SecurityLogEntity)
	err := t.HostStorage.Get(ctx).ByKey("%s:user:security-log:%s", userId, utc.Format(DDMMYYYYhhmmss)).ToProto(event)
	if err != nil {
		return false, err
	}
	return event.EventName != "", nil
}

func (t *implSecurityLogService) EnumEvents(ctx context.Context, userId string, cb func(item *pb.SecurityLogEntity) bool) error {

	userId = utils.NormalizeUserId(userId)
	if userId == "" {
		return errors.New("userId is empty")
	}

	return t.HostStorage.Enumerate(ctx).ByPrefix("%s:user:security-log:", userId).
		WithBatchSize(BatchSize).
		DoProto(func() proto.Message {
			return new(pb.SecurityLogEntity)
		}, func(entry *store.ProtoEntry) bool {
			if v, ok := entry.Value.(*pb.SecurityLogEntity); ok {
				return cb(v)
			}
			return true
		})

}

