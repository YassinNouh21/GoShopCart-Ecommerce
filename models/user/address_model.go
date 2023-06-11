package user

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
   Address represents a user's address, including the address ID, street, city, state, postal code, and country code.

   Fields:
   - AddressID: Unique identifier for the address.
   - Street: Street name of the address. Required field and validated accordingly.
   - City: City of the address. Required field and validated accordingly.
   - State: State or province of the address. Required field and validated accordingly.
   - PostalCode: Postal code of the address. Required field and validated accordingly.
   - CountryCode: Country code of the address. Required field and validated accordingly.
*/

type Address struct {
	AddressID   primitive.ObjectID `bson:"_id"`
	Street      string             `json:"street" validate:"required" bson:"street"`
	City        string             `json:"city" validate:"required" bson:"city"`
	State       string             `json:"state" validate:"required" bson:"state"`
	PostalCode  string             `json:"postal_code" validate:"required" bson:"postal_code"`
	CountryCode string             `json:"country_code" validate:"required" bson:"country_code"`
}
