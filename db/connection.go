package db

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	env "gospotify.com/env"
)

var (
	Client     *mongo.Client
	ClientErr  error
	ClientOnce sync.Once
)

func DbClient() (*mongo.Client, error) {
	ClientOnce.Do(
		func() {
			// mongodb config and connection
			serverAPI := options.ServerAPI(options.ServerAPIVersion1)
			opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s.hfa617f.mongodb.net/?retryWrites=true&w=majority&appName=%s", env.ClusterName, env.UserPassword, strings.ToLower(env.ClusterName), env.ClusterName)).SetServerAPIOptions(serverAPI)

			// create a new client and connect to the server
			Client, ClientErr = mongo.Connect(context.TODO(), opts)
			if ClientErr != nil {
				panic(ClientErr)
			}

			// send a ping to confirm a successful connection
			if err := Client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
				panic(err)
			}
			fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
		})
	return Client, nil
}
