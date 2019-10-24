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
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands"
	"github.com/digitalocean/doctl/commands/account"
	"github.com/digitalocean/doctl/commands/compute"
	"github.com/digitalocean/doctl/commands/compute/droplet"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func main() {
	v := viper.New()

	root := commands.NewRootCmd(v)
	root.PersistentFlags().Parse(os.Args[1:])

	var token string
	switch c := v.GetString("context"); c {
	case doctl.ArgDefaultContext:
		token = v.GetString(doctl.ArgAccessToken)
	default:
		contexts := v.GetStringMapString("auth-contexts")
		t, ok := contexts[c]
		if !ok {
			panic("HOW????")
		}

		token = t
	}

	authClient, err := NewGodoClient(v.GetBool("trace"), v.GetString("api-url"), token)
	if err != nil {
		//this error is a problem we need to show cobra help plus our error
		fmt.Println(err)
	}

	displayer := displayers.Displayer{
		OutputType: v.GetString(doctl.ArgOutput),
		Out:        os.Stdout,
	}

	root.AddCommand(account.NewAccountCmd(v, do.NewAccountService(authClient), displayer))

	computeCmd := compute.NewCommand()
	computeCmd.AddCommand(droplet.NewCommand(v, do.NewDropletsService(authClient), displayer))

	root.AddCommand(computeCmd)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewGodoClient needs to support `trace`.
func NewGodoClient(trace bool, url, token string) (*godo.Client, error) {
	if token == "" {
		return nil, fmt.Errorf("access token is required. (hint: run 'doctl auth init')")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	return godo.New(oauthClient, godo.SetUserAgent("ME"), godo.SetBaseURL(url))
}
