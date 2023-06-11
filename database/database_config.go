package database

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoInstance() *mongo.Client {
	godotenv.Load(".env")
	uri := os.Getenv("MONGO_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	if err := client.Database("e-commerce").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	return client
}

var MongoDBInstance *mongo.Client = MongoInstance()

func GetCollectionMongoDB(collectionName string) *mongo.Collection {
	collection := MongoDBInstance.Database("e-commerce").Collection(collectionName)
	return collection
}

func InitializeMongoDBCollections() {
	mongoDBCollectionProducts := MongoDBInstance.Database("e-commerce").Collection("products")
	mongoDBCollectionUsers := MongoDBInstance.Database("e-commerce").Collection("users")

	InitializeDatabase(mongoDBCollectionUsers, mongoDBCollectionProducts)
}
