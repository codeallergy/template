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
	"github.com/pkg/errors"
	"github.com/codeallergy/glue"
	"github.com/codeallergy/template/pkg/api"
	"github.com/codeallergy/sprint"
)

type implAdminCommand struct {
	Context           glue.Context             `inject`
	Application       sprint.Application        `inject`
	ApplicationFlags   sprint.ApplicationFlags   `inject`
}

func AdminCommand() sprint.Command {
	return &implAdminCommand{}
}

func (t *implAdminCommand) BeanName() string {
	return "admin"
}

func (t *implAdminCommand) Desc() string {
	return "admin commands: [list, add, remove]"
}

func (t *implAdminCommand) Run(args []string) error {
	if len(args) == 0 {
		return errors.Errorf("admin command needs argument, %s", t.Desc())
	}
	cmd := args[0]
	args = args[1:]

	return doWithAdminClient(t.Context, func(client api.AdminClient) error {
		content, err := client.AdminCommand(cmd, args)
		if err == nil {
			println(content)
		}
		return err
	})

}
