package controllers

import (
	"context"
	"go-api/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	pointsToSpend := spendRequest["points"]

	collection := client.Database("myapp_db").Collection("transactions")
	transactions, err := models.GetTransactions(ctx, collection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch transactions"})
	}

	spendResults, err := models.SpendPoints(ctx, collection, transactions, pointsToSpend)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
