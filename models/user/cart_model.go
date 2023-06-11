package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	Cart represents a user's cart item, which includes the cart ID, product ID, quantity, and timestamps for creation and updates.

	CartWithoutId is a variant of the Cart model without the CartID field, used when the cart item doesn't 	require an explicit ID.

	Fields:
	- CartID: The unique identifier of the cart item.
	- ProductID: The identifier of the associated product.
	- Quantity: The quantity of the product in the cart.
	- CreatedAt: The timestamp indicating when the cart item was created.
	- UpdatedAt: The timestamp indicating when the cart item was last updated.

*/

type Cart struct {
	CartID    primitive.ObjectID `json:"cart_id" bson:"_id" validate:"required"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id" validate:"required"`
	Quantity  int                `json:"quantity" bson:"quantity" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" validate:"required"`
}

type CartWithoutId struct {
	CartID    primitive.ObjectID `json:"cart_id" bson:"_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id" validate:"required"`
	Quantity  int                `json:"quantity" bson:"quantity" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" validate:"required"`
}
