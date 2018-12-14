package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type ProblemsetCreateRequest struct{}
type AddTraditionalProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	Name         string    `json:"name"`
}
type SubmitTraditionalProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	ProblemName  string    `json:"problem_name"`
	Language     string    `json:"language"`
	Code         string    `json:"code"`
}

func Test(client *http.Client, method string, ep string, data interface{}) (result interface{}) {
	fmt.Printf("Doing %s to %s with data %+v\n", method, ep, data)
	var body io.Reader
	if data != nil {
		bodyContent, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(bodyContent)
	}
	req, err := http.NewRequest(method, "http://127.0.0.1:5900"+ep, body)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Printf("Result: %+v\n", result)
	fmt.Printf("Headers: %+v\n", resp.Header)
	u, _ := url.Parse("http://127.0.0.1:5900")
	client.Jar.SetCookies(u, resp.Cookies())
	return
}

func main() {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	Test(client, "POST", "/api/auth/register", RegisterRequest{"aaa", "B"})
	Test(client, "POST", "/api/auth/login", LoginRequest{"aaa", "B"})
	u, _ := url.Parse("http://127.0.0.1:5900/")
	fmt.Println(cookieJar.Cookies(u))
	resp1 := Test(client, "POST", "/api/problemset/create", ProblemsetCreateRequest{})
	problemsetId := uuid.MustParse(resp1.(map[string]interface{})["data"].(map[string]interface{})["problemset_id"].(string))
	fmt.Println(resp1)
	fmt.Println(problemsetId)
	resp2 := Test(client, "POST", "/api/problemset/add-traditional", AddTraditionalProblemRequest{
		ProblemsetId: problemsetId,
		Name:         "1",
	})
	problemId := uuid.MustParse(resp2.(map[string]interface{})["data"].(map[string]interface{})["problem_id"].(string))
	fmt.Println(problemId)
	Test(client, "POST", "/api/problemset/submit-traditional", SubmitTraditionalProblemRequest{
		ProblemsetId: problemsetId,
		ProblemName:  "1",
		Language:     "cpp",
		Code:         "hellello",
	})
}
