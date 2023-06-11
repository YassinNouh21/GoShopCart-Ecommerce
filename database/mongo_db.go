package database

import "go.mongodb.org/mongo-driver/mongo"

// DatabaseCollection holds the database collections.
type DatabaseCollection struct {
	UserCollection    *mongo.Collection
	ProductCollection *mongo.Collection
}

// DB holds the instance of the DatabaseCollection used in the project.
var DB *DatabaseCollection

// InitializeDatabase initializes the database collections with the provided user and product collections.
func InitializeDatabase(userCollection, productCollection *mongo.Collection) {
	DB = &DatabaseCollection{
		UserCollection:    userCollection,
		ProductCollection: productCollection,
	}
}
