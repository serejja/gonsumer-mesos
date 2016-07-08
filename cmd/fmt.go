package cmd

import (
	"fmt"
	"github.com/serejja/gonsumer-mesos/framework"
	"strings"
)

func FmtGroups(groups []*framework.Group, indent int) string {
	s := Indent(indent) + "groups:\n"

	stringGroups := make([]string, 0, len(groups))
	for _, group := range groups {
		stringGroups = append(stringGroups, FmtGroup(group, indent+1))
	}

	return s + strings.Join(stringGroups, "\n")
}

func FmtGroup(group *framework.Group, indent int) string {
	s := Indent(indent) + fmt.Sprintf("ID: %s\n", group.ID)
	s += Indent(indent) + fmt.Sprintf("subscription: %s\n", strings.Join(group.Subscriptions, ","))
	s += Indent(indent) + fmt.Sprintf("bootstrap brokers: %s\n", strings.Join(group.BootstrapBrokers, ","))
	s += Indent(indent) + FmtConsumers(group.Consumers, indent+1)

	return s
}

func FmtConsumers(consumers []*framework.Consumer, indent int) string {
	s := Indent(indent) + "consumers:\n"

	stringConsumers := make([]string, 0, len(consumers))
	for _, consumer := range consumers {
		stringConsumers = append(stringConsumers, FmtConsumer(consumer, indent+1))
	}

	return s + strings.Join(stringConsumers, "\n")
}

func FmtConsumer(consumer *framework.Consumer, indent int) string {
	return ""
}

func Indent(indent int) string {
	s := ""
	for i := 0; i < indent; i++ {
		s = s + "  "
	}

	return s
}
