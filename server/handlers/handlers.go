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
	router.Path("/api/page/problems").Methods("GET").Handler(s.WrapHandler(Get_Problems, true))
	router.Path("/api/page/problems/add-problem").Methods("POST").Handler(s.WrapHandler(Handle_Problems_Add_Problem, true))
	router.Path("/api/page/problem/create").Methods("GET").Handler(s.WrapHandler(Get_Problem_Create, true))
	router.Path("/api/page/problem/create/create").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Create, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}").Methods("GET").Handler(s.WrapHandler(Get_Problem, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/add-judge-traditional").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Add_Judge_Traditional, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/judge/traditional/submit").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Submit_Judge_Traditional, true))
	router.Path("/api/page/debug/add-judger").Methods("POST").Handler(s.WrapDebugHandler(Handle_Debug_Add_Judger))
	router.PathPrefix("/api").Methods("GET").Handler(s.WrapHandler(Handle_Not_Found, true))
}
