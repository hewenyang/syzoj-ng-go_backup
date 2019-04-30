package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/service"
	apiService "github.com/syzoj/syzoj-ng-go/service/api/service"
	problemService "github.com/syzoj/syzoj-ng-go/service/problem/service"
	userProfileService "github.com/syzoj/syzoj-ng-go/service/user-profile/service"
	userService "github.com/syzoj/syzoj-ng-go/service/user/service"
)

var log = logrus.StandardLogger()

func main() {
	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	wg.Add(1)
	go func() {
		<-c
		cancelFunc()
		signal.Stop(c)
		wg.Done()
	}()
	manager := service.NewServiceManager("services.json")
	manager.AddService(userService.NewUserService(&userService.Config{
		MySQL:        "test:@/test",
		KafkaBrokers: []string{"localhost:9092"},
	}))
	manager.AddService(apiService.NewApiService(&apiService.Config{
		Debug:    true,
		HttpAddr: ":3124",
	}))
	manager.AddService(userProfileService.NewUserProfileService(&userProfileService.Config{
		MySQL:        "test:@/test",
		KafkaBrokers: []string{"localhost:9092"},
	}))
	manager.AddService(problemService.NewProblemService(&problemService.Config{
		MySQL:        "test:@/test",
		KafkaBrokers: []string{"localhost:9092"},
	}))
	if err := manager.Migrate(); err != nil {
		log.WithError(err).Error("Failed to migrate")
		os.Exit(1)
	}
	if err := manager.Run(ctx); err != nil && err != context.Canceled {
		log.WithError(err).Error("Failed to run service")
		os.Exit(1)
	}
	wg.Wait()
	log.Info("Done")
	var numBytes int = 32768
	stack := make([]byte, numBytes)
	for runtime.Stack(stack, true) >= numBytes {
		numBytes *= 2
		stack = make([]byte, numBytes)
	}
	os.Stderr.Write(stack)
}
