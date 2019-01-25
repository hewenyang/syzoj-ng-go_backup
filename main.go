package main

import (
	"context"
    "database/sql"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
    "fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
    _ "github.com/go-sql-driver/mysql"

	"github.com/syzoj/syzoj-ng-go/app/api"
	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/contest"
	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
    "github.com/syzoj/syzoj-ng-go/tool/import"
)

var log = logrus.StandardLogger()

type syzoj_config struct {
	Mongo     string     `json:"mongo"`
	Addr      string     `json:"addr"`
	JudgeAddr string     `json:"judge_addr"`
	Api       api.Config `json:"api_server"`
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
    if len(os.Args) <= 1 {
        fmt.Println("Must specify a subcommand: run or import")
        return
    }
    switch os.Args[1] {
    case "run":
        cmdRun()
    case "import":
        cmdImport()
    default:
        fmt.Println("Must specify a subcommand: run or import")
        return
    }
}

func cmdImport() {
    importFlagSet := flag.NewFlagSet("import", flag.ExitOnError)
    configPtr := importFlagSet.String("config", "config.json", "Sets the config file")
    mysqlPtr := importFlagSet.String("mysql", "root:@/syzoj", "MySQL database to import from")
    importFlagSet.Parse(os.Args[2:])
    
	var err error
	var configData []byte
	if configData, err = ioutil.ReadFile(*configPtr); err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	var config *syzoj_config
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Fatal("Error parsing config file: ", err)
	}
   
	log.Info("Connecting to MongoDB")
	var mongoClient *mongo.Client
	if mongoClient, err = mongo.Connect(context.Background(), config.Mongo); err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		log.Info("Disconnecting from MongoDB")
		mongoClient.Disconnect(ctx)
	}()

    log.Info("Connecting to MySQL")
    var mysql *sql.DB
    if mysql, err = sql.Open("mysql", *mysqlPtr); err != nil {
        log.Fatal("Error connecting to MySQL: ", err)
    }
    defer func() {
        log.Info("Disconnecting from MySQL")
        mysql.Close()
    }()

    tool_import.ImportMySQL(mongoClient, mysql)
}

func cmdRun() {
    runFlagSet := flag.NewFlagSet("run", flag.ExitOnError)
    configPtr := runFlagSet.String("config", "config.json", "Sets the config file")
	runFlagSet.Parse(os.Args[2:])

	var err error
	var configData []byte
	if configData, err = ioutil.ReadFile(*configPtr); err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	var config *syzoj_config
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	log.Info("Connecting to MongoDB")
	var mongoClient *mongo.Client
	if mongoClient, err = mongo.Connect(context.Background(), config.Mongo); err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		log.Info("Disconnecting from MongoDB")
		mongoClient.Disconnect(ctx)
	}()

	var grpcServer *grpc.Server = grpc.NewServer()

	log.Info("Start judge service")
	var judgeService judge.Service
	if judgeService, err = judge.CreateJudgeService(mongoClient); err != nil {
		log.Fatal("Error starting judge service: ", err)
	}
	defer func() {
		log.Info("Stopping judge service")
		judgeService.Close()
	}()
	judge_api.RegisterJudgeServer(grpcServer, judgeService)

	reflection.Register(grpcServer)
	go func() {
		log.Info("Setting up gRPC service")
		lis, err := net.Listen("tcp", "0.0.0.0:3073")
		if err != nil {
			log.Fatal("Failed to listen: ", err)
		}
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC service: ", err)
		}
	}()

    log.Info("Start contest service")
    var contestService = contest.NewContestService(mongoClient)
    if err = contestService.Init(context.Background()); err != nil {
        log.Fatal("Failed to start contest service")
    }
    defer func() {
        log.Info("Stopping contest service")
        contestService.Close()
    }()

	log.Info("Setting up api server")
	var apiServer *api.ApiServer
	if apiServer, err = api.CreateApiServer(mongoClient, judgeService, contestService, config.Api); err != nil {
		log.Fatal("Error intializing api server: ", err)
	}

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiServer)
	//router.Handle("/judge-traditional", tjudgeService)
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
