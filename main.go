//go:generate make generate
//go:generate python3 gtag.py MYGTAG assets/

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

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/codeallergy/glue"
	"github.com/codeallergy/template/pkg/assets"
	"github.com/codeallergy/template/pkg/assetsgz"
	"github.com/codeallergy/template/pkg/client"
	"github.com/codeallergy/template/pkg/cmd"
	"github.com/codeallergy/template/pkg/resources"
	"github.com/codeallergy/template/pkg/server"
	"github.com/codeallergy/template/pkg/service"
	"github.com/codeallergy/sprintframework/pkg/app"
	sprintclient "github.com/codeallergy/sprintframework/pkg/client"
	sprintcmd "github.com/codeallergy/sprintframework/pkg/cmd"
	sprintcore "github.com/codeallergy/sprintframework/pkg/core"
	sprintserver "github.com/codeallergy/sprintframework/pkg/server"
	"os"
	"time"
)

var (
	Version string
	Build   string
)

var AppAssets = &glue.ResourceSource{
	Name: "assets",
	AssetNames: assets.AssetNames(),
	AssetFiles: assets.AssetFile(),
}

var AppGzipAssets = &glue.ResourceSource{
	Name: "assets-gzip",
	AssetNames: assetsgz.AssetNames(),
	AssetFiles: assetsgz.AssetFile(),
}

var AppResources = &glue.ResourceSource{
	Name: "resources",
	AssetNames: resources.AssetNames(),
	AssetFiles: resources.AssetFile(),
}

func doMain() (err error) {

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = errors.Errorf("%v", v)
			}
		}
	}()

	return app.Application("template",
		app.WithVersion(Version),
		app.WithBuild(Build),
		app.Beans(app.DefaultApplicationBeans, sprintcmd.DefaultCommands, AppAssets, AppGzipAssets, AppResources, cmd.Commands),
		app.Core(sprintcore.CoreScanner(
			sprintcore.BadgerStorageFactory("config-storage"),
			sprintcore.BadgerStorageFactory("host-storage"),
			sprintcore.LumberjackFactory(),
			sprintcore.AutoupdateService(),
			service.UserService(),
			service.SecurityLogService(),
			service.PageService(),
		)),
		app.Server(sprintserver.ServerScanner(
			sprintserver.AuthorizationMiddleware(),
			sprintserver.GrpcServerFactory("control-grpc-server"),
			sprintserver.ControlServer(),
			server.UIGrpcServer(),
			sprintserver.HttpServerFactory("control-gateway-server"),
			sprintserver.TlsConfigFactory("tls-config"),
		)),
		app.Client(sprintclient.ControlClientScanner(
			sprintclient.AnyTlsConfigFactory("tls-config"),
			client.AdminClient(),
		)),
	).Run(os.Args[1:])

}

func main() {

	if err := doMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(100 * time.Millisecond)
}
