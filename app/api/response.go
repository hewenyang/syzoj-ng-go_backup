package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/app/session"
)

var log = logrus.StandardLogger()

type ApiErrorResponse struct {
	Error string `json:"error"`
}
type ApiErrorResponseWithSession struct {
	Error   string          `json:"error,omitempty"`
	Session SessionResponse `json:"session"`
}

func writeError(w http.ResponseWriter, r *http.Request, err error, sess *session.Session) {
	log.Infof("Error handling request %s: %s", r.URL, err)
	var err2 error
	defer func() {
		if err2 != nil {
			log.WithField("error", err2).Warning("Failed to write error")
		}
	}()
	switch v := err.(type) {
	case *ApiError:
		w.WriteHeader(v.Code)
		encoder := json.NewEncoder(w)
		if sess != nil {
			err2 = encoder.Encode(ApiErrorResponseWithSession{Error: v.Message, Session: getSessionResponse(sess)})
		} else {
			err2 = encoder.Encode(ApiErrorResponse{v.Message})
		}
		
	default:
		encoder := json.NewEncoder(w)
		if sess != nil {
			err2 = encoder.Encode(ApiErrorResponseWithSession{Error: v.Error(), Session: getSessionResponse(sess)})
		} else {
			err2 = encoder.Encode(ApiErrorResponse{v.Error()})
		}
	}
}

type ApiSuccessResponseWithSession struct {
	Data    interface{}     `json:"data"`
	Session SessionResponse `json:"session"`
}
type ApiSuccessResponse struct {
	Data    interface{}     `json:"data"`
}
type SessionResponse struct {
	UserId   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
}

func writeResponse(w http.ResponseWriter, data interface{}, sess *session.Session) {
	encoder := json.NewEncoder(w)
	var err error
	defer func() {
		if err != nil {
			log.WithField("error", err).Warning("Failed to write response")
		}
	}()
	if sess != nil {
		err = encoder.Encode(ApiSuccessResponseWithSession{Data: data, Session: getSessionResponse(sess)})
	} else {
		err = encoder.Encode(ApiSuccessResponse{Data: data})
	}
}

// sess must not be nil
func getSessionResponse(sess *session.Session) SessionResponse {
	return SessionResponse{
		UserId:   sess.AuthUserId,
		UserName: sess.UserName,
	}
}