package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"orderbook/core/models"
)

type OrderRepository struct {
	Collection *mongo.Collection
}

var OrderRepo *OrderRepository
var ctx = context.TODO()

func InitOrderRepo(db *mongo.Database) {
	OrderRepo = &OrderRepository{db.Collection(models.OrderCollectionName)}
}

func (r *OrderRepository) Create(o *models.Order) error {
	_, err := r.Collection.InsertOne(ctx, o)
	if err != nil {
		log.Print("failed to add", o.String(), " error:", err.Error())
	}
	return err
}

func (r *OrderRepository) Update(o *models.Order) error {
	o.UpdatedAt = time.Now()
	_, err := r.Collection.ReplaceOne(ctx, bson.D{{"_id", o.ID}}, o)
	if err != nil {
		log.Print("failed to update", o.String(), " error:", err.Error())
	}
	return err
}

func (r *OrderRepository) GetAll() ([]*models.Order, error) {
	filter := bson.D{{}}
	return r.GetByFilter(filter)
}

func (r *OrderRepository) GetByFilter(filter interface{}) ([]*models.Order, error) {
	var orders []*models.Order

	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return orders, err
	}

	for cur.Next(ctx) {
		var o models.Order
		err := cur.Decode(&o)
		if err != nil {
			return orders, err
		}

		orders = append(orders, &o)
	}

	if err := cur.Err(); err != nil {
		return orders, err
	}

	_ = cur.Close(ctx)

	if len(orders) == 0 {
		return orders, mongo.ErrNoDocuments
	}

	return orders, nil
}
