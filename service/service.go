package service

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type Service struct {
	info        *ServiceInfo
	ctx         context.Context
	cancel      func()
	startup     chan struct{}
	startupOnce sync.Once
	stop        chan struct{}
	stopOnce    sync.Once
}

type ServiceContext struct {
	s *Service
}

// ServiceBuilder is a struct used to build ServiceInfo.
type ServiceBuilder struct {
	Name    string
	Version string
	Object  ServiceObject
}

// ServiceInfo describes a service. It is immutable.
type ServiceInfo struct {
	name    string
	version string
	object  ServiceObject
}

type ServiceObject interface {
	Main(context.Context, *ServiceContext)
	Migrate(string) error
}

func (s ServiceBuilder) Build() *ServiceInfo {
	if s.Name == "" {
		panic("ServiceBuilder: Name is empty")
	}
	if s.Version == "" {
		panic("ServiceBuilder: Version is empty")
	}
	if s.Object == nil {
		panic("ServiceBuilder: Main is nil")
	}
	return &ServiceInfo{
		name:    s.Name,
		version: s.Version,
		object:  s.Object,
	}
}

func (s *ServiceInfo) GetName() string {
	return s.name
}

func (s *ServiceInfo) GetVersion() string {
	return s.version
}

func NewService(ctx context.Context, info *ServiceInfo) *Service {
	s := &Service{}
	s.info = info
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.startup = make(chan struct{})
	s.stop = make(chan struct{})
	c := &ServiceContext{s: s}
	go func() {
		log.Infof("Starting service %s", info.GetName())
		defer func() {
			log.Infof("Stopped service %s", info.GetName())
			s.doStartup()
			s.doStop()
		}()
		info.object.Main(s.ctx, c)
	}()
	return s
}

func (s *Service) StartupChan() <-chan struct{} {
	return s.startup
}

func (s *Service) StopChan() <-chan struct{} {
	return s.stop
}

func (s *Service) doStartup() {
	s.startupOnce.Do(func() { close(s.startup) })
}

func (s *Service) doStop() {
	s.stopOnce.Do(func() { close(s.stop) })
}

func (s *ServiceContext) StartupDone() {
	s.s.doStartup()
}
