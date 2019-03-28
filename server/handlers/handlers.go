package handlers

import (
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/server"
)

var log = logrus.StandardLogger()

func RegisterHandlers(s *server.ApiServer) {
	router := s.Router()
	loginPage := router.PathPrefix("/api/page/login").Subrouter()
	loginPage.Path("").Methods("GET").Handler(s.WrapHandler(Get_Login, true))
	loginPage.Path("/login").Methods("POST").Handler(s.WrapHandler(Handle_Login, true))
	registerPage := router.PathPrefix("/api/page/register").Subrouter()
	registerPage.Path("").Methods("GET").Handler(s.WrapHandler(Get_Register, true))
	registerPage.Path("/register").Methods("POST").Handler(s.WrapHandler(Handle_Register, true))
    router.Path("/api/page/").Methods("GET").Handler(s.WrapHandler(Get_Index, true))
	router.Path("/api/page/problems").Methods("GET").Handler(s.WrapHandler(Get_Problems, true))
	router.Path("/api/page/problem/create").Methods("GET").Handler(s.WrapHandler(Get_Problem_Create, true))
	router.Path("/api/page/problem/create/create").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Create, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}").Methods("GET").Handler(s.WrapHandler(Get_Problem, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/add-statement").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Add_Statement, true))
	router.Path("/api/page/problem/{problem_id:[0-9A-Za-z-_]{16}}/set-public").Methods("POST").Handler(s.WrapHandler(Handle_Problem_Set_Public, true))
	router.PathPrefix("/api").Methods("GET").Handler(s.WrapHandler(Handle_Not_Found, true))
}
