// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // This is a user defined method to close resources.
// // This method closes mongoDB connection and cancel context.
// func close(client *mongo.Client, ctx context.Context,
// 	cancel context.CancelFunc) {

// 	defer cancel()

// 	defer func() {
// 		if err := client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()
// }

// // This is a user defined method that returns mongo.Client,
// // context.Context, context.CancelFunc and error.
// // mongo.Client will be used for further database operation.
// // context.Context will be used set deadlines for process.
// // context.CancelFunc will be used to cancel context and
// // resource associated with it.
// func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

// 	ctx, cancel := context.WithTimeout(context.Background(),
// 		30*time.Second)
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
// 	return client, ctx, cancel, err
// }

// // insertMany is a user defined method, used to insert
// // documents into collection returns result of
// // InsertMany and error if any.
// func insertMany(client *mongo.Client, ctx context.Context, dataBase, col string, docs []interface{}) (*mongo.InsertManyResult, error) {

// 	// select database and collection ith Client.Database
// 	// method and Database.Collection method
// 	collection := client.Database(dataBase).Collection(col)

// 	// InsertMany accept two argument of type Context
// 	// and of empty interface
// 	result, err := collection.InsertMany(ctx, docs)
// 	return result, err
// }

// func main() {

// 	// get Client, Context, CancelFunc and err from connect method.
// 	client, ctx, cancel, err := connect("mongodb://localhost:27017")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Release resource when main function is returned.
// 	defer close(client, ctx, cancel)

// 	// Now will be inserting multiple documents into
// 	// the collection. create  a object of type slice
// 	// of interface to store multiple  documents
// 	var comments []interface{}

// 	// Storing into interface list.
// 	comments = []interface{}{
// 		bson.D{
// 			{"commentid", "63499463"},
// 			{"commentText", "The Force member"},
// 		},
// 		bson.D{
// 			{"commentid", "60111106"},
// 			{"commentText", "Another the Force member"},
// 		},
// 	}

// 	// insertMany insert a list of documents into
// 	// the collection. insertMany accepts client,
// 	// context, database name collection name
// 	// and slice of interface. returns error
// 	// if any and result of multi document insertion.
// 	insertManyResult, err := insertMany(client, ctx, "local",
// 		"comments", comments)

// 	// handle the error
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Result of InsertMany")

// 	// print the insertion ids of the multiple
// 	// documents, if they are inserted.
// 	for id := range insertManyResult.InsertedIDs {
// 		fmt.Println(id)
// 	}
// }
