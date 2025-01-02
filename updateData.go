package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns
// mongo.Client, context.Context,
// context.CancelFunc and error.
// mongo.Client will be used for further database
// operation.context.Context will be used set
// deadlines for process. context.CancelFunc will
// be used to cancel context and resource
// associated with it.
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// UpdateOne is a user defined method, that update
// a single document matching the filter.
// This methods accepts client, context, database,
// collection, filter and update filter and update
// is of type interface this method returns
// UpdateResult and an error if any.
func UpdateOne(client *mongo.Client, ctx context.Context, dataBase,
	col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	// select the database and the collection
	collection := client.Database(dataBase).Collection(col)

	// A single document that match with the
	// filter will get updated.
	// update contains the filed which should get updated.
	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

// UpdateMany is a user defined method, that update
// a multiple document matching the filter.
// This methods accepts client, context, database,
// collection, filter and update filter and update
// is of type interface this method returns
// UpdateResult and an error if any.
func UpdateMany(client *mongo.Client, ctx context.Context,
	dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	// select the database and the collection
	collection := client.Database(dataBase).Collection(col)

	// All the documents that match with the filter will
	// get updated.
	// update contains the filed which should get updated.
	result, err = collection.UpdateMany(ctx, filter, update)
	return
}

func main() {

	// get Client, Context, CancelFunc and err from connect method.
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Free the resource when main function in returned
	defer close(client, ctx, cancel)

	// filter object is used to select a single
	// document matching that matches.
	filter := bson.D{{"commentid", "22749003"}}

	// The field of the document that need to updated.
	update := bson.D{
		{"$set", bson.D{
			{"commentText", "This is the updated msg"},
		}},
	}

	// Returns result of updated document and a error.
	result, err := UpdateOne(client, ctx, "local",
		"comments", filter, update)

	// handle error
	if err != nil {
		panic(err)
	}

	// print count of documents that affected
	fmt.Println("update single document")
	fmt.Println(result.ModifiedCount)

	// filter = bson.D{
	//     {"computer", bson.D{{"$lt", 100}}},
	// }
	// update = bson.D{
	//     {"$set", bson.D{
	//         {"computer", 100},
	//     }},
	// }

	// // Returns result of updated document and a error.
	// result, err = Update(client, ctx, "gfg",
	//                      "marks", filter, update)

	// // handle error
	// if err != nil {
	//     panic(err)
	// }

	// // print count of documents that affected
	// fmt.Println("update multiple document")
	// fmt.Println(result.ModifiedCount)
}
