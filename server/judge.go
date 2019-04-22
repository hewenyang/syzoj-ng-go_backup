package server

import (
	"container/list"
	"context"
	"errors"
	"io"
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

	mutex         sync.Mutex // hold when changing judger or queue state
	taskCnt       int        // for taskId only
	condQueue     sync.Cond  // notified when judger or queue state changes
	abort         bool       // Whether the judge server is shutting down. When this is set to true, all operations will fail
	judgers       map[database.JudgerRef]*judger
	judgeRequests map[int]*JudgeRequest
	judgerList    list.List // The list of judgers waiting for a task. Type: *judger
	judgeQueue    list.List // The list of jge requests. Type: *JudgeRequest
}

type judger struct {
	*JudgeServer
	id        database.JudgerRef
	condState sync.Cond     // notified when judger state changes
	state     int           // 0 is just initialized, 1 is fetching tasks, 2 is judging, -1 is aborted
	e         *list.Element // Set to nil when judger is removed from queue
	req       *JudgeRequest
}

type JudgeRequest struct {
	s         *JudgeServer
	taskId    int
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

// Methods that must be called with lock held.
func (j *JudgeServer) getJudger(id database.JudgerRef) *judger {
	if o, found := j.judgers[id]; found {
		return o
	}
	o := &judger{}
	o.JudgeServer = j
	o.id = id
	o.condState.L = &j.mutex
	j.judgers[id] = o
	return o
}

func (j *JudgeServer) newJudgeRequest() *JudgeRequest {
	req := &JudgeRequest{}
	req.s = j
	req.taskId = j.taskCnt
	j.taskCnt++
	req.condState.L = &j.mutex
	j.judgeRequests[req.taskId] = req
	return req
}

func (j *JudgeServer) removeJudger(x *judger) {
	if x.state != 1 {
		panic("removeJudger: judger state is not 1")
	}
	j.judgerList.Remove(x.e)
	x.e = nil
}

func (j *JudgeServer) removeRequest(req *JudgeRequest) {
	if req.e == nil {
		panic("removeRequest: request not in queue")
	}
	j.judgeQueue.Remove(req.e)
	req.e = nil
}

func (j *JudgeServer) pushJudger(x *judger) {
	if x.state != 0 {
		panic("pushJudger: judger state is not 0")
	}
	x.state = 1
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

func (j *JudgeRequest) setResponse(resp *judge.JudgeResponse) {
	if j.done {
		panic("setMessage: judge already done")
	}
	j.done = true
	j.resp = resp
	j.condState.Broadcast()
}

func (j *JudgeRequest) setProgress(resp *judge.JudgeResponse) {
	if j.done {
		panic("setMessage: judge already done")
	}
	j.resp = resp
	j.condState.Broadcast()
}

func (j *JudgeRequest) setMessage(msg string) {
	j.setResponse(&judge.JudgeResponse{Response: &judge.JudgeResponse_String_{String_: &judge.JudgeStringResponse{Message: proto.String(msg)}}})
}

func (j *JudgeRequest) abort() {
	if j.e != nil {
		j.s.judgeQueue.Remove(j.e)
		j.e = nil
	}
	if !j.done {
		j.setMessage("judge aborted")
	}
	delete(j.s.judgeRequests, j.taskId)
}

func (j *judger) abort() {
	j.state = -1
	if j.e != nil {
		j.judgerList.Remove(j.e)
	}
	if j.req != nil {
		if !j.req.done {
			j.req.setMessage("judge aborted")
		}
		j.req = nil
	}
	delete(j.judgers, j.id)
}

func (j *JudgeServer) init() {
	j.ctx, j.cancelFunc = context.WithCancel(context.Background())
	j.condQueue.L = &j.mutex
	j.judgers = make(map[database.JudgerRef]*judger)
	j.judgeRequests = make(map[int]*JudgeRequest)
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
				if judger.state != 1 {
					panic("judge: judger in queue but state is not 1")
				}
				j.removeJudger(judger)
				j.removeRequest(req)
				judger.req = req
				judger.state = 2
				judger.condState.Broadcast()
			default:
				j.condQueue.Wait()
			}
		}
		// The queue is aborted, tell all judgers and judge requests
		for _, judger := range j.judgers {
			judger.abort()
		}
		for _, req := range j.judgeRequests {
			req.abort()
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
	j.mutex.Lock()
	if j.abort {
		j.mutex.Unlock()
		return nil, ErrAborted
	}
	j.pushRequest(r)
	j.mutex.Unlock()
	ctx, cancelFunc := context.WithCancel(ctx)
	// Abort the judge request when context is cancelled.
	go func() {
		<-ctx.Done()
		j.mutex.Lock()
		r.abort()
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

func (r *JudgeRequest) GetResult(last *judge.JudgeResponse) (*judge.JudgeResponse, error) {
	r.s.mutex.Lock()
	for r.resp == last && !r.done {
		r.condState.Wait()
	}
	res := r.resp
	r.s.mutex.Unlock()
	if res != last {
		return res, nil
	} else {
		return last, io.EOF
	}
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
	j.mutex.Lock()
	defer j.mutex.Unlock()
	judger := j.getJudger(judgerRef)
	if judger.state != 0 {
		return nil, ErrAborted
	}
	j.pushJudger(judger)
	// Remove judger from queue if context is cancelled.
	ctx2, cancelFunc := context.WithCancel(ctx)
	succ = false
	go func() {
		<-ctx2.Done()
		j.mutex.Lock()
		defer j.mutex.Unlock()
		if !succ {
			judger.abort()
		}
	}()
wait:
	switch judger.state {
	case -1:
		cancelFunc()
		return nil, errors.New("Aborted")
	case 0:
		panic("FetchTask: invalid state")
	case 2:
		succ = true
		cancelFunc()
		return &judge.FetchTaskResponse{Task: judger.req.req}, nil
	default:
		judger.condState.Wait()
		goto wait
	}
}

func (j judgeService) HandleTask(s judge.JudgeService_HandleTaskServer) error {
	ctx := j.ctx
	req, err := s.Recv()
	if err != nil {
		return err
	}
	judgerRef, succ := j.authJudger(ctx, req.Auth)
	if !succ {
		return errors.New("Judger auth failed")
	}
	j.mutex.Lock()
	judger, found := j.judgers[judgerRef]
	if !found {
		return errors.New("HandleTask without a task")
	}
	if judger.state != 2 {
		return errors.New("HandleTask without a task")
	}
	j.mutex.Unlock()
	for {
		req, err := s.Recv()
		if err != nil {
			goto fail
		}
		j.mutex.Lock()
		if judger.state == -1 {
			j.mutex.Unlock()
			return ErrAborted
		}
		if judger.state != 2 {
			panic("HandleTask: state is not 2")
		}
		if judger.req.done {
			break
		}
		if req.GetDone() {
			judger.req.setResponse(req.Response)
			j.mutex.Unlock()
			break
		} else {
			judger.req.setProgress(req.Response)
			j.mutex.Unlock()
		}
	}
	j.mutex.Lock()
	judger.abort()
	j.mutex.Unlock()
	s.SendAndClose(&judge.HandleTaskResponse{})
	return nil
fail:
	j.mutex.Lock()
	judger.abort()
	j.mutex.Unlock()
	return ErrAborted
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
