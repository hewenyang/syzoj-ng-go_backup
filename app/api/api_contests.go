package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contests(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var cursor *mongo.Cursor
	if cursor, err = c.Server().mongodb.Collection("contest").Find(c.Context(), bson.D{}, mongo_options.Find().SetProjection(bson.D{{"name", 1}})); err != nil {
		panic(err)
	}
	defer cursor.Close(c.Context())
	resp := new(model.ContestsResponse)
	for cursor.Next(c.Context()) {
		entry := new(model.ContestsResponseContestEntry)
		contestModel := new(model.Contest)
		if err = cursor.Decode(contestModel); err != nil {
			panic(err)
		}
		entry.Contest = contestModel
		resp.Contests = append(resp.Contests, entry)
	}
	if err = cursor.Err(); err != nil {
		panic(err)
	}
	c.SendValue(resp)
	return
}
