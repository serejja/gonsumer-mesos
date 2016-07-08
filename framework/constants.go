package framework

import "time"

const (
	DefaultFrameworkMaster  = "127.0.0.1:5050"
	DefaultFrameworkName    = "gonsumer"
	DefaultFrameworkRole    = "*"
	DefaultFrameworkTimeout = 365 * 24 * time.Hour
	DefaultFrameworkStorage = "file:/tmp/gonsumer.json"
)

const (
	ParamGroupID          = "group-id"
	ParamSubscription     = "subscription"
	ParamBootstrapBrokers = "bootstrap-brokers"
)
