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

var client *mongo.Client

func SetClient(mongoClient *mongo.Client) {
	client = mongoClient
}

func AddPoints(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var transaction models.Transaction
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	collection := client.Database("myapp_db").Collection("transactions")

	_, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert transaction"})
	}

	return c.NoContent(http.StatusOK)
}

func SpendPoints(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var spendRequest map[string]int
	if err := c.Bind(&spendRequest); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	pointsToSpend := spendRequest["points"]

	collection := client.Database("myapp_db").Collection("transactions")

	// Fetch transactions sorted by timestamp (oldest first)
	findOptions := options.Find().SetSort(bson.D{{"timestamp", 1}})
	cur, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch transactions")
	}
	defer cur.Close(ctx)

	var transactions []models.Transaction
	if err = cur.All(ctx, &transactions); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse transactions")
	}

	// Calculate total available points for each payer
	payerPoints := make(map[string]int)
	for _, transaction := range transactions {
		payerPoints[transaction.Payer] += transaction.Points
	}

	// Check if points to spend exceed available points
	totalAvailablePoints := 0
	for _, points := range payerPoints {
		totalAvailablePoints += points
	}
	if pointsToSpend > totalAvailablePoints {
		return c.String(http.StatusBadRequest, "Insufficient points")
	}

	spendResults, err := models.SpendPoints(ctx, collection, transactions, pointsToSpend)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, spendResults)
}

func GetBalance(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("myapp_db").Collection("transactions")
	balance, err := models.GetBalance(ctx, collection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch balance"})
	}

	return c.JSON(http.StatusOK, balance)
}
