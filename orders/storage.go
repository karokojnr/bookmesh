package main

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseName   = "bookmesh-orders"
	collectionName = "orders"
)

var orders = make([]*pb.Order, 0)

type storage struct {
	db *mongo.Client
}

func NewStorage(db *mongo.Client) *storage {
	return &storage{db}
}

func (s *storage) Create(ctx context.Context, o Order) (primitive.ObjectID, error) {
	col := s.db.Database(databaseName).Collection(collectionName)

	newOrder, err := col.InsertOne(ctx, o)

	id := newOrder.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *storage) Get(ctx context.Context, orderId, customerId string) (*Order, error) {
	col := s.db.Database(databaseName).Collection(collectionName)

	oID, _ := primitive.ObjectIDFromHex(orderId)

	var o Order
	err := col.FindOne(ctx, bson.M{
		"_id":         oID,
		"customer_id": customerId,
	}).Decode(&o)

	return &o, err
}

func (s *storage) Update(ctx context.Context, orderId string, order *pb.Order) error {
	col := s.db.Database(databaseName).Collection(collectionName)

	oID, _ := primitive.ObjectIDFromHex(orderId)

	_, err := col.UpdateOne(ctx,
		bson.M{"_id": oID},
		bson.M{"$set": bson.M{
			"payment_link": order.PaymentLink,
			"status":       order.Status,
		}})

	return err
}
