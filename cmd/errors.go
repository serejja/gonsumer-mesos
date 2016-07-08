package cmd

import "errors"

var ErrApiRequired = errors.New("Unspecified gonsumer-mesos API server address. Use --api flag or GM_API env to set.")

var ErrGroupIDRequired = errors.New("Group --id flag is required.")
