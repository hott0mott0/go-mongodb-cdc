package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb://root:password@mongo:27017/?replicaSet=my-replica-set&authSource=admin"
	databaseName   = "purchase"
	collectionName = "subscriptionOrders"
	projectID      = "emulator"
	topicID        = "event-topic-local"
)

func main() {
	ctx := context.Background()

	// MongoDBクライアントの作成
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}
	defer mongoClient.Disconnect(ctx)
	mongoCollection := mongoClient.Database(databaseName).Collection(collectionName)

	log.Println("Connected to MongoDB")

	// PubSubクライアントの作成
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create PubSub client: %v", err)
	}
	defer pubsubClient.Close()
	topic := pubsubClient.Topic(topicID)

	log.Println("Sucessed to create PubSub client")

	item, err := mongoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to find items: %v", err)
	}
	defer item.Close(ctx)

	log.Println("Sucessed to find items")

	// ChangeStreamの開始
	changeStream, err := mongoCollection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open change stream: %v", err))
	}
	defer changeStream.Close(ctx)

	log.Printf("Start watching changes and resume token is %s", changeStream.ResumeToken().String())

	if err := publishChangeEventToPubSub(ctx, changeStream, topic); err != nil {
		panic(fmt.Sprintf("Failed to publish change event to PubSub: %v", err))
	}
}

func publishChangeEventToPubSub(ctx context.Context, stream *mongo.ChangeStream, topic *pubsub.Topic) error {
	for stream.Next(ctx) {
		var changeEvent bson.M
		if err := stream.Decode(&changeEvent); err != nil {
			log.Printf("Failed to decode change event: %v", err)
			return err
		}

		eventData, err := json.Marshal(changeEvent)
		if err != nil {
			log.Printf("Failed to marshal change event: %v", err)
			return err
		}

		log.Printf("Change event: %s", string(eventData))

		result := topic.Publish(ctx, &pubsub.Message{
			Data: eventData,
		})

		id, err := result.Get(ctx)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			return err
		}
		log.Printf("Published message with ID: %s", id)
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return nil
}
