package product

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
	Package product defines the Product model used in the ecommerce application.

	The Product struct represents a single product with its attributes.

	Fields:
	- ProductID: The unique identifier for the product. It is represented as a primitive.ObjectID.
	- ProductName: The name of the product. It is a required field.
	- Price: The price of the product. It is a required field and represented as a float32.
	- Rating: The rating of the product. It is a required field and represented as a float32.
	- ImageUrl: The URL of the product's image. It is a required field.

	This Product model is used to represent individual products in the ecommerce application.
*/

type Product struct {
	ProductID   primitive.ObjectID `json:"product_id" bson:"_id" validate:"required"`
	ProductName string             `json:"product_name" bson:"product_name" validate:"required"`
	Price       float32            `json:"price" bson:"price" validate:"required"`
	Rating      float32            `json:"rating" bson:"rating" validate:"required"`
	ImageUrl    string             `json:"image" bson:"image_url" validate:"required"`
}
