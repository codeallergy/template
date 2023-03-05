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

package cmd

import (
	"github.com/codeallergy/glue"
	"github.com/codeallergy/template/pkg/api"
	"github.com/codeallergy/sprint"
	"github.com/pkg/errors"
	"log"
)

func doWithAdminClient(parent glue.Context, cb func(client api.AdminClient) error) error {

	var verbose bool
	list := parent.Bean(sprint.ApplicationFlagsClass, glue.DefaultLevel)
	if len(list) > 0 {
		flags := list[0].Object().(sprint.ApplicationFlags)
		if flags.Verbose() {
			verbose = true
	}
	}

	list = parent.Bean(sprint.ClientScannerClass, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("application context should have one client scanner, but found '%d'", len(list))
	}
	bean := list[0]

	scanner, ok := bean.Object().(sprint.ClientScanner)
	if !ok {
		return errors.Errorf("invalid object '%v' found instead of sprint.ClientScanner in application context", bean.Class())
	}

	beans := scanner.ClientBeans()
	if verbose {
		verbose := glue.Verbose{ Log: log.Default() }
		beans = append([]interface{}{verbose}, beans...)
	}

	ctx, err := parent.Extend(beans...)
	if err != nil {
		return err
	}
	defer ctx.Close()

	list = ctx.Bean(api.AdminClientClass, glue.DefaultLevel)
	if len(list) != 1 {
		return errors.Errorf("client context shoulw have one api.AdminClient inference, but found '%d'", len(list))
	}
	bean = list[0]

	if client, ok := bean.Object().(api.AdminClient); ok {
		return cb(client)
	} else {
		return errors.Errorf("invalid object '%v' found instead of api.AdminClient in client context", bean.Class())
	}

}