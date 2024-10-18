package store

import (
	"context"

	"github.com/karokojnr/bookmesh-orders/internal/types"
	pb "github.com/karokojnr/bookmesh-shared/proto"
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

func (s *storage) Create(ctx context.Context, o types.Order) (primitive.ObjectID, error) {
	col := s.db.Database(databaseName).Collection(collectionName)

	newOrder, err := col.InsertOne(ctx, o)

	id := newOrder.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *storage) Get(ctx context.Context, orderId, customerId string) (*types.Order, error) {
	col := s.db.Database(databaseName).Collection(collectionName)

	oId, _ := primitive.ObjectIDFromHex(orderId)

	var o types.Order
	err := col.FindOne(ctx, bson.M{
		"_id":         oId,
		"customer_id": customerId,
	}).Decode(&o)

	return &o, err
}

func (s *storage) Update(ctx context.Context, orderId string, order *pb.Order) error {
	col := s.db.Database(databaseName).Collection(collectionName)

	oId, _ := primitive.ObjectIDFromHex(orderId)

	_, err := col.UpdateOne(ctx,
		bson.M{"_id": oId},
		bson.M{"$set": bson.M{
			"payment_link": order.PaymentLink,
			"status":       order.Status,
		}})

	return err
}
