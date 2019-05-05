package main

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// PaymentRepository is an interface which defines the methods must be implemented by a specific repository that persist payments to storage
type PaymentRepository interface {
	InsertPayment(payment Payment) (err error)

	UpdatePayment(payment Payment) (err error)

	DeletePayment(paymentID string) (err error)

	GetPayment(paymentID string) (payment Payment, err error)

	GetAllPayments() (payments []Payment, err error)
}

type mongoClient struct {
	client *mongo.Client
}

func (m *mongoClient) InsertPayment(payment Payment) (err error) {
	collection := getCollection(m.client)

	_, err = collection.InsertOne(getContextWithTimeout(), payment)
	if err != nil {
		log.Printf("Unexpected error while inserting: %s", err.Error())
		return &PersistenceError{}
	}

	return err
}

func (m *mongoClient) UpdatePayment(payment Payment) (err error) {
	collection := getCollection(m.client)

	currentVersion := payment.Version
	payment.Version = payment.Version + 1

	// Here we use version of payment for optimistic locking
	filter := bson.M{"_id": payment.ID, "version": currentVersion}
	update := bson.M{"$set": payment}

	result, err := collection.UpdateOne(getContextWithTimeout(), filter, update)
	if err != nil {
		log.Printf("Unexpected error while updating: %s", err.Error())
		return &PersistenceError{}
	}

	if result.MatchedCount == 0 {
		_, err = m.GetPayment(payment.ID)
		if err != nil {
			return err
		}
		return &PaymentVersionConflictError{payment.ID, currentVersion}
	}
	return err
}

func (m *mongoClient) DeletePayment(paymentID string) (err error) {
	collection := getCollection(m.client)

	filter := bson.M{"_id": paymentID}

	result, err := collection.DeleteOne(getContextWithTimeout(), filter)

	if err != nil {
		log.Printf("Unexpected error while deleting: %s", err.Error())
		return &PersistenceError{}
	}

	if result.DeletedCount == 0 {
		return &PaymentNotFoundError{paymentID}
	}

	return err
}

func (m *mongoClient) GetPayment(paymentID string) (payment Payment, err error) {
	collection := getCollection(m.client)

	filter := bson.M{"_id": paymentID}
	err = collection.FindOne(getContextWithTimeout(), filter).Decode(&payment)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return payment, &PaymentNotFoundError{paymentID}
		}
		log.Printf("Unexpected error while loading: %s", err.Error())
	}
	return payment, err
}

func (m *mongoClient) GetAllPayments() (payments []Payment, err error) {
	ctx := getContextWithTimeout()
	collection := getCollection(m.client)

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

func getContextWithTimeout() context.Context {
	duration := time.Duration(viper.GetInt(mongoDbTimeout)) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), duration)
	return ctx
}

func getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("account_book").Collection("payments")
}

func initializeMongoRepository() (PaymentRepository, *mongo.Client) {
	host := viper.GetString(mongoDbHost)
	port := viper.GetString(mongoDbPort)

	log.Printf("Connecting to MongoDB [%s:%s] ... ", host, port)
	ctx := getContextWithTimeout()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port))
	if err != nil {
		log.Fatalf("Failed to establish connection to MongoDB [%s:%s]: %s", host, port, err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to establish connection to MongoDB [%s:%s]: %s", host, port, err.Error())
	}

	repository := &mongoClient{client: client}
	log.Printf("Connection to MongoDB [%s:%s] - OK", host, port)

	return repository, client
}

func shutdownMongoRepository(client *mongo.Client) {
	log.Println("Disconnecting from to MongoDB ... ")
	_ = client.Disconnect(getContextWithTimeout())
	log.Println("Disconnected from to MongoDB")
}
