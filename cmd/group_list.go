package cmd

import (
	"fmt"
	"github.com/serejja/gonsumer-mesos/api"
	"github.com/urfave/cli"
)

func GroupListAction(c *cli.Context) error {
	apiURL := ResolveApi(c)
	if apiURL == "" {
		return ErrApiRequired
	}

	client := api.NewClient(apiURL)
	groups, err := client.ListGroups()
	if err != nil {
		return err
	}

	fmt.Println(FmtGroups(groups, 0))
	return nil
}
