package framework

import (
	"encoding/json"
	"sync"
)

type Cluster interface {
	SetFrameworkID(id string)
	GetFrameworkID() string

	AddGroup(group *Group)
	GetGroup(id string) *Group
	ExistsGroup(id string) bool
	GetGroups() []*Group
}

type gonsumerClusterJSON struct {
	FrameworkID string   `json:"framework_id"`
	Groups      []*Group `json:"groups"`
}

type GonsumerCluster struct {
	lock sync.Mutex

	frameworkID string
	groups      map[string]*Group
}

func NewGonsumerCluster() *GonsumerCluster {
	return &GonsumerCluster{
		groups: make(map[string]*Group),
	}
}

func (c *GonsumerCluster) SetFrameworkID(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.frameworkID = id
}

func (c *GonsumerCluster) GetFrameworkID() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.frameworkID
}

func (c *GonsumerCluster) AddGroup(group *Group) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.groups[group.ID] = group
}

func (c *GonsumerCluster) GetGroup(id string) *Group {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.groups[id]
}

func (c *GonsumerCluster) ExistsGroup(id string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, exists := c.groups[id]
	return exists
}

func (c *GonsumerCluster) GetGroups() []*Group {
	c.lock.Lock()
	defer c.lock.Unlock()

	groups := make([]*Group, 0, len(c.groups))
	for _, group := range c.groups {
		groups = append(groups, group)
	}

	return groups
}

func (c *GonsumerCluster) MarshalJSON() ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cluster := gonsumerClusterJSON{
		FrameworkID: c.frameworkID,
		Groups:      make([]*Group, 0, len(c.groups)),
	}

	for _, group := range c.groups {
		cluster.Groups = append(cluster.Groups, group)
	}

	return json.Marshal(cluster)
}

func (c *GonsumerCluster) UnmarshalJSON(bytes []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var cluster gonsumerClusterJSON
	err := json.Unmarshal(bytes, &cluster)
	if err != nil {
		return err
	}

	c.frameworkID = cluster.FrameworkID
	c.groups = make(map[string]*Group)
	for _, group := range cluster.Groups {
		c.groups[group.ID] = group
	}

	return nil
}

type Group struct {
	ID               string   `json:"id"`
	Subscriptions    []string `json:"subscriptions"`
	BootstrapBrokers []string `json:"bootstrap_brokers"`

	Consumers []*Consumer `json:"consumers"`
}

type Consumer struct {
	ID string `json:"id"`
	//Assignments
}
