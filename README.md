# orders-api 
A Microservice that handles the management of customer orders. It uses Redis to store and retrieve data. It exposes a set of RESTful APIs through which other microservices or clients can interact with it. It is written in Go, utilizing the [chi](https://pkg.go.dev/github.com/go-chi/chi/v5) library.

### Requirements

- Golang v1.20 or higher
- Redis v7

## Getting Started
To run the application:
1. Clone this repository
2. Run `go mod download` command to download all the dependencies.
3. Set up Redis.
4. Provide values for SERVER_PORT and REDIS_ADDR in the `.env` file. Otherwise, the defaults `:8000` and `localhost:6379` are used as sever port and redis address respectively.
5. Run `go run main.go` command to start the microservice. You can now make http requests using the API below.

## API Endpoints

| Endpoint | Method | Description |
| -------- | ------ | ----------- |
| /orders   | POST   | Insert a new order |
| /orders   | GET    | Get list of all orders |
| /orders/{id} | GET    | Retrieve an order by id |
| /orders/{id} | PUT    | Update an order by id |
| /orders/{id} | DELETE | Delete an order by id |

## Usage Examples
### 1. Insert a new Order
#### *Request:*
POST /orders
```json
{
    "customer_id": "36c37006-7577-4456-b295-3eafc8b17100",
    "line_items": [
        {
            "item_id": "23d77a68-8246-43c0-a0ea-ec81dd1742cd",
            "quantity": 4,
            "price": 7511
        },
        {
            "item_id": "2e0687cf-57ee-48d9-99d9-1d45b1dd9f30",
            "quantity": 2,
            "price": 2214
        }
    ]
}
```
#### *Response:*
201 Created
```json
{
    "order_id": 5294788243102287312,
    "customer_id": "36c37006-7577-4456-b295-3eafc8b17100",
    "line_items": [
        {
            "item_id": "23d77a68-8246-43c0-a0ea-ec81dd1742cd",
            "quantity": 4,
            "price": 7511
        },
        {
            "item_id": "2e0687cf-57ee-48d9-99d9-1d45b1dd9f30",
            "quantity": 2,
            "price": 2214
        }
    ],
    "created_at": "2023-12-27T11:35:04.03181922Z",
    "shipped_at": null,
    "completed_at": null
}
```
### 2. List all Orders
#### *Request:*
GET /orders
#### *Response:*
200 OK
```json
{
    "items": [
        
        {
            "order_id": 13286624006422104208,
            "customer_id": "c1bb76c0-f779-4e8a-aee3-11edcf46ccf7",
            "line_items": [
                {
                    "item_id": "44df2743-a68a-4a17-a111-4c8ff242b6d8",
                    "quantity": 10,
                    "price": 4110
                },
                {
                    "item_id": "9eabe57d-2b27-40ea-880b-9665adcbcb95",
                    "quantity": 8,
                    "price": 9516
                }
            ],
            "created_at": "2023-12-25T19:01:56.165092312Z",
            "shipped_at": null,
            "completed_at": null
        },
        {
            "order_id": 5294788243102287312,
            "customer_id": "36c37006-7577-4456-b295-3eafc8b17100",
            "line_items": [
                {
                    "item_id": "23d77a68-8246-43c0-a0ea-ec81dd1742cd",
                    "quantity": 4,
                    "price": 7511
                },
                {
                    "item_id": "2e0687cf-57ee-48d9-99d9-1d45b1dd9f30",
                    "quantity": 2,
                    "price": 2214
                }
            ],
            "created_at": "2023-12-27T11:35:04.03181922Z",
            "shipped_at": null,
            "completed_at": null
        }
    ]
}
```
### 3. Retrieve an Order by id
#### *Request:*
GET /orders/5294788243102287312
#### *Response:*
200 OK
```json
{
    "order_id": 5294788243102287312,
    "customer_id": "36c37006-7577-4456-b295-3eafc8b17100",
    "line_items": [
        {
            "item_id": "23d77a68-8246-43c0-a0ea-ec81dd1742cd",
            "quantity": 4,
            "price": 7511
        },
        {
            "item_id": "2e0687cf-57ee-48d9-99d9-1d45b1dd9f30",
            "quantity": 2,
            "price": 2214
        }
    ],
    "created_at": "2023-12-27T11:35:04.03181922Z",
    "shipped_at": null,
    "completed_at": null
}
```
### 4. Update an Order by id
#### *Request:*
PUT /orders/5294788243102287312

The `status` field in the request body should only be set to `shipped` or `completed`.
```json
{
    "status": "shipped"
}
```
#### *Response:*
200 OK
```json
{
    "order_id": 5294788243102287312,
    "customer_id": "36c37006-7577-4456-b295-3eafc8b17100",
    "line_items": [
        {
            "item_id": "23d77a68-8246-43c0-a0ea-ec81dd1742cd",
            "quantity": 4,
            "price": 7511
        },
        {
            "item_id": "2e0687cf-57ee-48d9-99d9-1d45b1dd9f30",
            "quantity": 2,
            "price": 2214
        }
    ],
    "created_at": "2023-12-27T11:35:04.03181922Z",
    "shipped_at": "2023-12-28T08:48:10.320806736Z",
    "completed_at": null
}
```
### 5. Delete an Order by id
#### *Request:*
DELETE /orders/13286624006422104208
#### *Response:*
200 OK