package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
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
		stream, err := client.HandleTask(ctx)
		if err != nil {
			log.WithError(err).Error("Failed to handle task")
		}
		reader := bufio.NewReader(os.Stdin)
		var res *judge.JudgeResponse
		for {
			fmt.Print("Enter result: ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n")
			var done bool
			if text != "" {
				res = &judge.JudgeResponse{Response: &judge.JudgeResponse_String_{String_: &judge.JudgeStringResponse{Message: proto.String(text)}}}
				done = false
			} else {
				done = true
			}
			err = stream.Send(&judge.HandleTaskRequest{Auth: auth, Response: res, Done: proto.Bool(done)})
			if err != nil {
				log.WithError(err).Error("Failed to handle task")
				continue
			}
			if done {
				break
			}
		}
		resp2, err := stream.CloseAndRecv()
	}
}
