/*
Copyright 2018 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands"
	"github.com/digitalocean/doctl/commands/account"
	"github.com/digitalocean/doctl/do"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()

	root := commands.NewRootCmd(v)
	root.PersistentFlags().Parse(os.Args[1:])

	var token string
	switch c := viper.GetString("context"); c {
	case doctl.ArgDefaultContext:
		token = viper.GetString(doctl.ArgAccessToken)
	default:
		contexts := viper.GetStringMapString("auth-contexts")
		token = contexts[c]
	}

	lc := &doctl.LiveConfig{}
	authClient, err := lc.GetGodoClient(viper.GetBool("trace"), token)
	if err != nil {
		//this error is a problem we need to show cobra help plus our error
		fmt.Println(err)
	}

	root.AddCommand(account.NewAccountCmd(v, do.NewAccountService(authClient)))

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
