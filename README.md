# Orderbook
Orderbook is a simple implementation of limit-order-book, written in Go. It simulates the ordering system using Kafka and MongoDB. It also stores and updates orderbooks using some in-memory algorithms and data structures.

## Abstract
Orderbook consists of two applications: one for producing random orders and sending them to a Kafka topic, and one for consuming the orders, storing them in MongoDB, and updating the orderbook.

1. Run `cmd/orderproducer/main.go` for the producer application. This will create a Kafka producer and a scheduler that will produce a random order every second and send it to the orders topic. The orders are in four sample branches for testing purposes.

2. Run `cmd/server/main.go` for the main application. This will create a Kafka consumer that will read the orders from the orders topic and process them. It also exposes an endpoint for getting the orderbook.

`GET /orderbook`: returns the entire orderbook as a JSON object.
#### input
```json
{
  "symbol": "ETHIRT"
  "limit":  10
}
```
#### output
```json
{
    "data": {
        "bids": [
            [
                "1550006",
                "0.325300"
            ]
        ],
        "asks": [
            [
                "1550008",
                "0.679650"
            ]
        ],
        "minAsk": 1550008,
        "maxBid": 1550006
    }
}
```

The orderbook is stored in memory using three data models: PricePoint, Order, and Orderbook. 
### PricePoint
Each `PricePoint` represents a discreet limit price. For products that need to be priced at higher precision, this might not be a scalable solution, but for this exercise, I'm going to assume whatever we're trading is 8 or 9 digits is the max we'll need and can easily be handled in memory. Each PricePoint just contains pointers to the first and last order entered at that price for ease of executing and appending new orders.
```go
type PricePoint struct {
	OrderHead *Order
	OrderTail *Order
}
```
### Order
Each `Order` is either a buy or sell, and has a limit price and amount. also, each of them is linked to the next order at the same price point so we can ensure orders are examined in the order they are entered.
```go
type Order struct {
	ID            string 
	Side          Side        // Buy, Sell    
	Symbol        Symbol      // BTCUSDT, ETHUSDT, ...
	Amount        float64
	Price         uint32
	Next          *Order
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
```
### Orderbook
`OrderBooks` does most of the heavy lifting by keeping track of the current maximum bid and minimum ask, an index mapping order IDs to pointers, and an array of all possible price points.
```go
type Orderbook struct {
	MinAsk      uint32
	MaxBid      uint32
	OrderIndex  map[ObjectID]*Order
	PricePoints [MaxPrice]*PricePoint
}
```

## Algorithms and Complexity

These are the basic operations for an orderbook:
- Inserting a new order into the book.
- Filling an order by matching it with the best available opposite order.
- Updating an existing order by changing its price or volume.
- Deleting an existing order from the book.

But the inserting and the filling operations are just implemented in this project.

### Inserting
Inserting a new order into the book is just a matter of appending it to the list of orders at its price point, updating the order book's bid or asking if necessary, adding an entry in the order index, and inserting the record to Mongo. These are all O(1).

### Filling
Filling orders in the case of a Sell is done by starting at the max bid and iterating over all open orders until we either fill the order or exhaust the open orders at that price point. Then we move down to the next price point and repeat. 

The performance of Filling depends on how sparse the order book is at the time. In the worst case, we'd need to iterate over every price point, however as the number of orders in the book increases the number of price points we need to examine approaches 1. All of the other operations are constant time so this can also be done in O(1).

## Next steps
Some possible improvements and extensions for this project:
- Using a sorted set data structure or Redis to store the price points, instead of an array. This would reduce the space complexity and handle the cases where the price range or precision is not fixed.
- Making the queries more ACID-compliant and the Kafka consumer more scalable and fault-tolerant.

















