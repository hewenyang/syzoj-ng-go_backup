package service

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

type ServiceManager struct {
	dbPath   string
	info     *serviceManagerData
	services []*ServiceInfo
}

type serviceManagerData struct {
	ServiceVersion map[string]string `json:"service_version"`
}

type serviceManagerService struct {
	Version string `json:"version"`
}

func NewServiceManager(dbPath string) *ServiceManager {
	return &ServiceManager{dbPath: dbPath}
}

func (s *ServiceManager) loadDB() error {
	file, err := os.Open(s.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			s.info = &serviceManagerData{}
			return nil
		}
		return err
	}
	if err := json.NewDecoder(file).Decode(&s.info); err != nil {
		return err
	}
	file.Close()
	return nil
}

func (s *ServiceManager) saveDB() error {
	file, err := os.OpenFile(s.dbPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(s.info); err != nil {
		return err
	}
	file.Close()
	return nil
}

func (s *ServiceManager) Migrate() error {
	if err := s.loadDB(); err != nil {
		return err
	}
	if s.info.ServiceVersion == nil {
		s.info.ServiceVersion = make(map[string]string)
	}
	ctx := context.Background()
	for _, i := range s.services {
		v := s.info.ServiceVersion[i.GetName()]
		if i.GetVersion() != v {
			sv := NewService(ctx, i)
			err := sv.Migrate(v)
			if err != nil {
				log.WithError(err).Errorf("Failed to migrate service %s", i.GetName())
				s.saveDB()
				return err
			}
			s.info.ServiceVersion[i.GetName()] = i.GetVersion()
		}
	}
	return s.saveDB()
}

func (s *ServiceManager) AddService(info *ServiceInfo) {
	s.services = append(s.services, info)
}

func (s *ServiceManager) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, i := range s.services {
		sv := NewService(ctx, i)
		sv.Run()
		<-sv.StartupChan()
		wg.Add(1)
		go func() {
			<-sv.StopChan()
			wg.Done()
		}()
	}
	wg.Wait()
	return ctx.Err()
}
