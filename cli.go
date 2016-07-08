/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package main

import (
	"fmt"
	"github.com/serejja/gonsumer-mesos/cmd"
	"github.com/serejja/gonsumer-mesos/framework"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gonsumer-mesos"
	app.HelpName = "gonsumer-mesos"
	app.Usage = "Gonsumer Mesos CLI"
	app.UsageText = "gonsumer-mesos command [command options] [arguments...]"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:  "framework",
			Usage: "Launch Mesos framework",
			Flags: []cli.Flag{
				apiFlag,
				cli.StringFlag{
					Name:  cmd.FrameworkMasterFlag,
					Usage: "Mesos Master address in form <ip>:<port>.",
					Value: framework.DefaultFrameworkMaster,
				},
				cli.StringFlag{
					Name:  cmd.FrameworkNameFlag,
					Usage: "Mesos framework name.",
					Value: framework.DefaultFrameworkName,
				},
				cli.StringFlag{
					Name:  cmd.FrameworkRoleFlag,
					Usage: "Mesos framework role.",
					Value: framework.DefaultFrameworkRole,
				},
				cli.DurationFlag{
					Name:  cmd.FrameworkTimeoutFlag,
					Usage: "Mesos framework timeout.",
					Value: framework.DefaultFrameworkTimeout,
				},
				cli.StringFlag{
					Name:  cmd.FrameworkStorageFlag,
					Usage: "Storage for cluster state.",
					Value: framework.DefaultFrameworkStorage,
				},
				cli.StringFlag{
					Name:  cmd.FrameworkUserFlag,
					Usage: "Mesos user. Defaults to current system user.",
				},
				cli.StringFlag{
					Name:  cmd.FrameworkBindIPFlag,
					Usage: "Scheduler driver binding IP address. Optional.",
				},
			},
			Action: cmd.FrameworkAction,
		},
		{
			Name:  "group",
			Usage: "Manage consumer groups",
			Subcommands: []cli.Command{
				{
					Category: "group",
					Name:     "add",
					Usage:    "Add consumer group",
					Action:   cmd.GroupAddAction,
					Flags: []cli.Flag{
						apiFlag,
						cli.StringFlag{
							Name:  cmd.GroupIDFlag,
							Usage: "Group ID to identify a set of consumers. Required.",
						},
						cli.StringFlag{
							Name:  cmd.GroupSubscriptionFlag,
							Usage: "Group subscription expression.",
						},
						cli.StringFlag{
							Name:  cmd.GroupBootstrapBrokersFlag,
							Usage: "Group bootstrap Kafka brokers to discover cluster.",
						},
					},
				},
				{
					Category: "group",
					Name:     "update",
					Usage:    "Update consumer group configuration",
					Action:   cmd.GroupUpdateAction,
				},
				{
					Category: "group",
					Name:     "start",
					Usage:    "Start consumer group",
					Action:   cmd.GroupStartAction,
				},
				{
					Category: "group",
					Name:     "stop",
					Usage:    "Stop consumer group",
					Action:   cmd.GroupStopAction,
				},
				{
					Category: "group",
					Name:     "remove",
					Usage:    "Remove consumer group",
					Action:   cmd.GroupRemoveAction,
				},
				{
					Category: "group",
					Name:     "list",
					Usage:    "List consumer groups",
					Action:   cmd.GroupListAction,
					Flags: []cli.Flag{
						apiFlag,
					},
				},
			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		err := cli.ShowAppHelp(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var apiFlag = cli.StringFlag{
	Name:  cmd.ApiFlag,
	Usage: "host:port address for gonsumer-mesos API server. Required.",
}
