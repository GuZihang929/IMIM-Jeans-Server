package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://152.32.213.10:27017"))
	if err != nil {
		panic(err)
	}
	fmt.Println("mongo连接成功")
	return client
}
