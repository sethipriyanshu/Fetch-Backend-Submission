package controllers

import (
	"context"
	"go-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// declare mongoDB global variable
var client *mongo.Client

// set the client to be used by controller functions
func SetClient(mongoClient *mongo.Client) {
	client = mongoClient
}

// handle the /add endpoint
func AddPoints(c echo.Context) error {
	// 10s timeout in case the code takes to long to execute
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// bind the json body to a struct
	var transaction models.Transaction
	if err := c.Bind(&transaction); err != nil {
		// handle 400 error
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// get the transactions collection from mongoDB
	collection := client.Database("myapp_db").Collection("transactions")

	// insert the new transaction
	_, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		// error if insertion fails
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert transaction"})
	}

	// 200 response if success
	return c.NoContent(http.StatusOK)
}

// handle the /spend endpoint
func SpendPoints(c echo.Context) error {
	// 10s timeout in case code takes too long to execute
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// bind the json body to a map
	var spendRequest map[string]int
	if err := c.Bind(&spendRequest); err != nil {
		// handle 400 error
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	pointsToSpend := spendRequest["points"]

	// get the transactions collection from mongoDB
	collection := client.Database("myapp_db").Collection("transactions")

	// sort and fetch transactions according to timestamp
	findOptions := options.Find().SetSort(bson.D{{"timestamp", 1}})
	cur, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		// error if fetching fails
		return c.String(http.StatusInternalServerError, "Failed to fetch transactions")
	}

	defer cur.Close(ctx)

	// parse transaction in a slice
	var transactions []models.Transaction
	if err = cur.All(ctx, &transactions); err != nil {
		// error if parsing fails
		return c.String(http.StatusInternalServerError, "Failed to parse transactions")
	}

	// calculate total available points for each payer
	payerPoints := make(map[string]int)
	for _, transaction := range transactions {
		payerPoints[transaction.Payer] += transaction.Points
	}

	// checking for negative points
	totalAvailablePoints := 0
	for _, points := range payerPoints {
		totalAvailablePoints += points
	}
	if pointsToSpend > totalAvailablePoints {
		// insufficient points error
		return c.String(http.StatusBadRequest, "Insufficient points")
	}

	// execute /spend
	spendResults, err := models.SpendPoints(ctx, collection, transactions, pointsToSpend)
	if err != nil {
		// error if spending fails
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// return 200 and results after spending
	return c.JSON(http.StatusOK, spendResults)
}

// handle /balance endpoint
func GetBalance(c echo.Context) error {
	// 10s timeout in case the code takes too long to execute
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get the transactions collection from mongodb
	collection := client.Database("myapp_db").Collection("transactions")

	// fetch points for all payers
	balance, err := models.GetBalance(ctx, collection)
	if err != nil {
		// error if fetching fails
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch balance"})
	}

	// Return 200 and results
	return c.JSON(http.StatusOK, balance)
}
