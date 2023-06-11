package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/* User Package user provides the User model for representing user data.

 The User struct represents a user entity in the application, including their personal information,
authentication credentials, access tokens, and related data such as addresses, orders, and cart items.
 Fields:
- ID: The unique identifier of the user.
- FirstName: The first name of the user. Must be between 3 and 20 characters.
- LastName: The last name of the user.
- Email: The email address of the user.
- Password: The password of the user. Must be at least 6 characters.
- Token: The access token associated with the user.
- RefreshToken: The refresh token associated with the user.
- CreatedAt: The timestamp indicating the creation time of the user.
- UpdatedAt: The timestamp indicating the last update time of the user.
- UserID: The user ID associated with the user.
- AddressDetails: The list of addresses associated with the user.
- OrderStatus: The list of order statuses associated with the user.
- UserCart: The list of cart items associated with the user.
*/

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	FirstName      string             `json:"first_name" validate:"required,min=3,max=20"`
	LastName       string             `json:"last_name" bson:"last_name"`
	Email          string             `json:"email" validate:"email" bson:"email"`
	Password       string             `json:"password" bson:"password" validate:"required,min=6"`
	Token          string             `json:"token" bson:"token"`
	RefreshToken   string             `json:"refreshtoken" bson:"refreshtoken"`
	CreatedAt      time.Time          `json:"createdat" bson:"createdat"`
	UpdatedAt      time.Time          `json:"updatedat" bson:"updatedat"`
	UserID         string             `json:"userid"  bson:"userid"`
	AddressDetails []Address          `json:"address" bson:"addressdetails"`
	OrderStatus    []Order            `json:"orderstatus"`
	UserCart       []Cart             `json:"usercart"`
}
