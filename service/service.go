package service

import (
	"context"
	"io/ioutil"
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
	l *logrus.Logger
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
	Migrate(context.Context, *ServiceContext, string) error
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
	return s
}

func (s *Service) Run() {
	c := &ServiceContext{s: s}
	l := logrus.New()
	l.Out = ioutil.Discard
	l.AddHook(&logHook{l: logrus.StandardLogger(), name: s.info.GetName()})
	c.l = l
	go func() {
		log.Infof("Starting service %s", s.info.GetName())
		defer func() {
			log.Infof("Stopped service %s", s.info.GetName())
			s.doStartup()
			s.doStop()
		}()
		s.info.object.Main(s.ctx, c)
	}()
}

func (s *Service) Migrate(prevVersion string) error {
	c := &ServiceContext{s: s}
	l := logrus.New()
	l.Out = ioutil.Discard
	l.AddHook(&logHook{l: logrus.StandardLogger(), name: s.info.GetName()})
	c.l = l
	return s.info.object.Migrate(s.ctx, c, prevVersion)
}

func (s *Service) Stop() {
	s.cancel()
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

func (s *ServiceContext) GetLogger() *logrus.Logger {
	return s.l
}

type logHook struct {
	l    *logrus.Logger
	name string
}

func (h *logHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}
}

func (h *logHook) Fire(e *logrus.Entry) error {
	h.l.WithFields(e.Data).WithField("service", h.name).WithTime(e.Time).Log(e.Level, e.Message)
	return nil
}
