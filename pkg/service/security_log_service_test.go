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
	"github.com/stretchr/testify/require"
	"github.com/codeallergy/template/pkg/service"
	"sync"
	"testing"
	"time"
)

func TestDateFormatter(t *testing.T) {

	var el eventList
	err := el.addEvent()
	require.NoError(t, err)
	err = el.addEvent()
	require.NoError(t, err)

	cnt := 0
	el.eventMap.Range(func(key, value interface{}) bool {
		//fmt.Printf("key = %v\n", key)
		cnt++
		return true
	})

	require.Equal(t, 2, cnt)

}

type eventList struct {
	eventMap sync.Map
}

func (t *eventList) addEvent() error {

	current := time.Now()
tryAgain:
	utc := current.UTC()

	if has, err := t.hasEvent(utc); err != nil {
		return err
	} else if has {
		current = current.Add(time.Millisecond)
		//fmt.Printf("current %v\n", current)
		goto tryAgain
	}

	key := utc.Format(service.DDMMYYYYhhmmss)
	t.eventMap.Store(key, true)
	return nil
}

func (t *eventList) hasEvent(utc time.Time) (bool, error) {
	key := utc.Format(service.DDMMYYYYhhmmss)
	if _ , ok := t.eventMap.Load(key); ok {
		return true, nil
	}
	return false, nil
}