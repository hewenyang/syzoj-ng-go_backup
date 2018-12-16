package impl_leveldb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/judge"
)

func (*judgeService) getProblem(db dbGetter, problemId uuid.UUID) (problem *judge.Problem, err error) {
	var data []byte
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if data, err = db.Get(keyProblem, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = judge.ErrProblemNotExist
		}
		return
	}
	problem = new(judge.Problem)
	if err = json.Unmarshal(data, problem); problem != nil {
		return
	}
	return
}

func (*judgeService) putProblem(db dbPutter, problemId uuid.UUID, problem *judge.Problem) (err error) {
	var data []byte
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if data, err = json.Marshal(problem); err != nil {
		return
	}
	if err = db.Put(keyProblem, data, nil); err != nil {
		return
	}
	return
}

func (*judgeService) deleteProblem(db dbDeleter, problemId uuid.UUID) (err error) {
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if err = db.Delete(keyProblem, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = judge.ErrProblemNotExist
		}
		return
	}
	return
}

func (s *judgeService) CreateProblem(info *judge.Problem) (id uuid.UUID, err error) {
	if id, err = uuid.NewRandom(); err != nil {
		return
	}

	if err = s.putProblem(s.db, id, info); err != nil {
		return
	}
	return
}

func (s *judgeService) InitProblemGit(id uuid.UUID) (gitId uuid.UUID, token string, err error) {
	s.problemLock.Lock()
	defer s.problemLock.Unlock()

	var info *judge.Problem
	if info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	if gitId = info.GitRepo; gitId == (uuid.UUID{}) {
		if gitId, err = s.git.CreateRepository("judge"); err != nil {
			return
		}
		info.GitRepo = gitId
	}
	if token, err = s.git.ResetToken(gitId); err != nil {
		return
	}
	info.GitToken = token
	if err = s.putProblem(s.db, id, info); err != nil {
		return
	}

	return
}

func (s *judgeService) GetProblem(id uuid.UUID) (info *judge.Problem, err error) {
	if info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	return
}

func (s *judgeService) UpdateProblem(id uuid.UUID, info *judge.Problem) (err error) {
	var org_info *judge.Problem
	s.problemLock.Lock()
	defer s.problemLock.Unlock()
	if org_info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	if org_info.Version != info.Version {
		err = judge.ErrConcurrentUpdate
		return
	}
	info.Version++
	if err = s.putProblem(s.db, id, info); err != nil {
		return
	}
	return
}

func (s *judgeService) DeleteProblem(id uuid.UUID) (err error) {
	if err = s.deleteProblem(s.db, id); err != nil {
		return
	}
	return
}
