package framework

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	mesos "github.com/mesos/mesos-go/scheduler"
	"github.com/yanzay/log"
	"net"
	"strings"
	"time"
)

type GonsumerFrameworkConfig struct {
	Api              string
	Master           string
	FrameworkName    string
	FrameworkRole    string
	FrameworkStorage string
	FrameworkTimeout time.Duration
	User             string
	BindIP           string
}

func NewConfig() GonsumerFrameworkConfig {
	return GonsumerFrameworkConfig{
		FrameworkName:    "gonsumer",
		FrameworkRole:    "*",
		FrameworkStorage: "file:/tmp/gonsumer.json",
		FrameworkTimeout: 365 * 24 * time.Hour,
		Master:           "127.0.0.1:5050",
	}
}

type Framework struct {
	config    GonsumerFrameworkConfig
	driver    mesos.SchedulerDriver
	scheduler Scheduler
	server    Server
}

func New(config GonsumerFrameworkConfig) (*Framework, error) {
	storage, err := NewStorage(config.FrameworkStorage)
	if err != nil {
		return nil, err
	}

	scheduler, err := NewScheduler(storage)
	if err != nil {
		return nil, err
	}

	driver, err := newSchedulerDriver(scheduler, config)
	if err != nil {
		return nil, err
	}

	server := NewHttpServer(listenAddr(config.Api), scheduler)

	return &Framework{
		config:    config,
		driver:    driver,
		scheduler: scheduler,
		server:    server,
	}, nil
}

func (f *Framework) Start() error {
	go f.server.Start()

	if status, err := f.driver.Run(); err != nil {
		log.Infof("Framework stopped with status %s and error: %s\n", status.String(), err)
		return err
	}

	return nil
}

func newSchedulerDriver(gonsumerScheduler *GonsumerScheduler, config GonsumerFrameworkConfig) (mesos.SchedulerDriver, error) {
	frameworkInfo := &mesosproto.FrameworkInfo{
		User:            proto.String(config.User),
		Name:            proto.String(config.FrameworkName),
		Role:            proto.String(config.FrameworkRole),
		FailoverTimeout: proto.Float64(float64(config.FrameworkTimeout / 1e9)),
		Checkpoint:      proto.Bool(true),
	}

	frameworkID := gonsumerScheduler.Cluster().GetFrameworkID()
	if frameworkID != "" {
		frameworkInfo.Id = util.NewFrameworkID(frameworkID)
	}

	driverConfig := mesos.DriverConfig{
		Scheduler: gonsumerScheduler,
		Framework: frameworkInfo,
		Master:    config.Master,
	}

	if config.BindIP != "" {
		driverConfig.BindingAddress = net.ParseIP(config.BindIP)
	}

	driver, err := mesos.NewMesosSchedulerDriver(driverConfig)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func listenAddr(address string) string {
	if strings.HasPrefix(address, "http://") {
		address = address[len("http://"):]
	}

	colonIndex := strings.LastIndex(address, ":")
	if colonIndex != -1 {
		address = "0.0.0.0" + address[colonIndex:]
	}

	return address
}
