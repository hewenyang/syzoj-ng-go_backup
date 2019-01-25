package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_ProblemDb(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	query := bson.D{}
	if c.FormValue("my") != "" {
		if !c.Session.LoggedIn() {
			return ErrNotLoggedIn
		}
		query = append(query, bson.E{"owner", c.Session.AuthUserUid})
	}
	var cursor mongo.Cursor
	if cursor, err = c.Server().mongodb.Collection("problem").Find(c.Context(), query,
		mongo_options.Find().SetProjection(bson.D{{"_id", "1"}, {"title", 1}, {"create_time", 1}}),
	); err != nil {
		panic(err)
	}
	defer cursor.Close(c.Context())

	arena := new(fastjson.Arena)
	result := arena.NewObject()
	problems := arena.NewArray()
	item := 0
	for cursor.Next(c.Context()) {
		var problem model.Problem
		if err = cursor.Decode(&problem); err != nil {
			return
		}
		value := arena.NewObject()
		value.Set("id", arena.NewString(EncodeObjectID(problem.Id)))
		value.Set("title", arena.NewString(problem.Title))
		value.Set("create_time", arena.NewString(problem.CreateTime.String()))
		problems.SetArrayItem(item, value)
		item += 1
        if item >= 100 {
            break
        }
	}
	if err = cursor.Err(); err != nil {
		panic(err)
	}
	result.Set("problems", problems)
	c.SendValue(result)
	return
}
