package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Package user defines the data model for the Order entity.
The Order struct represents an order placed by a user. It contains information such as the order ID, discount, payment method, creation and update timestamps, and the quantity of items in the order.

Fields:
- OrderID: The unique identifier for the order.
- Discount: The discount applied to the order.
- PaymentMethod: The payment method used for the order.
- CreatedAt: The timestamp when the order was created.
- UpdatedAt: The timestamp when the order was last updated.
- Quantity: The quantity of items in the order.
*/
type Order struct {
	OrderID       primitive.ObjectID `json:"order_id" validate:"required" bson:"product_id"`
	Discount      float32            `json:"discount" validate:"required" bson:"discount"`
	PaymentMethod string             `json:"payment_method" validate:"required" bson:"payment_method"`
	CreatedAt     time.Time          `json:"created_at" validate:"required" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" validate:"required" bson:"updated_at"`
	Quantity      int                `json:"quantity" validate:"required" bson:"quantity"`
}
