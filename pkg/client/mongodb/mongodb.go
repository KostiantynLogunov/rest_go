package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var mongoDbUrl string
	var isAuth bool

	if authDB == "" {
		authDB = database
	}

	if username == "" && password == "" {
		mongoDbUrl = fmt.Sprintf("mongodb://%s%s", host, port) //if we don't have authorization to db
	} else {
		mongoDbUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port) //if we have authorization to db
	}

	clientOptions := options.Client().ApplyURI(mongoDbUrl)
	if isAuth {
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	} else {

	}
	// Connect
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb due to error : %v", err)
	}

	// Ping
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping to mongodb due to error : %v", err)
	}

	return client.Database(database), nil
}
