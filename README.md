# GoShopCart-Ecommerce

## Description

This project is an e-commerce back-end developed using the Golang Gin framework. It provides various APIs for user authentication, profile management, address management, cart management, and product management.

## Installation

1. Clone the repository:
   
   ```bash
   git clone https://github.com/YassinNouh21/GoShopCart-Ecommerce.git
   ```

2. Install the dependencies:

   ```bash
   go mod download
   ```

## Usage

1. Build the project:

   ```bash
   go build
   ```

2. Run the application:

   ```bash
   ./your-project
   ```

3. The application will start running on `http://localhost:8080`.

## Database Schema

The following diagram represents the database schema of the GoShopCart E-commerce API:

![Database Schema](https://github.com/YassinNouh21/GoShopCart-Ecommerce/blob/main/GoShopCart-Ecommerce.png?raw=true)

Please refer to the diagram for a visual representation of the relationships between the different entities in the database.


## Endpoints

The following endpoints are available in the application:

- `POST   /auth/signin` - Signs in the user.
- `POST   /auth/signup` - Signs up a new user.
- `POST   /auth/tokenrefresh` - Refreshes the authentication token.
- `GET    /user/profile` - Retrieves the user's profile information.
- `POST   /user/profile/update` - Updates the user's profile information.
- `GET    /user/address` - Retrieves the user's address information.
- `POST   /user/address` - Adds a new address for the user.
- `DELETE /user/address` - Deletes all addresses of the user.
- `DELETE /user/address/:address_id` - Deletes a specific address of the user.
- `GET    /user/cart` - Retrieves the user's cart information.
- `POST   /user/cart` - Adds a product to the user's cart
- `DELETE /user/cart` - Deletes all products from the user's cart.
- `PUT    /user/cart/:cart_id` - Updates a specific product in the user's cart.
- `POST   /product/` - Creates a new product.
- `GET    /product/:id` - Retrieves a specific product.
- `PUT    /product/:id` - Updates a specific product.
- `DELETE /product/:id` - Deletes a specific product.
- `GET    /product/price` - Retrieves products within a price range.
- `GET    /product/price/:price` - Retrieves products by price.
- `GET    /product/keyword` - Retrieves products by keyword.

## Contributing

Contributions to this project are welcome. Feel free to open a pull request or submit any issues you may encounter.

