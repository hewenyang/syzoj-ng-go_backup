package handlers

import (
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/server"
)

var log = logrus.StandardLogger()

func RegisterHandlers(s *server.ApiServer) {
	router := s.Router()
	router.Path("/api/page/").Methods("GET").Handler(s.WrapHandler(Get_Index, true))
	router.Path("/api/page/register").Methods("GET").Handler(s.WrapHandler(Get_Register, true))
	router.Path("/api/page/register/register").Methods("POST").Handler(s.WrapHandler(Handle_Register, true))
	router.Path("/api/page/login").Methods("GET").Handler(s.WrapHandler(Get_Login, true))
	router.Path("/api/page/login/login").Methods("POST").Handler(s.WrapHandler(Handle_Login, true))
	router.Path("/api/page/problemset/create").Methods("GET").Handler(s.WrapHandler(Get_Problemset_Create, true))
	router.Path("/api/page/problemset/create/create").Methods("POST").Handler(s.WrapHandler(Handle_Problemset_Create, true))
	router.Path("/api/page/problemset/{problemset_id:[0-9A-Za-z-_]{16}}").Methods("GET").Handler(s.WrapHandler(Get_Problemset, true))
	router.Path("/api/page/problemset/{problemset_id:[0-9A-Za-z-_]{16}}/create-problem").Methods("POST").Handler(s.WrapHandler(Handle_Problemset_Create_Problem, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}").Methods("GET").Handler(s.WrapHandler(Get_Problem, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/remove-judge").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Remove_Judge, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/edit-statement").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Edit_Statement, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/add-judge-traditional").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Add_Judge_Traditional, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/judge/traditional/submit").Methods("GET").Handler(s.WrapWsHandler(Handle_Problem_Submit_Judge_Traditional))
	router.Path("/api/page/debug").Methods("GET").Handler(s.WrapDebugHandler(Get_Debug))
	router.Path("/api/page/debug/add-judger/add-judger").Methods("POST").Handler(s.WrapDebugHandler(Handle_Debug_Add_Judger))
	router.PathPrefix("/api").Methods("GET").Handler(s.WrapHandler(Handle_Not_Found, true))
}
