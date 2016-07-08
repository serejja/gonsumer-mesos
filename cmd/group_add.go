package cmd

import (
	"github.com/serejja/gonsumer-mesos/api"
	"github.com/urfave/cli"
)

func GroupAddAction(c *cli.Context) error {
	apiURL := ResolveApi(c)
	if apiURL == "" {
		return ErrApiRequired
	}

	if !c.IsSet(GroupIDFlag) {
		return ErrGroupIDRequired
	}

	groupID := c.String(GroupIDFlag)
	subscription := c.String(GroupSubscriptionFlag)
	bootstrapBrokers := c.String(GroupBootstrapBrokersFlag)

	client := api.NewClient(apiURL)
	return client.AddGroup(groupID, subscription, bootstrapBrokers)
}
