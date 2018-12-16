package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/api"
	"github.com/syzoj/syzoj-ng-go/app/auth"
	auth_impl "github.com/syzoj/syzoj-ng-go/app/auth/impl_leveldb"
	"github.com/syzoj/syzoj-ng-go/app/git"
	git_impl "github.com/syzoj/syzoj-ng-go/app/git/impl_leveldb"
	"github.com/syzoj/syzoj-ng-go/app/judge"
	judge_impl "github.com/syzoj/syzoj-ng-go/app/judge/impl_leveldb"
	"github.com/syzoj/syzoj-ng-go/app/problemset"
	problemset_impl "github.com/syzoj/syzoj-ng-go/app/problemset/impl_leveldb"
	"github.com/syzoj/syzoj-ng-go/app/session"
	session_impl "github.com/syzoj/syzoj-ng-go/app/session/impl_leveldb"
)

var log = logrus.StandardLogger()

type syzoj_config struct {
	Database string `json:"database"`
	Addr     string `json:"addr"`
	GitPath  string `json:"git_path"`
	LevelDB  string `json:"leveldb_path"`
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// Parse config
	configPtr := flag.String("config", "config.json", "Sets the config file")

	flag.Parse()

	var err error
	var configData []byte
	if configData, err = ioutil.ReadFile(*configPtr); err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	var config *syzoj_config
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	log.Info("Opening LevelDB")
	var ldb *leveldb.DB
	if ldb, err = leveldb.OpenFile(config.LevelDB, nil); err != nil {
		log.Fatal("Error opening LevelDB: ", err)
	}
	defer func() {
		log.Info("Shutting down LevelDB")
		ldb.Close()
	}()

	log.Info("Setting up session service")
	var sessService session.Service
	if sessService, err = session_impl.NewLevelDBSessionService(ldb); err != nil {
		log.Fatal("Error initializing session service: ", err)
	}
	defer func() {
		log.Info("Shutting down session service")
		sessService.Close()
	}()

	log.Info("Setting up auth service")
	var authService auth.Service
	if authService, err = auth_impl.NewLevelDBAuthService(ldb); err != nil {
		log.Fatal("Error initializing auth service: ", err)
	}
	defer func() {
		log.Info("Shutting down auth service")
		authService.Close()
	}()

	log.Info("Setting up git service")
	var gitService git.GitService
	if gitService, err = git_impl.NewLevelDBGitService(ldb, config.GitPath); err != nil {
		log.Fatal("Error initializing git service: ", err)
	}
	defer func() {
		log.Info("Shutting down git service")
		gitService.Close()
	}()

	log.Info("Setting up judge service")
	var tjudgeService judge.Service
	if tjudgeService, err = judge_impl.NewJudgeService(ldb, gitService); err != nil {
		log.Fatal("Error initializing traditional judge service: ", err)
	}
	defer func() {
		log.Info("Shutting down traditional judge service")
		tjudgeService.Close()
	}()
	gitService.AttachHookHandler("judge", tjudgeService.GetGitHandler())

	log.Info("Setting up problemset service")
	var problemsetService problemset.Service
	if problemsetService = problemset_impl.NewLevelDBProblemset(ldb, tjudgeService); err != nil {
		log.Fatal("Error initializing regular problemset service: ", err)
	}
	defer func() {
		log.Info("Shutting down problemset service")
		problemsetService.Close()
	}()

	log.Info("Setting up api server")
	var apiServer *api.ApiServer
	if apiServer, err = api.CreateApiServer(sessService, authService, problemsetService, tjudgeService); err != nil {
		log.Fatal("Error intializing api server: ", err)
	}

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiServer)
	router.Handle("/judge-traditional", tjudgeService)
	router.PathPrefix("/git").Handler(gitService)
	router.Handle("/", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Addr:         config.Addr,
		Handler:      router,
		WriteTimeout: time.Second * 10,
	}
	go func() {
		log.Infof("Starting web server at %s", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Web server failed unexpectedly: ", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	server.Shutdown(context.Background())
}
