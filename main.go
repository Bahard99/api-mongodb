package main

import (
	"encoding/json"
	"context"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	InsertedID		string	`json:"insertedid"`
	DeletedCount	int		`json:"deletedcount"`
}

func main()  {
	app := fiber.New()
	app.Get("/person/:id?", getPerson)
	app.Post("/person", createPerson)
	app.Put("/person/:id", updatePerson)
	app.Delete("/person/:id", deletePerson)
	app.Listen(port)
}

func createPerson(c *fiber.Ctx)  {
	collection, err := getMongoDBCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var person Person
	json.Unmarshal([]byte(c.Body()), &person)

	res, err := collection.InsertOne(context.Background(), person)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var response Response
	resp, _ := json.Marshal(res)

	json.Unmarshal(resp, &response)
	c.Send("Success Inserted with ID ", response.InsertedID)
}

func getPerson(c *fiber.Ctx)  {
	collection, err := getMongoDBCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// idcoba, _ := primitive.ObjectIDFromHex("0")
	var filter bson.M = bson.M{}

	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		// c.SendStatus(404)
		c.Status(404).Send("NO DATA FOUND")
		return
	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

func updatePerson(c *fiber.Ctx)  {
	collection, err := getMongoDBCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	
	var person Person
	json.Unmarshal([]byte(c.Body()), &person)

	update := bson.M{
		"$set": person,
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var response Response
	resp, _ := json.Marshal(res)

	json.Unmarshal(resp, &response)
	c.Send("Success Update with ID ", response.InsertedID)
}

func deletePerson(c *fiber.Ctx)  {
	collection, err := getMongoDBCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var response Response
	resp, _ := json.Marshal(res)

	json.Unmarshal(resp, &response)
	c.Send("Success Delete ", response.DeletedCount, " Data")
}