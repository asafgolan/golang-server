package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todos Todos

// Give us some seed data
func getCollection(collectionName string) Todos {
	// set ctx ,open connection
	todos = nil
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DBCONNECTIONSTR")))
	//set collection , set ctx
	log.Printf("hey mr tambourine man")
	collection := client.Database("todo").Collection(collectionName)
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	//query collection
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		jsonString, _ := json.Marshal(result)
		// convert json to struct
		s := Todo{}
		json.Unmarshal(jsonString, &s)
		todos = append(todos, s)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return todos
}

//helper functions
func RemoveQoutsFromStr(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func RepoFindTodo(id string) Todo {
	ids, err := primitive.ObjectIDFromHex(id)
	var todo Todo
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DBCONNECTIONSTR")))
	//set collection , set ctx
	log.Printf("hey mr tambourine man")
	collection := client.Database("todo").Collection("todo")
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	//query collection
	cur, err := collection.Find(ctx, bson.M{"_id": ids})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		jsonString, _ := json.Marshal(result)
		// convert json to struct

		json.Unmarshal(jsonString, &todo)
		//todos = append(todos, s)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return todo

}

func RepoCreateTodo(t Todo) Todo {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DBCONNECTIONSTR")))
	if err != nil {
		panic(err)
	}
	//set collection , set ctx

	collection := client.Database("todo").Collection("todo")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, t)

	if err != nil {
		panic(err)
	}

	//fix response with boolean
	if res != nil {
		t.Id = res.InsertedID.(primitive.ObjectID)
		todos = append(todos, t)
	}
	return t
}

func DestroyTodo(id string) int64 {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DBCONNECTIONSTR")))
	//set collection , set ctx
	collection := client.Database("todo").Collection("todo")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	ids, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	// delete document
	res, err := collection.DeleteOne(ctx, bson.M{"_id": ids})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted count: %d\n", res.DeletedCount)

	return res.DeletedCount

}

func UpdateOne(t Todo) Todo {

	//fmt.Println(b)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DBCONNECTIONSTR")))
	if err != nil {
		log.Fatal(err)
	}
	//set collection , set ctx
	col := client.Database("todo").Collection("todo")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// set filters and updates
	//log.Fatal("ggogogo ",t.Id, t.String())
	filter := bson.M{"_id": t.Id}
	update := bson.M{"$set": t}

	res, err := col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	} else {
		//fix the ID
		log.Print("bobo ", res)
		return t
	}
	return t
}
