package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/model/judge"
)

var log = logrus.StandardLogger()

func main() {
	conn, err := grpc.Dial("127.0.0.1:3073", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := judge.NewJudgeServiceClient(conn)
	ctx := context.Background()
	auth := &judge.JudgerAuth{
		JudgerId:    proto.String("k7STX0H9lQ6bsY_R"),
		JudgerToken: proto.String("iq6GvQXcVdw0TEKLEStF3iL7smDyGfj2wJC3_EaEywA="),
	}
	for {
		log.Info("Fetching task")
		resp, err := client.FetchTask(ctx, &judge.FetchTaskRequest{Auth: auth})
		if err != nil {
			log.WithError(err).Error("Failed to fetch task")
			time.Sleep(time.Second * 5)
			continue
		}
		task := resp.Task
		log.Info("Fetched task: ", task)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter result: ")
		text, _ := reader.ReadString('\n')
		_, err = client.HandleTask(ctx, &judge.HandleTaskRequest{Auth: auth, Response: &judge.JudgeResponse{Response: &judge.JudgeResponse_String_{String_: &judge.JudgeStringResponse{Message: proto.String(text)}}}})
		if err != nil {
			log.WithError(err).Error("Failed to handle task")
			continue
		}
	}
	fmt.Println("vim-go")
}
