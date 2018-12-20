package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/auth"
	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/problemset"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type ApiServer struct {
	router            *mux.Router
	sessService       session.Service
	authService       auth.Service
	problemsetService problemset.Service
	judgeService      judge.Service
}

var defaultUserId = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func CreateApiServer(sessService session.Service, authService auth.Service, problemsetService problemset.Service, judgeService judge.Service) (*ApiServer, error) {
	srv := &ApiServer{
		sessService:       sessService,
		authService:       authService,
		problemsetService: problemsetService,
		judgeService:      judgeService,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/auth/register", srv.HandleAuthRegister).Methods("POST")
	router.HandleFunc("/api/auth/login", srv.HandleAuthLogin).Methods("POST")
	router.HandleFunc("/api/auth/logout", srv.HandleAuthLogout).Methods("POST")
	router.HandleFunc("/api/problemset/create", srv.HandleCreateProblemset).Methods("POST")
	router.HandleFunc("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/add", srv.HandleProblemsetAdd).Methods("POST")
	router.HandleFunc("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/list", srv.HandleProblemsetList).Methods("GET")
	router.HandleFunc("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/view", srv.HandleProblemsetView).Methods("GET")
	router.HandleFunc("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/submit", srv.HandleProblemsetSubmit).Methods("POST")
	router.HandleFunc("/api/problem/create", srv.HandleProblemCreate).Methods("POST")
	router.HandleFunc("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/view", srv.HandleProblemView).Methods("GET")
	router.HandleFunc("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/reset-token", srv.HandleResetProblemToken).Methods("POST")
	router.HandleFunc("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/update", srv.HandleProblemUpdate).Methods("POST")
	router.HandleFunc("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/change-title", srv.HandleProblemChangeTitle).Methods("POST")
	srv.router = router
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *ApiServer) ensureSession(w http.ResponseWriter, r *http.Request) (uuid.UUID, *session.Session, error) {
	var sessId uuid.UUID
	if cookie, err := r.Cookie("SYZOJSESSION"); err == nil {
		sessId, _ = uuid.Parse(cookie.Value)
	}
	if sess, err := srv.sessService.GetSession(sessId); err != nil {
		if err != session.ErrSessionNotFound {
			return sessId, sess, err
		}
		if sessId, sess, err := srv.sessService.NewSession(); err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "SYZOJSESSION",
				Value:    sessId.String(),
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24 * 30),
			})
			return sessId, sess, err
		} else {
			return sessId, sess, err
		}
	} else {
		return sessId, sess, nil
	}
}
