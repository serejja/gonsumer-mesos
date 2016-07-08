package cmd

import (
	"github.com/serejja/gonsumer-mesos/framework"
	"github.com/urfave/cli"
	"os"
)

const (
	FrameworkMasterFlag  = "master"
	FrameworkNameFlag    = "framework-name"
	FrameworkRoleFlag    = "framework-role"
	FrameworkTimeoutFlag = "framework-timeout"
	FrameworkStorageFlag = "storage"
	FrameworkUserFlag    = "user"
	FrameworkBindIPFlag  = "bind-ip"

	ApiFlag = "api"
	ApiEnv  = "GM_API"

	GroupIDFlag               = "id"
	GroupSubscriptionFlag     = "subscription"
	GroupBootstrapBrokersFlag = "bootstrap-brokers"
)

func FrameworkAction(c *cli.Context) error {
	config := framework.NewConfig()
	config.Api = ResolveApi(c)
	if config.Api == "" {
		return ErrApiRequired
	}

	config.Master = c.String(FrameworkMasterFlag)
	config.FrameworkName = c.String(FrameworkNameFlag)
	config.FrameworkRole = c.String(FrameworkRoleFlag)
	config.FrameworkTimeout = c.Duration(FrameworkTimeoutFlag)
	config.FrameworkStorage = c.String(FrameworkStorageFlag)
	config.User = c.String(FrameworkUserFlag)
	config.BindIP = c.String(FrameworkBindIPFlag)

	gonsumerFramework, err := framework.New(config)
	if err != nil {
		return err
	}

	return gonsumerFramework.Start()
}

func ResolveApi(c *cli.Context) string {
	api := os.Getenv(ApiEnv)
	if api == "" {
		return c.String(ApiFlag)
	}

	return api
}
