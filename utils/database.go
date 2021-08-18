package utils

import (
	"context"
	"time"

	"github.com/jintaokoong/go-tcl/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateClient(config structs.Config) (client *mongo.Client, err error) {
	return mongo.NewClient(options.Client().ApplyURI(config.Database.ConnectionString))
}

func CreateDatabaseContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
