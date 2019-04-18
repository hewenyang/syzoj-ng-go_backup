package server

import (
	"container/list"
	"context"
	"errors"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model/judge"
)

type JudgeServer struct {
	s *Server

	// The base context for goroutines with the same lifetime as the judge server.
	// Goroutines should check for ctx.Done() signal and shutdown.
	ctx        context.Context
	cancelFunc func()
	wg         sync.WaitGroup

	mutex      sync.Mutex // hold when changing judger or queue state
	condQueue  sync.Cond  // notified when judger or queue state changes
	abort      bool
	judgers    map[database.JudgerRef]*judger
	judgerList list.List // The list of judgers waiting for a task. Type: *judger
	judgeQueue list.List // The list of jge requests. Type: *JudgeRequest
}

type judger struct {
	*JudgeServer
	id        database.JudgerRef
	condState sync.Cond     // notified when judger state changes
	e         *list.Element // Set to nil when judger is removed from queue
	abort     bool
	req       *JudgeRequest
}

type JudgeRequest struct {
	s         *JudgeServer
	condState sync.Cond     // notified when judge request state changes
	e         *list.Element // Set to nil when request is removed from queue
	req       *judge.JudgeRequest
	done      bool
	resp      *judge.JudgeResponse
}

// The new functions does not need locks
func (s *Server) newJudgeServer() *JudgeServer {
	j := &JudgeServer{}
	j.s = s
	j.init()
	return j
}

func (j *JudgeServer) newJudger(id database.JudgerRef) *judger {
	o := &judger{}
	o.id = id
	o.condState.L = &j.mutex
	return o
}

func (j *JudgeServer) newJudgeRequest() *JudgeRequest {
	req := &JudgeRequest{}
	req.s = j
	req.condState.L = &j.mutex
	return req
}

// Primitive methods that must be called with lock held.
func (j *JudgeServer) removeJudger(x *judger) {
	if x.e == nil {
		panic("removeJudger: judger not in queue")
	}
	if x != j.judgerList.Remove(x.e).(*judger) {
		panic("removeJudger: invalid judger list iterator")
	}
	x.e = nil
}

func (j *JudgeServer) removeRequest(req *JudgeRequest) {
	if req.e == nil {
		panic("removeRequest: request not in queue")
	}
	if req != j.judgeQueue.Remove(req.e).(*JudgeRequest) {
		panic("removeRequest: invalid judge request list iterator")
	}
	req.e = nil
}

func (j *JudgeServer) pushJudger(x *judger) {
	if x.e != nil {
		panic("pushJudger: already in queue")
	}
	x.e = j.judgerList.PushBack(x)
	j.condQueue.Broadcast()
}

func (j *JudgeServer) pushRequest(req *JudgeRequest) {
	if req.e != nil {
		panic("pushRequest: request already in queue")
	}
	req.e = j.judgeQueue.PushBack(req)
	j.condQueue.Broadcast()
}

func (req *JudgeRequest) setResponse(resp *judge.JudgeResponse) {
	if req.done {
		panic("setMessage: judge already done")
	}
	req.done = true
	req.resp = resp
	req.condState.Broadcast()
}

func (req *JudgeRequest) setMessage(msg string) {
	req.setResponse(&judge.JudgeResponse{Response: &judge.JudgeResponse_String_{String_: &judge.JudgeStringResponse{Message: proto.String(msg)}}})
}

func (j *JudgeServer) init() {
	j.ctx, j.cancelFunc = context.WithCancel(context.Background())
	j.condQueue.L = &j.mutex
	j.judgers = make(map[database.JudgerRef]*judger)
	// Abort at context done.
	go func() {
		<-j.ctx.Done()
		j.mutex.Lock()
		j.abort = true
		j.condQueue.Broadcast()
		j.mutex.Unlock()
	}()
	// Start a goroutine to match judge requests with judgers.
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		j.mutex.Lock()
		defer j.mutex.Unlock()
	loop:
		for {
			switch {
			case j.abort:
				break loop
			case j.judgerList.Len() > 0 && j.judgeQueue.Len() > 0:
				judger := j.judgerList.Front().Value.(*judger)
				req := j.judgeQueue.Front().Value.(*JudgeRequest)
				j.removeJudger(judger)
				j.removeRequest(req)
				judger.req = req
				judger.condState.Broadcast()
			default:
				j.condQueue.Wait()
			}
		}
		// The queue is aborted, tell all judgers and judge requests
		for {
			switch {
			case j.judgerList.Front() != nil:
				e := j.judgerList.Front()
				judger := e.Value.(*judger)
				judger.abort = true
				judger.condState.Broadcast()
				j.judgerList.Remove(e)
			case j.judgeQueue.Front() != nil:
				e := j.judgeQueue.Front()
				req := e.Value.(*JudgeRequest)
				if !req.done {
					req.setMessage("judge aborted")
				}
				j.judgeQueue.Remove(e)
			default:
				return
			}
		}
	}()
}

func (j *JudgeServer) close() {
	j.cancelFunc()
	j.wg.Wait()
}

// Requests the judge server to judge a submission. The server will find a judger to handle the task.
// The context's lifetime will be considered the request's lifetime.
func (j *JudgeServer) JudgeSubmission(ctx context.Context, req *judge.JudgeRequest) (*JudgeRequest, error) {
	r := j.newJudgeRequest()
	r.req = req
	ctx, cancelFunc := context.WithCancel(ctx)
	j.mutex.Lock()
	j.pushRequest(r)
	j.mutex.Unlock()
	// Abort the judge request when context is cancelled.
	go func() {
		<-ctx.Done()
		j.mutex.Lock()
		// Remove the judge request from the queue.
		if r.e != nil {
			j.removeRequest(r)
		}
		// Cancel the request.
		if !r.done {
			r.setMessage("judge aborted")
		}
		j.mutex.Unlock()
	}()
	go func() {
		j.mutex.Lock()
		defer j.mutex.Unlock()
		for !r.done {
			r.condState.Wait()
		}
		cancelFunc()
	}()
	return r, nil
}

func (r *JudgeRequest) GetResult() (*judge.JudgeResponse, error) {
	r.s.mutex.Lock()
	for !r.done {
		r.condState.Wait()
	}
	res := r.resp
	r.s.mutex.Unlock()
	return res, nil
}

// RPC service
type judgeService struct {
	*JudgeServer
}

func (j *JudgeServer) GetService() judge.JudgeServiceServer {
	return judgeService{j}
}

func (j judgeService) FetchTask(ctx context.Context, req *judge.FetchTaskRequest) (*judge.FetchTaskResponse, error) {
	judgerRef, succ := j.authJudger(ctx, req.Auth)
	if !succ {
		return nil, errors.New("Judger auth failed")
	}
	judger := j.newJudger(judgerRef)
	j.wg.Add(1)
	defer j.wg.Done()
	j.mutex.Lock()
	defer j.mutex.Unlock()
	if _, found := j.judgers[judgerRef]; found {
		return nil, errors.New("Conflicting judger detected")
	}
	j.judgers[judgerRef] = judger
	j.pushJudger(judger)
	// Remove judger from queue if context is cancelled.
	ctx2, cancelFunc := context.WithCancel(ctx)
	go func() {
		<-ctx2.Done()
		j.mutex.Lock()
		defer j.mutex.Unlock()
		if judger.e != nil {
			j.removeJudger(judger)
			delete(j.judgers, judgerRef)
		}
	}()
wait:
	switch {
	case judger.abort:
		cancelFunc()
		return nil, errors.New("Aborted")
	case judger.req != nil:
		cancelFunc()
		return &judge.FetchTaskResponse{Task: judger.req.req}, nil
	default:
		judger.condState.Wait()
		goto wait
	}
}

func (j judgeService) HandleTask(ctx context.Context, req *judge.HandleTaskRequest) (*judge.HandleTaskResponse, error) {
	if req.Response == nil {
		return nil, errors.New("Bad request")
	}
	judgerRef, succ := j.authJudger(ctx, req.Auth)
	if !succ {
		return nil, errors.New("Judger auth failed")
	}
	j.mutex.Lock()
	defer j.mutex.Unlock()
	judger, found := j.judgers[judgerRef]
	if !found {
		return nil, errors.New("HandleTask without a task")
	}
	if judger.req == nil {
		return nil, errors.New("HandleTask without a task")
	}
	r := judger.req
	r.done = true
	r.resp = req.Response
	delete(j.judgers, judgerRef)
	r.condState.Broadcast()
	return &judge.HandleTaskResponse{}, nil
}

//// Helper methods
// Authenticates the judger using JudgerAuth message.
// Returns the judger id and true on success, zero id and false on fail.
func (j *JudgeServer) authJudger(ctx context.Context, auth *judge.JudgerAuth) (database.JudgerRef, bool) {
	if auth == nil {
		return database.JudgerRef(""), false
	}
	judgerId := database.JudgerRef(auth.GetJudgerId())
	judger, err := j.s.GetDB().GetJudger(ctx, judgerId)
	if err != nil {
		log.WithError(err).Error("Failed to get judger")
		return database.JudgerRef(""), false
	}
	if judger == nil {
		return database.JudgerRef(""), false
	}
	if auth.JudgerToken == nil || judger.Token == nil || judger.GetToken() != auth.GetJudgerToken() {
		return database.JudgerRef(""), false
	}
	return judgerId, true
}
