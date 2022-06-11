package main

import (
	"encoding/json"
	"context"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	response, _ := json.Marshal(res)
	c.Send(response)
}

func getPerson(c *fiber.Ctx)  {
	collection, err := getMongoDBCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

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
		c.SendStatus(404)
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

	response, _ := json.Marshal(res)
	c.Send(response)
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

	jsonResponse, _ := json.Marshal(res)
	c.Send(jsonResponse)
}