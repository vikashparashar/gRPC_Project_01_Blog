package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func StratDatabase() (*mongo.Client, *mongo.Collection) {
	cli, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatalln("error while setting up client for mongodb", err)
	}
	err = cli.Connect(context.Background())
	if err != nil {
		log.Fatalln("error while connection to db", err)
	}
	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalln("error while ping", err)
	}
	coll := cli.Database("testdb").Collection("test_collection")
	n, err := coll.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		panic("can't read total no. of document available in collection")
	}
	fmt.Printf("_____ Total No. Of Document Available Into Database Are : %v_____\n", n)
	return cli, coll
}
