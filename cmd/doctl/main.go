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
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands/account"
	"github.com/digitalocean/doctl/commands/account/details"
	"github.com/digitalocean/doctl/commands/account/ratelimit"
	"github.com/digitalocean/doctl/commands/compute"
	"github.com/digitalocean/doctl/commands/compute/volume"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/digitalocean/godo"
	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func main() {
	v := viper.New()
	args := os.Args[1:]

	fs := flag.NewFlagSet("doctl-global-flags", flag.ExitOnError)
	// fs.StringP("config", "c", filepath.Join(configHome(), "config.yaml"), "config file")
	// v.BindPFlag("config", fs.Lookup("config"))

	fs.StringP("api-url", "u", "", "Override default API V2 endpoint")
	v.BindPFlag("api-url", fs.Lookup("api-url"))

	fs.StringP(doctl.ArgAccessToken, "t", "", "API V2 Access Token")
	v.BindPFlag(doctl.ArgAccessToken, fs.Lookup(doctl.ArgAccessToken))

	fs.StringP(doctl.ArgOutput, "o", "text", "output format [text|json]")
	v.BindPFlag("output", fs.Lookup(doctl.ArgOutput))

	fs.StringP(doctl.ArgContext, "", "default", "authentication context")
	v.BindPFlag("context", fs.Lookup(doctl.ArgContext))

	fs.BoolP("trace", "", false, "trace api access")
	fs.BoolP(doctl.ArgVerbose, "v", false, "verbose output")

	fs.Parse(args)

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

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := &cli.CLI{
		Args:     args,
		Name:     "doctl",
		HelpFunc: constructHelp(fs),
		Commands: map[string]cli.CommandFactory{
			"account": func() (cli.Command, error) { return account.New(), nil },
			"account get": func() (cli.Command, error) {
				return details.New(v, do.NewAccountService(authClient), displayer, ui), nil
			},
			"account ratelimit": func() (cli.Command, error) {
				return ratelimit.New(v, do.NewAccountService(authClient), displayer, ui), nil
			},
			"compute":             func() (cli.Command, error) { return compute.New(), nil },
			"compute volume":      func() (cli.Command, error) { return volume.New(), nil },
			"compute volume list": func() (cli.Command, error) { return volume.New(), nil },
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		log.Fatalf("failed running cli command: %s", err)
	}

	os.Exit(exitCode)
}

// NewGodoClient needs to support `trace`.
func NewGodoClient(trace bool, url, token string) (*godo.Client, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	return godo.New(oauthClient, godo.SetUserAgent("ME"), godo.SetBaseURL(url))
}

func constructHelp(fs *flag.FlagSet) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {

		var (
			keys   []string
			maxLen int
		)
		for k := range commands {
			keys = append(keys, k)
			if len(k) > maxLen {
				maxLen = len(k)
			}
		}
		sort.Strings(keys)

		var justCommands bytes.Buffer
		for _, key := range keys {
			cmdFunc := commands[key]

			cmd, _ := cmdFunc()

			formattedKey := fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxLen-len(key)))
			justCommands.WriteString(fmt.Sprintf("  %s  %s\n", formattedKey, cmd.Synopsis()))
		}

		return fmt.Sprintf(strings.TrimSpace(help), strings.TrimSpace(justCommands.String()), fs.FlagUsages())
	}
}

const help = `
doctl is a command line interface for the DigitalOcean API.

Usage:
  doctl [command]

Available Commands:
  %s

Flags:
%s

Use "doctl [command] --help" for more information about a command.
`
