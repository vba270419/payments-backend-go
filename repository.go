package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type PaymentRepository interface {
	InsertPayment(payment Payment) (err error)

	UpdatePayment(payment Payment) (err error)

	DeletePayment(paymentId string) (err error)

	GetPayment(paymentId string) (payment Payment, err error)

	GetAllPayments() (payments []Payment, err error)
}

type MongoClientProvider struct {
	client *mongo.Client
}

func GetCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("account_book").Collection("payments")
}

func InitializeMongoRepository() (PaymentRepository, *mongo.Client, error) {
	host := "localhost"
	port := "27017"

	log.Printf("Initialize connection to MongoDB [%s:%s] ... ", host, port)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatalf("Failed to establish connection to MongoDB [%s:%s]: %s", host, port, err.Error())
	}

	repository := &MongoClientProvider{client: mongoClient}
	log.Printf("Connection to MongoDB [%s:%s] successfully established", host, port)

	return repository, mongoClient, err
}

func ShutdownMongoRepository(client *mongo.Client) {
	log.Println("Disconnecting from to MongoDB ... ")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = client.Disconnect(ctx)
	log.Println("Disconnected from to MongoDB")
}

func (m *MongoClientProvider) InsertPayment(payment Payment) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := GetCollection(m.client)

	_, err = collection.InsertOne(ctx, payment)
	if err != nil {
		log.Printf("Unexpected error while inserting: %s", err.Error())
		return &PersistenceError{}
	}

	return err
}

func (m *MongoClientProvider) UpdatePayment(payment Payment) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := GetCollection(m.client)

	currentVersion := payment.Version
	payment.Version = payment.Version + 1

	// Here we use version of payment for optimistic locking
	filter := bson.M{"_id": payment.ID, "version": currentVersion}
	update := bson.M{"$set": payment}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Unexpected error while updating: %s", err.Error())
		return &PersistenceError{}
	}

	if result.MatchedCount == 0 {
		_, err = m.GetPayment(payment.ID)
		if err != nil {
			return err
		}
		return &VersionConflictError{payment.ID, currentVersion}
	}
	return err
}

func (m *MongoClientProvider) DeletePayment(paymentId string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := GetCollection(m.client)

	filter := bson.M{"_id": paymentId}

	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		log.Printf("Unexpected error while deleting: %s", err.Error())
		return &PersistenceError{}
	}

	if result.DeletedCount == 0 {
		return &NotFoundError{paymentId}
	}

	return err
}

func (m *MongoClientProvider) GetPayment(paymentId string) (payment Payment, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := GetCollection(m.client)

	filter := bson.M{"_id": paymentId}
	err = collection.FindOne(ctx, filter).Decode(&payment)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return payment, &NotFoundError{paymentId}
		}
		log.Printf("Unexpected error while loading: %s", err.Error())
	}
	return payment, err
}

func (m *MongoClientProvider) GetAllPayments() (payments []Payment, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := GetCollection(m.client)

	filter := bson.M{}
	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		log.Printf("Unexpected error while loading: %s", err.Error())
		return payments, &PersistenceError{}
	}

	for cursor.Next(ctx) {
		var payment Payment
		err = cursor.Decode(&payment)
		if err != nil {
			log.Printf("Unexpected error while loading: %s", err.Error())
			break
		}
		payments = append(payments, payment)
	}

	_ = cursor.Close(ctx)
	return payments, err
}
