package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func Test(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// 初始化
	mdb := client.Database("webook")
	col := mdb.Collection("articles")

	// 索引
	idx, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(idx)

	// 增
	res, err := col.InsertOne(ctx, Article{
		Id:      123,
		Title:   "title",
		Content: "content",
	})
	// _id ：自身文档字段
	fmt.Printf("id %d \n", res.InsertedID)

	// 查
	filter := bson.D{
		bson.E{Key: "id", Value: 123},
	}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("No document was found")
			return
		}
		panic(err)
	}
	fmt.Printf("%+v \n", art)

	// or 查询
	or := bson.A{
		bson.D{bson.E{Key: "id", Value: 123}},
		bson.D{bson.E{Key: "id", Value: 456}},
	}
	or = bson.A{bson.M{"id": 123}, bson.M{"id": 456}}
	orRes, err := col.Find(ctx, bson.D{bson.E{Key: "$or", Value: or}})
	if err != nil {
		panic(err)
	}
	var ars []Article
	err = orRes.All(ctx, &ars)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", ars)

	// and 查询
	and := bson.A{
		bson.D{bson.E{Key: "id", Value: 123}},
		bson.D{bson.E{Key: "title", Value: "title"}},
	}
	andRes, err := col.Find(ctx, bson.D{bson.E{Key: "$and", Value: and}})
	if err != nil {
		panic(err)
	}
	ars = []Article{}
	err = andRes.All(ctx, &ars)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", ars)

	// In 查询
	in := bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: bson.A{123, 456}}}}}
	in = bson.D{bson.E{Key: "id", Value: bson.M{"$in": bson.A{123, 456}}}}
	inRes, err := col.Find(ctx, in)
	if err != nil {
		panic(err)
	}
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", ars)

	// projecion 只查某些字段
	inRes, err = col.Find(ctx, in, options.Find().SetProjection(bson.M{
		"id":    1,
		"title": 1,
	}))
	if err != nil {
		panic(err)
	}
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", ars)

	// 改
	sets := bson.D{
		bson.E{
			Key: "$set",
			Value: bson.E{
				Key:   "title",
				Value: "new title",
			},
		},
	}
	updateRes, err := col.UpdateMany(ctx, filter, sets)
	if err != nil {
		panic(err)
	}
	fmt.Printf("affected %d \n", updateRes.ModifiedCount)
	updateRes, err = col.UpdateMany(ctx, filter, bson.D{
		bson.E{
			Key: "$set",
			Value: Article{
				Title:    "new title 2",
				AuthorId: 123,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("affected %d \n", updateRes.ModifiedCount)

	// 删除
	delRes, err := col.DeleteMany(ctx, filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("deleted %d \n", delRes.DeletedCount)

	// 删所有
	_, err = col.DeleteMany(ctx, bson.D{})

}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	CreateAt int64  `bson:"create_at,omitempty"`
	UpdateAt int64  `bson:"update_at,omitempty"`
}
